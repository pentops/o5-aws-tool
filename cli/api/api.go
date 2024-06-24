package api

import "github.com/pentops/o5-aws-tool/libo5"

type BaseCommand struct {
	API string `env:"O5_API" flag:"api"`
}

func (cfg *BaseCommand) Client() *libo5.API {
	client := libo5.NewAPI(cfg.API)
	return client
}
