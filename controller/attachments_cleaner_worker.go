package controller

import (
	"context"
	"github.com/Falokut/go-kit/log"
)

type AttachmentsCleanerWorker interface {
	Do(ctx context.Context) error
}

type AttachmentsCleanerWorkerHandler struct {
	worker AttachmentsCleanerWorker
	logger log.Logger
}

func NewAttachmentsCleanerWorkerHandler(worker AttachmentsCleanerWorker, logger log.Logger) AttachmentsCleanerWorkerHandler {
	return AttachmentsCleanerWorkerHandler{
		worker: worker,
		logger: logger,
	}
}

func (w AttachmentsCleanerWorkerHandler) Do(ctx context.Context) {
	ctx = log.ToContext(ctx, log.Any("workerName", "attachments_cleaner"))
	err := w.worker.Do(ctx)
	if err != nil {
		w.logger.Error(ctx, "attachments cleaner worker error", log.Any("err", err.Error()))
	}
}
