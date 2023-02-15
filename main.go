package main

import (
	"github.com/urfave/cli/v2"
	"hash-api/internal/cmd"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "hash",
		Usage: "retrieve uuid hash with ttl",
		Commands: []*cli.Command{
			cmd.ServeGrpc(),
			cmd.ServeHttp(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
