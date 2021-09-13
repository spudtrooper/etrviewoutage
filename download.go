package main

import (
	"flag"
	"log"
	"time"

	"github.com/spudtrooper/etrviewoutage/lib"
)

func main() {
	var (
		dataDir  = flag.String("data_dir", "data", "Output directory")
		headless = flag.Bool("headless", false, "Selenium headless")
		verbose  = flag.Bool("verbose", false, "Selenium verbose")
		pause    = flag.Duration("pause", 5*time.Second, "Amount to pause to let the page load")
	)
	flag.Parse()
	if err := lib.Download(*dataDir, *verbose, *headless, *pause); err != nil {
		log.Fatalf("Download: %v", err)
	}
}
