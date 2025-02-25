package domain

type LimitOffsetRequest struct {
	Limit  int32 `validate:"required,min=10,max=100"`
	Offset int32 `validate:"min=0"`
}
