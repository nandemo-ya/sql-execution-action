package main

import (
	"log"

	"github.com/nandemo-ya/sql-execution-action/pkg/cmd"
)

func main() {
	if err := cmd.CreateRoot().Execute(); err != nil {
		log.Fatal(err)
	}
}
