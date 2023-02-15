package service_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"hash-api/internal/domain"
	domainMock "hash-api/internal/domain/mock"
	"hash-api/internal/service"
	"testing"
	"time"
)

var testHashFunc = func(tm time.Time, ttl time.Duration) domain.Hash {
	return domain.NewHash(uuid.NewString(), tm, tm.Add(ttl))
}

func TestHashGenerator_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	const ttl = time.Millisecond

	repo := domainMock.NewMockHashKeeper(ctrl)
	repo.EXPECT().
		Store(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, hash domain.Hash) {
			assert.Equal(t, time.Millisecond, hash.ExpiresAt().Sub(hash.CreatedAt))
		}).
		Times(3).
		Return(nil)

	go func() {
		time.Sleep(time.Millisecond * 3)
		cancelFunc()
	}()

	handler := service.NewHashGenerator(repo, testHashFunc, ttl)
	_ = handler.Handle(ctx)
}
