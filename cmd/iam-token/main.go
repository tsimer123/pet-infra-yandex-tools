package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tsimer123/pet-infra-yandex-tools/internal/jwt"
	"github.com/tsimer123/pet-infra-yandex-tools/internal/options"
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

	o := options.NewOptionsFromEnv()
	t := jwt.NewJWT(o)
	fmt.Print(t.GetIAMToken())
}
