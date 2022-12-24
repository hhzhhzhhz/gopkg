package go_pool

import (
	"github.com/panjf2000/ants/v2"
)

type GoPool struct {
	*ants.Pool
}

func NewPool(size int, options ...ants.Option) (*GoPool, error) {
	pool, err := ants.NewPool(size, options...)
	if err != nil {
		return &GoPool{}, nil
	}
	return &GoPool{Pool: pool}, err
}
