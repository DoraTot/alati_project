package model

type IdempotencyRequest struct {
	Key string `json:"key"`
}

func (i *IdempotencyRequest) SetKey(key string) {
	i.Key = key
}

type IdempotencyRepository interface {
	Add(i *IdempotencyRequest) error
	Get(key string) (bool, error)
}
