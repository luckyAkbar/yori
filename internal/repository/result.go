package repository

import (
	"context"

	"github.com/luckyAkbar/yori/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type resultRepo struct {
	db *gorm.DB
}

func NewResultRepo(db *gorm.DB) model.ResultRepository {
	return &resultRepo{
		db,
	}
}

func (r *resultRepo) Create(ctx context.Context, res *model.Result) error {
	if err := r.db.WithContext(ctx).Create(res).Error; err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (r *resultRepo) Update(ctx context.Context, res *model.Result) error {
	if err := r.db.WithContext(ctx).Save(res).Error; err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (r *resultRepo) FindByID(ctx context.Context, id int64) (*model.Result, error) {
	res := &model.Result{}
	err := r.db.WithContext(ctx).Model(&model.Result{}).Where("id = ?", id).Take(res).Error
	switch err {
	default:
		logrus.Error(err)
		return nil, err
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	case nil:
		return res, nil
	}
}
