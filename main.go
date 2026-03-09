package main

import (
	"os"

	"github.com/chinmaymk/acli/cmd/acli"
)

func main() {
	if err := acli.Execute(); err != nil {
		os.Exit(1)
	}
}
