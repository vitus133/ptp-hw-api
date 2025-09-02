# PTP Configuration Parser

A Go-based parser and validator for Precision Time Protocol (PTP) clock chain configurations with hardware-specific plugin support.

## Overview

This tool parses, validates, and merges PTP configuration files written in YAML format. It implements a **base configuration approach** where hardware-specific plugin defaults provide comprehensive pin configurations, and user configurations can overlay specific settings without losing other defaults.

## Key Features

- **🔧 Hardware Plugin System**: Extensible plugin architecture for different hardware platforms
- **📋 Base Configuration Approach**: All hardware defaults included by default, user configs overlay specific changes
- **✅ Configuration Validation**: Comprehensive validation of clock chain configurations
- **🔄 Smart Merging**: Intelligent merging of user configurations with hardware plugin defaults
- **📖 Multiple Configuration Types**: Support for various PTP deployment scenarios

## Supported Hardware

Current hardware plugins include:

- **Intel E810** (`e810.yaml`) - Intel E810 Network Adapter with Time Synchronization (real)
- **Intel GNR-D** (`gnr-d.yaml`) - Intel GNR-D platform (real, but names and capabilities here are fictional)
- **Timestone** (`timestone.yaml`) - Timestone hardware platform (fictional)

## Installation & Usage

### Prerequisites

- Go 1.21 or later
- YAML configuration files

### Building

```bash
# Clone the repository
git clone <repository-url>
cd openAPI-ptp

# Build the binary
make build

# Or run directly with Go
go run . <config-file>
```

### Usage

```bash
# Parse and validate a configuration file
./ptp-config-parser examples/tgm-wpc-single.yaml

# Or run with Go
go run . examples/triple-t-bc-wpc.yaml

# Check version
./ptp-config-parser --version
```

### Output

The tool provides:
- **Validation Results**: Confirms configuration correctness
- **Merged Configuration**: Shows final config with plugin defaults applied
- **System Summary**: Overview of subsystems, sources, and conditions

## Configuration Format

### Basic Structure

```yaml
# Optional: Shared definitions for reuse
commonDefinitions:
  # Optional: Define aliases for clock IDs to simplify configs
  clockIdentifiers:
  - alias: "GM1"
    clockId: "0x112233fffe445566"
    description: "Leader DPLL"
  - alias: "FOL1"
    clockId: "0xc7cc7cfffe001122"
  
  eSyncDefinitions:
  - name: "10MHz-1PPS"
    esyncConfig:
      transferFrequency: 10000000
      embeddedSyncFrequency: 1
      dutyCyclePct: 25

# Required: System structure definition
structure:
- name: "System Name"
  hardwarePlugin: "e810"  # References hardware plugin
  ethernet:
  - ports: ["ens4f0"]
  dpll: 
    clockId: "GM1"  # using alias instead of the raw clock ID
    phaseInputs:
      GNSS_1PPS:
        frequency: 1
        description: "GPS reference"

# Optional: Behavioral rules and conditions
behavior:
  sources:
  - name: "PTP"
    clockId: "GM1"
    sourceType: "ptpTimeReceiver"
    boardLabel: "CVL_SDP22"
  
  conditions:
  - name: "Default Configuration"
    sources:
    - sourceName: "Default on profile (re)load"
      conditionType: "default"
    desiredStates:
    - clockId: "GM1"
      boardLabel: "GNSS_1PPS"
      eec:
        priority: 0
      pps:
        priority: 0
```

### Hardware Plugin Format

Hardware plugins define default configurations for specific devices:

```yaml
# Plugin metadata
pluginInfo:
  name: "e810"
  description: "Intel E810 Network Adapter with Time Synchronization"
  version: "1.0.0"
  vendor: "Intel"

# Default pin configurations
specificDefaults:
  GNSS_1PPS:
    eec:
      priority: 0
    pps:
      priority: 0
  SMA1:
    eec:
      priority: 3
    pps:
      priority: 3
  # ... additional pins
```

## Configuration Examples

The `examples/` directory contains various deployment scenarios:

- **`tgm-wpc-single.yaml`** - Single card WPC Grandmaster
- **`dual-wpc.yaml`** - Dual WPC setup
- **`triple-t-bc-wpc.yaml`** - Triple T-BC WPC configuration
- **`bidirectional.yaml`** - [Bidirectional synchronization setup](examples/bidirectional.png)
- **`dual-wpc-esync.yaml`** - Dual WPC with eSync outputs

## How It Works

### Base Configuration Approach

1. **Plugin Defaults Applied**: All pin configurations from the relevant hardware plugin are included as defaults
2. **User Overlays**: User configuration can override specific settings without losing other plugin defaults
3. **Auto-Generation**: If no default condition exists, one is automatically created with plugin defaults
4. **Smart Merging**: User settings take precedence while preserving comprehensive hardware defaults

### Example Merge Process

```yaml
# User Config (minimal)
structure:
- name: "My System"
  hardwarePlugin: "e810"
  dpll:
    clockId: "0x123456"
    phaseInputs:
      GNSS_1PPS:
        frequency: 1

# Plugin Defaults (comprehensive)
specificDefaults:
  GNSS_1PPS:
    eec: {priority: 0}
    pps: {priority: 0}
  SMA1:
    eec: {priority: 3}
    pps: {priority: 3}
  # ... all other pins

# Merged Result
# - All plugin pins included with defaults
# - User GNSS_1PPS settings preserved
# - Auto-generated default condition created
```

## Validation

The tool performs comprehensive validation including:

- **Structure Validation**: Ensures required fields are present
- **Clock ID Format**: Validates clock ID format and uniqueness
- **Hardware Plugin**: Verifies plugin existence and compatibility
- **Pin Configurations**: Validates pin settings and states
- **Behavioral Logic**: Checks source conditions and state consistency

## Development

### Project Structure

```
├── main.go              # CLI entry point
├── types.go             # Configuration data structures
├── plugin_manager.go    # Hardware plugin system
├── api_test.go          # Tests
├── examples/            # Example configurations
│   ├── tgm-wpc-single.yaml
│   ├── triple-t-bc-wpc.yaml
│   └── ...
├── plugins/             # Hardware plugin definitions
│   ├── e810.yaml
│   ├── gnr-d.yaml
│   └── timestone.yaml
└── Makefile            # Build automation
```

### Adding New Hardware Plugins

1. Create a new YAML file in `plugins/` directory
2. Define `pluginInfo` with name, description, version, vendor
3. Add `specificDefaults` with pin configurations
4. Reference the plugin name in user configurations

### Dependencies

- `gopkg.in/yaml.v3` - YAML parsing and marshaling

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

[Add your license information here]

## Related Technologies

- **PTP (IEEE 1588)** - Precision Time Protocol for network time synchronization
- **DPLL** - Digital Phase-Locked Loop for clock synchronization
- **eSync** - Embedded Synchronization for timing distribution
- **Network Interface Cards** - Hardware time synchronization capabilities

## Support

For issues and questions:
- Check the `examples/` directory for configuration patterns
- Review plugin documentation in `plugins/` directory
- Submit issues for bugs or feature requests