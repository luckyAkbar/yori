package model

import (
	"context"
	"mime/multipart"
)

type FileUsecase interface {
	Upload(ctx context.Context, file *multipart.FileHeader) error
}
