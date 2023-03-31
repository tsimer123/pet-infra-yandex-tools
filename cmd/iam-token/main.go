package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tsimer123/pet-infra-yandex-tools/internal/env"
	"github.com/tsimer123/pet-infra-yandex-tools/internal/github"
	"github.com/tsimer123/pet-infra-yandex-tools/internal/jwt"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   "2006-01-02 15:04:05.000",
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		DataKey:           "",
		FieldMap:          nil,
		CallerPrettyfier:  nil,
		PrettyPrint:       false,
	})
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)

	o := env.NewOptionsFromEnv()
	t := jwt.NewJWT(o)
	g := github.NewGithub(o)
	g.UpdateSecret(t.GetIAMToken())
	fmt.Print(t.GetIAMToken())
}
