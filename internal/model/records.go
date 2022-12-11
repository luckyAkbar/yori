package model

import (
	"context"
	"strings"
	"time"
)

type Record struct {
	ID             int64     `json:"id"`
	No             int64     `json:"no"`
	Nama           string    `json:"nama"`
	NoEngine       string    `json:"no_engine"`
	TglMohonFaktur time.Time `json:"tgl_mohon_faktur"`
	Fincoy         string    `json:"fincoy"`
	Type           string    `json:"type"`
	Ktp            string    `json:"ktp"`
	Kk             string    `json:"kk"`
	Dealer         string    `json:"dealer"`
	Bulan          int       `json:"bulan"`
	Tahun          int       `json:"tahun"`
}

func (r Record) TableName() string {
	return "new_records"
}

func (r *Record) GetPaymentType() string {
	fincoy := strings.Trim(r.Fincoy, " ")
	if fincoy == "" || fincoy == "-" {
		return "CASH"
	}

	return "KREDIT"
}

type RecordUsecase interface {
	FindByKTP(ctx context.Context, ktp string) ([]Record, error)
	FindByKK(ctx context.Context, kk string) ([]Record, error)
}

type RecordRepository interface {
	SaveBulk(ctx context.Context, bulk []Record) error
	FindByKTP(ctx context.Context, ktp string) ([]Record, error)
	FindByKK(ctx context.Context, kk string) ([]Record, error)
	FindByDateRange(ctx context.Context, start, finish time.Time) ([]Record, error)
	SaveBulkWithTransaction(ctx context.Context, bulk []Record) error
	CheckIsDealerExists(ctx context.Context, name string) error
}
