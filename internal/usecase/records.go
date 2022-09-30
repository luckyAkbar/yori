package usecase

import (
	"context"
	"log"

	"github.com/luckyAkbar/yori/internal/model"
	"github.com/luckyAkbar/yori/internal/repository"
)

type recordUsecase struct {
	repo model.RecordRepository
}

func NewRecordUsecase(repo model.RecordRepository) model.RecordUsecase {
	return &recordUsecase{
		repo,
	}
}

func (u *recordUsecase) FindByKTP(ctx context.Context, ktp string) (*model.Record, error) {
	record, err := u.repo.FindByKTP(ctx, ktp)
	switch err {
	default:
		log.Println("error:", err)
		return nil, ErrInternal
	case repository.ErrNotFound:
		return nil, ErrNotFound
	case nil:
		return record, nil
	}
}
