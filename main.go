package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"d3-domain-tool/internal/analyzer"
	"d3-domain-tool/internal/output"
)

func main() {
	var (
		domain = flag.String("domain", "", "Domain to analyze (required)")
		format = flag.String("format", "table", "Output format: table, json")
		help   = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help || *domain == "" {
		showUsage()
		return
	}

	cleanDomain := strings.TrimSpace(strings.ToLower(*domain))
	if cleanDomain == "" {
		fmt.Fprintf(os.Stderr, "Error: Domain cannot be empty\n")
		os.Exit(1)
	}

	analyzer := analyzer.New()
	result, err := analyzer.AnalyzeDomain(cleanDomain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error analyzing domain: %v\n", err)
		os.Exit(1)
	}

	formatter := output.NewFormatter(*format)
	if err := formatter.Display(result); err != nil {
		fmt.Fprintf(os.Stderr, "Error displaying results: %v\n", err)
		os.Exit(1)
	}
}

func showUsage() {
	fmt.Println("D3 Domain Analysis Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  d3-domain-tool -domain=<domain> [-format=table|json]")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  d3-domain-tool -domain=example.com")
	fmt.Println("  d3-domain-tool -domain=mydomain.eth -format=json")
	fmt.Println()
	fmt.Println("Features:")
	fmt.Println("  ✅ Check domain availability (DNS + blockchain)")
	fmt.Println("  🔍 WHOIS data and blockchain metadata")
	fmt.Println("  💰 Domain value estimation")
	fmt.Println("  📦 Clean CLI output")
}