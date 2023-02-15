package cmd_test

import (
	"context"
	"errors"
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"hash-api/internal/cmd"
	"hash-api/proto/gen/hash/schema"
	"io"
	"testing"
	"time"
)

const (
	grpcAddr = ":5107"
	hashTtl  = time.Second
)

func TestGrpcServerAccessible(t *testing.T) {
	cancelFunc, grpcConn, err := serveGrpcServer(context.Background(), grpcAddr, hashTtl.String())
	defer cancelFunc()

	if err != nil {
		t.Fail()
		return
	}

	defer grpcConn.Close()

	client := schema.NewHashApiClient(grpcConn)

	var header metadata.MD
	hash, err := client.Get(context.Background(), &emptypb.Empty{}, grpc.Header(&header))
	assert.Nil(t, err)
	assert.NotEmpty(t, hash.Hash)

	// check expires header
	expHeader := header.Get("expires")
	assert.NotEmpty(t, expHeader)
	assert.NotEmpty(t, hash.CreatedAt.AsTime())

	// check that ttl is valid
	expiresAt, err := time.Parse(time.RFC1123, expHeader[0])
	assert.Nil(t, err)
	assert.LessOrEqual(t, expiresAt.Sub(hash.CreatedAt.AsTime()), hashTtl)
}

func serveGrpcServer(ctx context.Context, addr, ttl string) (context.CancelFunc, *grpc.ClientConn, error) {
	app := &cli.App{Writer: io.Discard}

	flags := []string{"serve-grpc", "--grpc-addr", addr, "--hash-ttl", ttl}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse(flags)

	var cliCtx = cli.NewContext(app, set, nil)
	ctx, cancelFunc := context.WithCancel(context.Background())
	cliCtx.Context = ctx

	// run command
	go func() {
		command := cmd.ServeGrpc()
		err := command.Run(cliCtx, flags...)
		if err != nil {
			cancelFunc()
		}
	}()

	// check that server is accessible
	retries := 5

	var (
		conn *grpc.ClientConn
		err  error
	)

	for {
		select {
		case <-ctx.Done():
			return cancelFunc, conn, errors.New("context cancelled")
		case <-time.After(time.Millisecond):
			conn, err = grpc.DialContext(context.Background(), addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				retries--
			}

			if retries == 0 {
				return cancelFunc, conn, errors.New("failed to reach server")
			}

			if err == nil {
				return cancelFunc, conn, nil
			}
		}
	}
}
