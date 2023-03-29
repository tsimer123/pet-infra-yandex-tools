package main

import (
	"fmt"

	"github.com/tsimer123/pet-infra-yandex-vm-starter/internal/jwt"
	"github.com/tsimer123/pet-infra-yandex-vm-starter/internal/options"
)

func main() {
	o := options.NewOptionsFromEnv()
	t := jwt.NewJWT(o)
	fmt.Print(t.GetIAMToken())
}
