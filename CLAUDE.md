# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a CLI-based domain analysis tool written in Go that checks domain availability across traditional DNS and blockchain-based naming systems (ENS, Unstoppable Domains). The tool also provides WHOIS data, blockchain metadata, and domain value estimation.

## Development Commands

```bash
# Build the application
go build -o d3-domain-tool

# Run tests
go test ./...

# Format code
go fmt ./...

# Static analysis
go vet ./...

# Run the tool
./d3-domain-tool -domain=example.com
./d3-domain-tool -domain=mydomain.eth -format=json
```

## Architecture

The project follows a clean, modular architecture:

```
├── main.go                    # CLI interface and entry point
├── internal/
│   ├── analyzer/             # Core analysis orchestration
│   ├── checker/              # DNS availability checking
│   ├── whois/                # WHOIS data retrieval and parsing
│   ├── blockchain/           # Blockchain domain resolution (ENS, Unstoppable)
│   ├── valuation/            # Domain value estimation engine
│   └── output/               # Output formatting (table/JSON)
```

### Key Components

- **Analyzer**: Orchestrates all analysis components and returns unified results
- **DNS Checker**: Uses Go's net package to check traditional DNS records (A, MX, NS, TXT)
- **WHOIS Client**: Raw socket connections to WHOIS servers with parsing logic
- **Blockchain Checker**: Simulates ENS and Unstoppable Domains lookups (extensible for real blockchain integration)
- **Valuation Engine**: Multi-factor domain value estimation based on length, character quality, brandability, etc.
- **Output Formatter**: Clean table and JSON output formats

## Supported Domain Types

- **Traditional DNS**: .com, .net, .org, .io, .co, etc.
- **ENS**: .eth domains
- **Unstoppable Domains**: .crypto, .nft, .x, .wallet, .bitcoin, .dao, .888, .zil, .blockchain

## Usage Examples

```bash
# Traditional domain analysis with full WHOIS data
./d3-domain-tool -domain=google.com

# Blockchain domain check
./d3-domain-tool -domain=vitalik.eth

# JSON output for programmatic use
./d3-domain-tool -domain=example.org -format=json
```

## Testing

Tests are included for the valuation engine. Run with `go test ./...`. The project is designed to be easily testable with dependency injection patterns.