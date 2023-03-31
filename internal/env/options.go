package env

// @lint-ignore ST1003.

import (
	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

func NewOptionsFromEnv() *Options {
	o := Options{}

	return o.loadFromEnv()
}

type Options struct {
	YaJWTkeyID            string `env:"YA_JWT_KEY_ID"`
	YaJWTserviceAccountID string `env:"YA_JWT_SERVICE_ACCOUNT_ID"`
	YaJWTkey              string `env:"YA_JWT_KEY_BASE64"`
	GithubToken           string `env:"GITHUB_TOKEN"`
	GithubOwner           string `env:"GITHUB_OWNER"`
	GithubRepo            string `env:"GITHUB_REPO"`
	GithubSecretName      string `env:"GITHUB_SECRET_NAME"`
}

func (o *Options) loadFromEnv() *Options {
	options := Options{}
	if err := env.Parse(&options); err != nil {
		panic(err)
	}

	logrus.Infof("Options: %+v", options)

	return &options
}
