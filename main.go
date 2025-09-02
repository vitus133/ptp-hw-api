package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Version can be set during build time
var Version = "dev"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("PTP Hardware Configuration Parser")
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

	// Resolve clock aliases early
	if err := config.ResolveClockAliases(); err != nil {
		fmt.Printf("Error resolving clock aliases: %v\n", err)
		os.Exit(1)
	}

	// Load hardware plugins and apply defaults
	pluginManager, err := NewPluginManager("plugins")
	if err != nil {
		fmt.Printf("Warning: Failed to load plugins: %v\n", err)
		fmt.Println("Continuing without plugin defaults...")
	} else {
		fmt.Printf("Loaded %d hardware plugins: %v\n", len(pluginManager.ListPlugins()), pluginManager.ListPlugins())

		// Apply plugin defaults to user configuration
		if err := pluginManager.MergeUserConfigWithDefaults(&config); err != nil {
			fmt.Printf("Error applying plugin defaults: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Successfully applied hardware plugin defaults")

		// Output the merged configuration
		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Println("MERGED CONFIGURATION (User Config + Plugin Defaults)")
		fmt.Println(strings.Repeat("=", 60))

		mergedYAML, err := yaml.Marshal(&config)
		if err != nil {
			fmt.Printf("Warning: Failed to marshal merged config: %v\n", err)
		} else {
			fmt.Printf("%s\n", string(mergedYAML))
		}
		fmt.Println(strings.Repeat("=", 60))
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
