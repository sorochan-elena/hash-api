package di

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"hash-api/internal/config"
	"hash-api/internal/domain"
	grpcPort "hash-api/internal/ports/grpc"
	httpPort "hash-api/internal/ports/http"
	"hash-api/internal/repository"
	"hash-api/internal/service"
	"hash-api/proto/gen/hash/schema"
	"log"
	"time"
)

func MakeHttpServer(cfg config.App) *httpPort.Server {
	return httpPort.NewServer(cfg.HttpAddr, service.NewHashGetter(hashApiRepo(cfg), hashInMemoryRepo()))
}

func MakeGrpcServer(cfg config.App) (*grpcPort.Server, *service.HashGenerator) {
	repo := hashInMemoryRepo()
	grpcServer := grpcPort.NewServer(cfg.GrpcAddr, service.NewHashGetter(repo, repository.NewNullRepository()))
	hashGenerator := service.NewHashGenerator(repo, hashFunc(), cfg.HashTtl)

	return grpcServer, hashGenerator
}

func hashInMemoryRepo() *repository.HashInMemoryRepository {
	return repository.NewHashInMemoryRepository()
}

func hashApiRepo(cfg config.App) *repository.HashApiRepository {
	conn, err := grpc.DialContext(context.Background(), cfg.GrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial: %s", err)
	}

	return repository.NewHashApiRepository(schema.NewHashApiClient(conn))
}

func hashFunc() domain.HashFunc {
	return func(tm time.Time, ttl time.Duration) domain.Hash {
		return domain.NewHash(uuid.NewString(), tm, tm.Add(ttl))
	}
}
