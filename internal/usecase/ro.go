package usecase

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/kumparan/go-lib/utils"
	"github.com/luckyAkbar/yori/internal/helper"
	"github.com/luckyAkbar/yori/internal/model"
	"github.com/luckyAkbar/yori/internal/repository"
	"github.com/sirupsen/logrus"
)

type roUsecase struct {
	recordRepo model.RecordRepository
	resultRepo model.ResultRepository
}

func NewROUsecase(recordRepo model.RecordRepository, resultRepo model.ResultRepository) model.ROUsecase {
	return &roUsecase{
		recordRepo,
		resultRepo,
	}
}

func (u *roUsecase) GetResult(ctx context.Context, id int64) ([]model.AdvanceROCheckingResult, *model.Result, error) {
	var result []model.AdvanceROCheckingResult

	res, err := u.resultRepo.FindByID(ctx, id)
	switch err {
	default:
		logrus.Error(err)
		return result, nil, ErrInternal
	case repository.ErrNotFound:
		return result, nil, ErrNotFound
	case nil:
		break
	}

	if res.Result == "" {
		return result, nil, ErrInProgress
	}

	if err := json.Unmarshal([]byte(res.Result), &result); err != nil {
		logrus.Error(err)
		return result, nil, ErrInternal
	}

	return result, res, nil
}

func (u *roUsecase) HandleAdvanceChecking(ctx context.Context, input *model.AdvanceROCheckingInput, id int64) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"input": utils.Dump(input),
	})

	logger.Info("creating emptry result")

	result := &model.Result{
		ID: id,
	}

	logrus.Info("ini id result: ", result.ID)
	if err := u.resultRepo.Create(ctx, result); err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("finish creating empty result")

	logger.Info("start operation")

	base, err := u.recordRepo.FindByDateRange(ctx, input.BaseStart, input.BaseFinish)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("finish find base")

	if len(base) == 0 {
		return ErrNotFound
	}

	logger.Info("start find comparator")

	comparator, err := u.recordRepo.FindByDateRange(ctx, input.CompareStart, input.CompareFinish)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("finish find comparator")

	if len(comparator) == 0 {
		return ErrNotFound
	}

	var res []model.AdvanceROCheckingResult
	cleanBase := helper.RemoveDuplicateKTP(base)
	for i, rec := range cleanBase {
		logger.Info("processing record no: ", i)
		ro, err := u.checkRO(rec, comparator)
		if err != nil {
			continue
		}

		res = append(res, ro)
	}

	result.Result = utils.Dump(res)
	result.TotalBase = int64(len(base))
	result.TotalRO = int64(len(res))
	if err := u.resultRepo.Update(ctx, result); err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("finish process id: ", id)

	return nil
}

func (u *roUsecase) CheckIsDealerExists(ctx context.Context, name string) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":  utils.DumpIncomingContext(ctx),
		"name": name,
	})

	err := u.recordRepo.CheckIsDealerExists(ctx, name)
	switch err {
	default:
		logger.Error(err)
		return ErrInternal
	case repository.ErrNotFound:
		return ErrNotFound
	case nil:
		return nil
	}
}

func (u *roUsecase) checkRO(baseRO model.Record, comparator []model.Record) (model.AdvanceROCheckingResult, error) {
	result := &model.AdvanceROCheckingResult{
		KTP: baseRO.Ktp,
	}

	for _, rec := range comparator {
		if rec.Ktp == baseRO.Ktp {
			result.AddPayment(rec.GetPaymentType())
			result.AddDealer(rec.Dealer)
		}
	}

	if !result.IsRO() {
		return *result, errors.New("bukan RO")
	}

	result.CountNumOrders()

	return *result, nil
}
