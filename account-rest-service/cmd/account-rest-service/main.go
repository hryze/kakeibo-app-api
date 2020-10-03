package main

import (
	"fmt"
	"os"

	"github.com/paypay3/kakeibo-app-api/account-rest-service/router"
)

func main() {
	if err := router.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
