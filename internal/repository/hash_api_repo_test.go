package repository_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"hash-api/internal/domain"
	"hash-api/internal/repository"
	"hash-api/internal/repository/mock"
	"hash-api/proto/gen/hash/schema"
	"testing"
	"time"
)

func TestHashApiRepository_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		testHash = domain.NewHash("aaaa", time.Now().UTC(), time.Now().Add(time.Minute).UTC())
		apiErr   = errors.New("api err")
	)

	tests := []struct {
		name       string
		mockClient func(ctx context.Context, apiClient *mock_repository.MockapiClient)
		wantErr    error
		wantHash   domain.Hash
	}{
		{
			name: "ok / api call",
			mockClient: func(ctx context.Context, apiClient *mock_repository.MockapiClient) {
				apiClient.EXPECT().
					Get(ctx, gomock.Any(), gomock.Any()).
					Return(&schema.Hash{Hash: testHash.Content, CreatedAt: timestamppb.New(testHash.CreatedAt)}, nil)
			},
			wantHash: testHash,
		},
		{
			name: "err",
			mockClient: func(ctx context.Context, apiClient *mock_repository.MockapiClient) {
				apiClient.EXPECT().
					Get(ctx, gomock.Any(), gomock.Any()).
					Return(nil, apiErr)
			},
			wantErr: apiErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiClient := mock_repository.NewMockapiClient(ctrl)
			ctx := context.Background()

			tt.mockClient(ctx, apiClient)

			repo := repository.NewHashApiRepository(apiClient)
			hash, err := repo.Get(ctx)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantHash.Content, hash.Content)
			assert.Equal(t, tt.wantHash.CreatedAt, hash.CreatedAt)
		})
	}
}
