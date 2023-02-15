package flag

import (
	"github.com/urfave/cli/v2"
	"time"
)

const (
	GrpcAddrFlag = "grpc-addr"
	HashTtlFlag  = "hash-ttl"
)

func GrpcAddr() cli.Flag {
	return &cli.StringFlag{
		Name:    GrpcAddrFlag,
		Value:   ":5105",
		EnvVars: []string{"GRPC_ADDR"},
	}
}

const defaultHashTtl = time.Minute * 5

func HashTtl() cli.Flag {
	return &cli.DurationFlag{
		Name:    HashTtlFlag,
		Value:   defaultHashTtl,
		EnvVars: []string{"HASH_TTL"},
	}
}
