# PTP Configuration Parser - Build Guide

## Quick Start

### Using the Makefile

```bash
# Show all available commands
make help

# Build the application
make build

# Run tests
make test

# Run the application with bidirectional example
make run

# Complete development workflow
make dev
```

### Basic Commands

#### Development
- `make build` - Build the application
- `make test` - Run all tests
- `make run` - Run with bidirectional example
- `make clean` - Clean build artifacts

#### Code Quality
- `make fmt` - Format Go code
- `make vet` - Run go vet
- `make check` - Run format, vet, and test

#### Multi-platform Builds
- `make build-linux` - Build for Linux
- `make build-windows` - Build for Windows  
- `make build-darwin` - Build for macOS
- `make build-all` - Build for all platforms

#### Dependencies
- `make deps` - Download dependencies
- `make deps-update` - Update dependencies

#### Workflows
- `make dev` - Complete development workflow (clean, format, check, build)
- `make release` - Release workflow (clean, check, test, build all platforms)

### Usage Examples

#### Build and run:
```bash
make build
./bin/ptp-config-parser examples/bidirectional.yaml
```

#### Run with different examples:
```bash
make run-example EXAMPLE=dual-wpc.yaml
```

#### Get version:
```bash
./bin/ptp-config-parser --version
```

#### Run tests with coverage:
```bash
make test-coverage
```

#### Validate all examples:
```bash
make validate-examples
```

### Docker Usage

```bash
# Build Docker image
make docker-build

# Run in Docker
make docker-run
```

### Project Information

```bash
# Show project info
make info

# Show version
make version
```

## Manual Build (without Makefile)

If you prefer not to use the Makefile:

```bash
# Build
go build -o bin/ptp-config-parser main.go types.go

# Test  
go test -v ./...

# Run
go run main.go types.go examples/bidirectional.yaml
```