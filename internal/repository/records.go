package repository

import (
	"context"
	"log"
	"time"

	"github.com/kumparan/go-lib/utils"
	"github.com/luckyAkbar/yori/internal/model"
	"github.com/sirupsen/logrus"
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

func (r *recordRepository) SaveBulk(ctx context.Context, bulk []model.Record) error {
	logger := logrus.WithField("ctx", utils.DumpIncomingContext(ctx))
	if err := r.db.WithContext(ctx).CreateInBatches(&bulk, 100).Error; err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (r *recordRepository) SaveBulkWithTransaction(ctx context.Context, bulk []model.Record) error {
	logger := logrus.WithField("ctx", utils.DumpIncomingContext(ctx))

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).CreateInBatches(&bulk, 100).Error; err != nil {
			logger.Error(err)
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (r *recordRepository) CheckIsDealerExists(ctx context.Context, name string) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":  utils.DumpIncomingContext(ctx),
		"name": name,
	})

	record := &model.Record{}
	err := r.db.WithContext(ctx).Model(&model.Record{}).Where("dealer = ?", name).Take(record).Error
	switch err {
	default:
		logger.Error(err)
		return err
	case gorm.ErrRecordNotFound:
		return ErrNotFound
	case nil:
		return nil
	}
}

func (r *recordRepository) FindByDateRange(ctx context.Context, start, finish time.Time) ([]model.Record, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":    utils.DumpIncomingContext(ctx),
		"start":  start,
		"finish": finish,
	})

	records := []model.Record{}
	if err := r.db.WithContext(ctx).Model(&model.Record{}).Where("tgl_mohon_faktur >= ? AND tgl_mohon_faktur <= ?", start, finish).Find(&records).Error; err != nil {
		logger.Error(err)
		return records, err
	}

	return records, nil
}
