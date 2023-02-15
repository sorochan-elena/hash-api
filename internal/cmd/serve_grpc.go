package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"hash-api/internal/cmd/flag"
	"hash-api/internal/config"
	"hash-api/internal/di"
	"hash-api/internal/service"
	"log"
	"os"
	"os/signal"
)

func ServeGrpc() *cli.Command {
	return &cli.Command{
		Name:  "serve-grpc",
		Usage: "start grpc server",
		Flags: []cli.Flag{
			flag.GrpcAddr(),
			flag.HashTtl(),
		},
		Action: func(cliCtx *cli.Context) error {
			ctx, cancelFunc := signal.NotifyContext(cliCtx.Context, os.Interrupt, os.Kill)

			cfg := config.FromContext(cliCtx)
			server, gen := di.MakeGrpcServer(cfg)

			var genDoneCh = make(chan struct{})

			go func(s *service.HashGenerator) {
				if err := s.Handle(ctx); err != nil {
					log.Println("generator:", err)
				}
				log.Println("generator: stopped")
				close(genDoneCh)
			}(gen)

			log.Println(fmt.Sprintf("starting grpc server %s", cfg.GrpcAddr))
			if err := server.Start(ctx); err != nil {
				cancelFunc()
				log.Println("grpc server:", err)
			}

			<-ctx.Done()
			<-genDoneCh

			return ctx.Err()
		},
	}
}
