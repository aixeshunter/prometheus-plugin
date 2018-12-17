package main

import (
	"fmt"
	"os"

	"github.com/aixeshunter/prometheus-plugin/cmd/app"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
