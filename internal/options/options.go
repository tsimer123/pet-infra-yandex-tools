package options

// @lint-ignore ST1003.

import (
	"github.com/caarlos0/env/v6"
)

func NewOptionsFromEnv() *Options {
	o := Options{}

	return o.loadFromEnv()
}

type Options struct {
	GWTkeyID            string `env:"GWT_KEY_ID"`
	GWTserviceAccountID string `env:"GWT_SERVICE_ACCOUNT_ID"`
	GWTkeyFile          string `env:"GWT_KEY_FILE"`
}

func (o *Options) loadFromEnv() *Options {
	options := Options{}
	if err := env.Parse(&options); err != nil {
		panic(err)
	}

	return &options
}
