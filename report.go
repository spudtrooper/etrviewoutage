package main

import (
	"flag"
	"log"

	"github.com/spudtrooper/etrviewoutage/lib"
)

func main() {
	var (
		dataDir = flag.String("data_dir", "data", "Output directory")
	)
	flag.Parse()
	if err := lib.OutputAll(*dataDir); err != nil {
		log.Fatalf("OutputAll: %v", err)
	}
}
