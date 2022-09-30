package model

import "context"

type Record struct {
	No             int64  `json:"no"`
	Nama           string `json:"nama"`
	NoEngine       string `json:"no_engine"`
	TglMohonFaktur string `json:"tgl_mohon_faktur"`
	Fincoy         string `json:"fincoy"`
	Type           string `json:"type"`
	Ktp            string `json:"ktp"`
	Kk             string `json:"kk"`
	Dealer         string `json:"dealer"`
	Bulan          int    `json:"bulan"`
	Tahun          int    `json:"tahun"`
}

type RecordUsecase interface {
	FindByKTP(ctx context.Context, ktp string) ([]Record, error)
	FindByKK(ctx context.Context, kk string) ([]Record, error)
}

type RecordRepository interface {
	FindByKTP(ctx context.Context, ktp string) ([]Record, error)
	FindByKK(ctx context.Context, kk string) ([]Record, error)
}
