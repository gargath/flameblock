package main

import (
	"fmt"
	"os"

	"github.com/gargath/flameblock/pkg/collector"
)

func main() {
	fmt.Printf("Flameblock Collector version %s\n", VERSION)

	c := &collector.Server{}

	err := c.Start()
	if err != nil {
		fmt.Printf("Error starting server: %s", err)
		os.Exit(1)
	}
}
