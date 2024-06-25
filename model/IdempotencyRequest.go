package model

import "context"

type IdempotencyRequest struct {
	Key string `json:"key"`
}

func (i *IdempotencyRequest) SetKey(key string) {
	i.Key = key
}

type IdempotencyRepository interface {
	Add(i *IdempotencyRequest, ctx context.Context) error
	Get(key string, ctx context.Context) (bool, error)
}
