package flag

import "github.com/urfave/cli/v2"

const HttpAddrFlag = "http-addr"

func HttpAddr() cli.Flag {
	return &cli.StringFlag{
		Name:    HttpAddrFlag,
		Value:   ":8080",
		EnvVars: []string{"HTTP_ADDR"},
	}
}
