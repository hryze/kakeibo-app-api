package main

import (
	"fmt"
	"os"

	"github.com/hryze/kakeibo-app-api/todo-rest-service/router"
)

func main() {
	if err := router.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
