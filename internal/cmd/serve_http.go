package cmd

import (
	"github.com/urfave/cli/v2"
	"hash-api/internal/cmd/flag"
	"hash-api/internal/config"
	"hash-api/internal/di"
	"os"
	"os/signal"
)

func ServeHttp() *cli.Command {
	return &cli.Command{
		Name:  "serve-http",
		Usage: "start http server",
		Flags: []cli.Flag{
			flag.HttpAddr(),
			flag.GrpcAddr(),
		},
		Action: func(cliCtx *cli.Context) error {
			ctx, cancelFunc := signal.NotifyContext(cliCtx.Context, os.Interrupt, os.Kill)
			defer cancelFunc()

			cfg := config.FromContext(cliCtx)
			server := di.MakeHttpServer(cfg)

			server.Start(ctx)

			<-ctx.Done()

			return ctx.Err()
		},
	}
}
