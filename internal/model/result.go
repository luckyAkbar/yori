package model

import "context"

type Result struct {
	ID        int64  `json:"id"`
	Result    string `json:"result"`
	TotalBase int64  `json:"total_base"`
	TotalRO   int64  `json:"total_ro"`
}

type ResultRepository interface {
	Create(ctx context.Context, result *Result) error
	Update(ctx context.Context, result *Result) error
	FindByID(ctx context.Context, id int64) (*Result, error)
}
