package service

import (
	"context"
	"courses-service/entity"
	"sync"

	"github.com/pkg/errors"
)

const maxRecordsDelete = 10

type AttachmentsCleanerTxRunner interface {
	CleanAttachmentsTransaction(ctx context.Context, txFunc func(ctx context.Context, tx AttachmentCleanerTx) error) error
}

type AttachmentCleanerTx interface {
	RemoveOrphanedAttachments(ctx context.Context, maxRecordsDelete int) ([]string, error)
}

type AttachmentsCleanerWorker struct {
	txRunner    AttachmentsCleanerTxRunner
	fileStorage FileRepo
}

func NewAttachmentsCleanerWorker(txRunner AttachmentsCleanerTxRunner, fileStorage FileRepo) AttachmentsCleanerWorker {
	return AttachmentsCleanerWorker{
		txRunner:    txRunner,
		fileStorage: fileStorage,
	}
}

func (w AttachmentsCleanerWorker) Do(ctx context.Context) error {
	err := w.txRunner.CleanAttachmentsTransaction(ctx, w.do)
	if err != nil {
		return errors.WithMessage(err, "clean attachments transaction")
	}
	return nil
}

func (w AttachmentsCleanerWorker) do(ctx context.Context, tx AttachmentCleanerTx) error {
	urlsToDelete, err := tx.RemoveOrphanedAttachments(ctx, maxRecordsDelete)
	if err != nil {
		return errors.WithMessage(err, "remove orphaned attachments")
	}
	if len(urlsToDelete) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	errorsCh := make(chan error, len(urlsToDelete))

	for _, url := range urlsToDelete {
		urlCopy := url
		wg.Add(1)
		go func(url string) {
			err := w.deleteFile(ctx, url)
			if err != nil {
				errorsCh <- errors.WithMessage(err, "delete file")
			}
			wg.Done()
		}(urlCopy)
	}

	wg.Wait()
	close(errorsCh)

	var resultErr error
	for err := range errorsCh {
		resultErr = errors.Wrap(resultErr, err.Error())
	}
	return resultErr
}

func (w AttachmentsCleanerWorker) deleteFile(ctx context.Context, url string) error {
	err := w.fileStorage.DeleteFile(ctx, url)
	if err != nil && !errors.Is(err, entity.ErrFileNotFound) {
		return errors.WithMessage(err, "delete file")
	}
	return nil
}
