package model

import (
	"context"
	"fmt"
	"time"
)

type AdvanceROCheckingInput struct {
	TipeRO        string    `json:"tipe_ro"`
	BaseStart     time.Time `json:"base_start"`
	BaseFinish    time.Time `json:"base_finish"`
	CompareStart  time.Time `json:"compare_start"`
	CompareFinish time.Time `json:"compare_finish"`
}

type AdvanceROCheckingResult struct {
	KTP         string   `json:"ktp"`
	Payments    []string `json:"payments"`
	DealerNames []string `json:"dealer_names"`
	Tahun       []string `json:"tahun"`
	NumOrders   int      `json:"num_orders"`
}

func (arocr *AdvanceROCheckingResult) AddPayment(s string) {
	arocr.Payments = append(arocr.Payments, s)
}

func (arocr *AdvanceROCheckingResult) AddTahun(t int) {
	arocr.Tahun = append(arocr.Tahun, fmt.Sprintf("%d", t))
}

func (arocr *AdvanceROCheckingResult) AddDealer(s string) {
	arocr.DealerNames = append(arocr.DealerNames, s)
}

func (arocr *AdvanceROCheckingResult) CountNumOrders() {
	arocr.NumOrders = len(arocr.Payments)
}

func (arocr *AdvanceROCheckingResult) IsRO() bool {
	return len(arocr.DealerNames) > 0
}

type ROUsecase interface {
	GetResult(ctx context.Context, id int64) ([]AdvanceROCheckingResult, *Result, error)
	HandleAdvanceChecking(ctx context.Context, input *AdvanceROCheckingInput, id int64) error
	CheckIsDealerExists(ctx context.Context, name string) error
}
