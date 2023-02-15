package repository

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"hash-api/internal/domain"
	"hash-api/proto/gen/hash/schema"
	"log"
	"time"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE
type apiClient interface {
	Get(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*schema.Hash, error)
}

type HashApiRepository struct {
	client apiClient
}

func NewHashApiRepository(client apiClient) *HashApiRepository {
	return &HashApiRepository{client: client}
}

func (h *HashApiRepository) Get(ctx context.Context) (domain.Hash, error) {
	var header metadata.MD
	hash, err := h.client.Get(ctx, &emptypb.Empty{}, grpc.Header(&header))
	if err != nil {
		return domain.Hash{}, fmt.Errorf("client: %w", err)
	}

	expHeader := header.Get("expires")
	var expiresAt time.Time
	if len(expHeader) > 0 {
		expiresAt, err = time.Parse(time.RFC1123, expHeader[0])
		if err != nil {
			log.Printf("failed to get expire time from header: %s", err)
		}
	}

	// fallback to createdAt time in case if we had error while reading [expires] header
	if expiresAt.IsZero() {
		expiresAt = hash.CreatedAt.AsTime()
	}

	return domain.NewHash(hash.Hash, hash.CreatedAt.AsTime(), expiresAt), nil
}
