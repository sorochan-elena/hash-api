package repository

import (
	"context"
	"hash-api/internal/domain"
	"sync"
)

type HashInMemoryRepository struct {
	hash *domain.Hash
	mx   sync.RWMutex
}

func NewHashInMemoryRepository() *HashInMemoryRepository {
	return &HashInMemoryRepository{}
}

func (h *HashInMemoryRepository) Get(ctx context.Context) (domain.Hash, error) {
	select {
	case <-ctx.Done():
		return domain.Hash{}, ctx.Err()
	default:
	}

	h.mx.RLock()
	defer h.mx.RUnlock()

	if h.hash == nil {
		return domain.Hash{}, domain.ErrNilHash
	}
	if h.hash.Expired() {
		return domain.Hash{}, domain.ErrExpired
	}

	return *h.hash, nil
}

func (h *HashInMemoryRepository) Store(ctx context.Context, hash domain.Hash) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	h.mx.Lock()
	defer h.mx.Unlock()

	h.hash = &hash
	return nil
}
