package mutex

import (
	"context"
	"fmt"
	client3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

const (
	ttl = 5
)

type Mutex interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
	Close() error
}

func NewMutex(cli *client3.Client, pfx string) (Mutex, error) {
	if pfx == "" {
		return nil, fmt.Errorf("pfx is empty")
	}
	sin, err := concurrency.NewSession(cli, concurrency.WithTTL(ttl))
	if err != nil {
		return nil, err
	}
	return &mutex{sin: sin, mux: concurrency.NewMutex(sin, pfx)}, nil
}

type mutex struct {
	sin *concurrency.Session
	mux *concurrency.Mutex
}

func (m *mutex) Lock(ctx context.Context) error {
	if err := m.mux.Lock(ctx); err != nil {
		return err
	}
	return nil
}

func (m *mutex) Unlock(ctx context.Context) error {
	if err := m.mux.Unlock(ctx); err != nil {
		return err
	}
	return nil
}

func (m *mutex) Close() error {
	return m.sin.Close()
}
