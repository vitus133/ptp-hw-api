package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Version can be set during build time
var Version = "dev"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("PTP Configuration Parser")
		fmt.Printf("Version: %s\n", Version)
		fmt.Println("Usage: go run . <config-file>")
		fmt.Println("       go run . --version")
		os.Exit(1)
	}

	// Handle version flag
	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		fmt.Printf("ptp-config-parser v%s\n", Version)
		os.Exit(0)
	}

	configFile := os.Args[1]

	// Read file
	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Parse YAML
	var config ClockChain
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("Error parsing YAML: %v\n", err)
		os.Exit(1)
	}

	// Validate
	if err := config.Validate(); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		os.Exit(1)
	}

	// Print result
	fmt.Printf("Successfully parsed and validated: %s\n", configFile)
	fmt.Printf("%s\n", config.String())

	// Print subsystems
	for i, subsystem := range config.Structure {
		fmt.Printf("  %d. %s\n", i+1, subsystem.String())
	}
}
