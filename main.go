// Package main lints paths in precious config files
package main

import (
	"log"
	"os"

	"github.com/oalders/ppath/audit"
)

func main() {
	// Remove timestamps from logging
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	const requiredArgs = 2
	if len(os.Args) != requiredArgs {
		log.Fatal("Usage: ppath precious.toml")
	}

	filename := os.Args[1]
	config, err := audit.PreciousConfig(filename)
	if err != nil {
		log.Fatalf("Error parsing %s %v", filename, err)
	}

	success, err := audit.Paths(config)
	if err != nil {
		log.Fatal(err)
	}

	if success {
		log.Print("All paths OK")
		os.Exit(0)
	}

	os.Exit(1)
}
