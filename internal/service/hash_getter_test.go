package service_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"hash-api/internal/domain"
	domainMock "hash-api/internal/domain/mock"
	"hash-api/internal/service"
	serviceMock "hash-api/internal/service/mock"
	"testing"
	"time"
)

func TestHashGetter_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)

	var testHash = domain.NewHash("aaa", time.Now(), time.Now())

	tests := []struct {
		name       string
		mockFields func(ctx context.Context, repo *domainMock.MockHashGetter, cache *serviceMock.Mockcache)
		wantHash   domain.Hash
		wantErr    error
	}{
		{
			name: "ok / from cache",
			mockFields: func(ctx context.Context, repo *domainMock.MockHashGetter, cache *serviceMock.Mockcache) {
				cache.EXPECT().Get(ctx).Return(testHash, nil)
			},
			wantHash: testHash,
		},
		{
			name: "ok / from api",
			mockFields: func(ctx context.Context, repo *domainMock.MockHashGetter, cache *serviceMock.Mockcache) {
				cache.EXPECT().Get(ctx).Return(domain.Hash{}, domain.ErrExpired)
				repo.EXPECT().Get(ctx).Return(testHash, nil)
				cache.EXPECT().Store(ctx, testHash).Return(nil)
			},
			wantHash: testHash,
		},
		{
			name: "err / from api",
			mockFields: func(ctx context.Context, repo *domainMock.MockHashGetter, cache *serviceMock.Mockcache) {
				cache.EXPECT().Get(ctx).Return(domain.Hash{}, domain.ErrExpired)
				repo.EXPECT().Get(ctx).Return(testHash, domain.ErrNilHash)
			},
			wantErr: domain.ErrNilHash,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				ctx     = context.Background()
				cache   = serviceMock.NewMockcache(ctrl)
				repo    = domainMock.NewMockHashGetter(ctrl)
				handler = service.NewHashGetter(repo, cache)
			)

			tt.mockFields(ctx, repo, cache)

			hash, err := handler.Handle(ctx)
			assert.Equal(t, tt.wantHash, hash)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
