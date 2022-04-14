package databases

import (
	"context"
	"encoding/json"

	rediscli "github.com/go-redis/redis/v8"
)

const key string = "fibonacci"

type Redis interface {
	GetMaxFibonacci() (*uint64, error)
	GetSavedFibonacci() ([]uint64, error)
	SaveFibonacci([]uint64) error
}

type redis struct {
	client *rediscli.Client
	ctx    context.Context
}

func NewRedis(address string, ctx context.Context) *redis {
	rdb := rediscli.NewClient(&rediscli.Options{
		Addr: address,
		DB:   0,
	})

	return &redis{
		client: rdb,
		ctx:    ctx,
	}
}

func (r *redis) GetMaxFibonacci() (*uint64, error) {
	res, err := r.GetSavedFibonacci()
	if err != nil {
		return nil, err
	}

	length := len(res) - 1
	if length < 0 {
		length = 0
	}

	retVal := uint64(length)
	return &retVal, nil
}

func (r *redis) GetSavedFibonacci() ([]uint64, error) {
	raw, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var res []uint64
	err = json.Unmarshal([]byte(raw), &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *redis) SaveFibonacci(sequence []uint64) error {
	val, err := json.Marshal(sequence)
	if err != nil {
		return err
	}

	err = r.client.Set(r.ctx, key, val, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
