package cmd_test

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"hash-api/internal/cmd"
	"hash-api/internal/domain"
	"io"
	"net/http"
	"testing"
	"time"
)

const httpAddr = ":8084"

func TestHttpServerAccessible(t *testing.T) {
	cancelGrpcCmdFunc, grpcConn, err := serveGrpcServer(context.Background(), grpcAddr, hashTtl.String())
	if err != nil {
		t.Fail()
		return
	}
	defer cancelGrpcCmdFunc()
	defer grpcConn.Close()

	app := &cli.App{Writer: io.Discard}

	flags := []string{"serve-http", "--http-addr", httpAddr, "--grpc-addr", grpcAddr}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse(flags)

	var cliCtx = cli.NewContext(app, set, nil)
	ctx, cancelFunc := context.WithCancel(context.Background())
	cliCtx.Context = ctx
	defer cancelFunc()

	// run command
	go func() {
		command := cmd.ServeHttp()
		err = command.Run(cliCtx, flags...)
		if err != nil {
			cancelFunc()
		}
	}()

	// check that server is accessible
	retries := 5

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Millisecond):
			resp, apiErr := http.Get("http://localhost" + httpAddr + "/hash")
			if apiErr != nil {
				retries--
			}

			if retries == 0 {
				t.Fatalf("api: %s", err)
				return
			}

			if apiErr == nil {
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				body, _ := io.ReadAll(resp.Body)
				var h domain.Hash
				_ = json.Unmarshal(body, &h)
				assert.NotEmpty(t, h.Content)
				assert.NotEmpty(t, h.CreatedAt)
				return
			}
		}
	}
}
