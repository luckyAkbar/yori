package usecase

import (
	"context"
	"encoding/csv"
	"mime/multipart"

	"github.com/kumparan/go-utils"
	"github.com/luckyAkbar/yori/internal/helper"
	"github.com/luckyAkbar/yori/internal/model"
	"github.com/sirupsen/logrus"
)

type fileUsecase struct {
	recordRepo model.RecordRepository
}

func NewFileUsecase(recordRepo model.RecordRepository) model.FileUsecase {
	return &fileUsecase{
		recordRepo,
	}
}

func (u *fileUsecase) Upload(ctx context.Context, input *multipart.FileHeader) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":  utils.DumpIncomingContext(ctx),
		"file": utils.Dump(input),
	})

	file, err := input.Open()
	if err != nil {
		logger.Warn(err)
		return ErrBadRequest
	}

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		logger.Warn(err)
		return ErrBadRequest
	}

	if err := helper.ValidateCSVHeader(records[0]); err != nil {
		logger.Warn(err)
		return ErrBadRequest
	}

	data, err := helper.GenerateAndValidateRecords(records)
	if err != nil {
		logger.Error(err)
		return ErrBadRequest
	}

	if err := u.recordRepo.SaveBulkWithTransaction(ctx, data); err != nil {
		logger.Error(err)
		return ErrInternal
	}

	return nil
}
