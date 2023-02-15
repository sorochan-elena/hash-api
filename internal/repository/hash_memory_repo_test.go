package repository_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"hash-api/internal/domain"
	"hash-api/internal/repository"
	"testing"
	"time"
)

func TestHashInMemoryRepository_Get(t *testing.T) {
	var (
		testHash    = domain.NewHash("aaa", time.Now(), time.Now().Add(time.Second))
		expiredHash = domain.NewHash("aaa", time.Now().Add(time.Second), time.Now())
	)

	tests := []struct {
		name     string
		repo     *repository.HashInMemoryRepository
		wantErr  error
		wantHash domain.Hash
	}{
		{
			name:    "nil hash",
			repo:    repository.NewHashInMemoryRepository(),
			wantErr: domain.ErrNilHash,
		},
		{
			name: "expired",
			repo: func() *repository.HashInMemoryRepository {
				r := repository.NewHashInMemoryRepository()
				_ = r.Store(context.Background(), expiredHash)
				return r
			}(),
			wantErr: domain.ErrExpired,
		},
		{
			name: "non-nil hash",
			repo: func() *repository.HashInMemoryRepository {
				r := repository.NewHashInMemoryRepository()
				_ = r.Store(context.Background(), testHash)
				return r
			}(),
			wantHash: testHash,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := tt.repo.Get(context.Background())
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantHash, hash)
		})
	}
}
