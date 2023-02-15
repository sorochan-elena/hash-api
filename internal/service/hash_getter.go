package service

import (
	"context"
	"fmt"
	"hash-api/internal/domain"
	"log"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE
type cache interface {
	domain.HashGetter
	domain.HashKeeper
}

type HashGetter struct {
	repo  domain.HashGetter
	cache cache
}

func NewHashGetter(repo domain.HashGetter, cache cache) *HashGetter {
	return &HashGetter{repo: repo, cache: cache}
}

func (g *HashGetter) Handle(ctx context.Context) (domain.Hash, error) {
	// check if we have valid hash in cache
	if hash, err := g.cache.Get(ctx); err == nil {
		return hash, nil
	}

	// retrieve hash from repository
	hash, err := g.repo.Get(ctx)
	if err != nil {
		return domain.Hash{}, fmt.Errorf("get: %w", err)
	}

	// store cache
	if err = g.cache.Store(ctx, hash); err != nil {
		log.Printf("store: %s", err)
	}

	return hash, nil
}
