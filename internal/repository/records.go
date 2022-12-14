package repository

import (
	"context"
	"log"

	"github.com/luckyAkbar/yori/internal/model"
	"gorm.io/gorm"
)

type recordRepository struct {
	db *gorm.DB
}

func NewRecordRepository(db *gorm.DB) model.RecordRepository {
	return &recordRepository{
		db,
	}
}

func (r *recordRepository) FindByKTP(ctx context.Context, ktp string) ([]model.Record, error) {
	record := []model.Record{}
	err := r.db.WithContext(ctx).Model(&model.Record{}).Where("ktp = ?", ktp).Find(&record).Error
	switch err {
	default:
		log.Println("error:", err)
		return nil, err
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	case nil:
		return record, nil
	}
}

func (r *recordRepository) FindByKK(ctx context.Context, kk string) ([]model.Record, error) {
	record := []model.Record{}
	err := r.db.WithContext(ctx).Model(&model.Record{}).Where("kk = ?", kk).Find(&record).Error
	switch err {
	default:
		log.Println("error:", err)
		return nil, err
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	case nil:
		return record, nil
	}
}
