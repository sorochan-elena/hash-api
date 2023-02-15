package service

import (
	"context"
	"fmt"
	"hash-api/internal/domain"
	"time"
)

type HashGenerator struct {
	repo     domain.HashKeeper
	ttl      time.Duration
	hashFunc domain.HashFunc
}

func NewHashGenerator(repo domain.HashKeeper, hashFunc domain.HashFunc, ttl time.Duration) *HashGenerator {
	return &HashGenerator{
		repo:     repo,
		hashFunc: hashFunc,
		ttl:      ttl,
	}
}

func (g HashGenerator) Handle(ctx context.Context) error {
	// initialize first hash
	_ = g.repo.Store(ctx, g.hashFunc(time.Now().UTC(), g.ttl))

	// add ticker to handle ttl
	ticker := time.NewTicker(g.ttl)
	defer ticker.Stop()

	for {
		select {
		case tm := <-ticker.C:
			if err := g.repo.Store(ctx, g.hashFunc(tm.UTC(), g.ttl)); err != nil {
				return fmt.Errorf("store: %w", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
