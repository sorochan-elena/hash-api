package config

import (
	"github.com/urfave/cli/v2"
	"hash-api/internal/cmd/flag"
	"time"
)

type App struct {
	HttpAddr string
	GrpcAddr string
	HashTtl  time.Duration
}

func FromContext(c *cli.Context) App {
	return App{
		HttpAddr: c.String(flag.HttpAddrFlag),
		GrpcAddr: c.String(flag.GrpcAddrFlag),
		HashTtl:  c.Duration(flag.HashTtlFlag),
	}
}
