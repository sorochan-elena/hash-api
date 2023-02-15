package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"hash-api/internal/domain"
	"hash-api/proto/gen/hash/schema"
	"log"
	"net"
	"time"
)

type hashGetter interface {
	Handle(ctx context.Context) (domain.Hash, error)
}

type Server struct {
	addr       string
	server     *grpc.Server
	hashGetter hashGetter
}

func NewServer(addr string, hashGetter hashGetter, opts ...grpc.ServerOption) *Server {
	var srv = &Server{
		addr:       addr,
		server:     grpc.NewServer(opts...),
		hashGetter: hashGetter,
	}

	schema.RegisterHashApiServer(srv.server, srv)

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	go func() {
		if err = s.server.Serve(lis); err != nil {
			log.Printf("gRPC server Serve error: %s", err)
		}
	}()

	<-ctx.Done()
	s.server.GracefulStop()

	return nil
}

func (s *Server) Get(ctx context.Context, _ *emptypb.Empty) (*schema.Hash, error) {
	hash, err := s.hashGetter.Handle(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	header := metadata.Pairs("expires", hash.ExpiresAt().Format(time.RFC1123))

	if err = grpc.SendHeader(ctx, header); err != nil {
		log.Printf("send header: %s", err)
	}

	return &schema.Hash{
		Hash:      hash.Content,
		CreatedAt: timestamppb.New(hash.CreatedAt),
	}, nil
}
