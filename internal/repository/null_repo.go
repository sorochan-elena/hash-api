package repository

import (
	"context"
	"hash-api/internal/domain"
)

type NullRepo struct{}

func NewNullRepository() *NullRepo {
	return &NullRepo{}
}

func (n NullRepo) Get(_ context.Context) (domain.Hash, error) {
	return domain.Hash{}, domain.ErrNilHash
}

func (n NullRepo) Store(_ context.Context, _ domain.Hash) error {
	return nil
}
