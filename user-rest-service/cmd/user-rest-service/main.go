package main

import (
	"fmt"
	"os"

	"github.com/hryze/kakeibo-app-api/user-rest-service/infrastructure/router"
)

func main() {
	if err := router.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
