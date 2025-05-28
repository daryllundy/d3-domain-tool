# D3 Domain Analysis Tool

A CLI-based domain analysis tool that checks domain availability across traditional DNS and blockchain-based naming systems, provides WHOIS data, and estimates domain values.

## Features

- ✅ **Domain Availability**: Check availability across traditional DNS (.com, .net, etc.) and blockchain domains (.eth, .crypto, etc.)
- 🔶 **DOMA Protocol Integration**: Check tokenization status, DeFi usage, and cross-chain presence
- 🔍 **WHOIS Data**: Retrieve and display comprehensive WHOIS information
- ⛓️ **Blockchain Support**: ENS and Unstoppable Domains integration
- 💰 **Enhanced Domain Valuation**: Intelligent domain value estimation including DomainFi factors
- 📦 **Clean Output**: Beautiful CLI formatting with table and JSON output options

## Installation

```bash
# Clone the repository
git clone https://github.com/daryllundy/d3-domain-tool
cd d3-domain-tool

# Build the application
go build -o d3-domain-tool

# Run the tool
./d3-domain-tool -domain=example.com
```

## Usage

### Basic Usage

```bash
# Check a traditional domain
./d3-domain-tool -domain=example.com

# Check a blockchain domain
./d3-domain-tool -domain=mydomain.eth

# Get JSON output
./d3-domain-tool -domain=example.com -format=json
```

### Command Line Options

- `-domain`: Domain to analyze (required)
- `-format`: Output format - `table` (default) or `json`
- `-help`: Show help message

### Examples

```bash
# Traditional domain analysis
./d3-domain-tool -domain=mycompany.com

# ENS domain check
./d3-domain-tool -domain=vitalik.eth

# Unstoppable domain check
./d3-domain-tool -domain=myname.crypto

# JSON output for integration
./d3-domain-tool -domain=example.org -format=json
```

## Supported Domain Types

### Traditional DNS
- .com, .net, .org, .info, .biz, .name
- .io, .co, .me, .tv, .cc, .ws
- Many other TLDs

### Blockchain Domains
- **ENS**: .eth domains
- **Unstoppable Domains**: .crypto, .nft, .x, .wallet, .bitcoin, .dao, .888, .zil, .blockchain

## Output Information

The tool provides comprehensive analysis including:

- **Availability Status**: Whether the domain is available or taken
- **DOMA Protocol Integration**: Tokenization status, token rights, DeFi usage, cross-chain presence
- **WHOIS Data**: Registration details, expiry dates, name servers
- **Blockchain Metadata**: Owner addresses, resolver information, crypto addresses
- **Domain Valuation**: Estimated value with confidence level and reasoning (enhanced with DomainFi factors)
- **Valuation Factors**: Length, character quality, brandability, pronounceability

## Architecture

The application is built with a modular design:

- `main.go`: CLI interface and argument parsing
- `internal/analyzer`: Core analysis orchestration
- `internal/checker`: DNS availability checking
- `internal/whois`: WHOIS data retrieval
- `internal/blockchain`: Blockchain domain resolution
- `internal/doma`: DOMA Protocol integration and tokenization analysis
- `internal/valuation`: Domain value estimation engine
- `internal/output`: Output formatting (table/JSON)

## Development

```bash
# Run tests
go test ./...

# Build for different platforms
GOOS=linux go build -o d3-domain-tool-linux
GOOS=windows go build -o d3-domain-tool.exe

# Format code
go fmt ./...

# Static analysis
go vet ./...
```

## Future Enhancements

- FastAPI backend integration
- Additional blockchain domain support
- Historical domain data
- Bulk domain analysis
- Domain monitoring
- Integration with domain marketplaces

## License

MIT License
