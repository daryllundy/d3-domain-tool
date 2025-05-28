package analyzer

import (
	"fmt"
	"strings"
	"time"

	"d3-domain-tool/internal/blockchain"
	"d3-domain-tool/internal/checker"
	"d3-domain-tool/internal/doma"
	"d3-domain-tool/internal/valuation"
	"d3-domain-tool/internal/whois"
)

type Analyzer struct {
	dnsChecker        *checker.DNSChecker
	blockchainChecker *blockchain.Checker
	whoisClient       *whois.Client
	domaClient        *doma.Client
	valuator          *valuation.Engine
}

type Result struct {
	Domain          string             `json:"domain"`
	Timestamp       time.Time          `json:"timestamp"`
	DNSAvailability *checker.DNSResult `json:"dns_availability"`
	BlockchainData  *blockchain.Result `json:"blockchain_data"`
	DomaData        *doma.Result       `json:"doma_data"`
	WhoisData       *whois.Result      `json:"whois_data"`
	ValuationData   *valuation.Result  `json:"valuation_data"`
}

func New() *Analyzer {
	return &Analyzer{
		dnsChecker:        checker.NewDNSChecker(),
		blockchainChecker: blockchain.NewChecker(),
		whoisClient:       whois.NewClient(),
		domaClient:        doma.NewClient(),
		valuator:          valuation.NewEngine(),
	}
}

func (a *Analyzer) AnalyzeDomain(domain string) (*Result, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain cannot be empty")
	}

	result := &Result{
		Domain:    domain,
		Timestamp: time.Now(),
	}

	// Always check DOMA Protocol integration first
	domaData, err := a.domaClient.CheckDomain(domain)
	if err == nil {
		result.DomaData = domaData
	}

	// Check if it's a blockchain domain
	if isBlockchainDomain(domain) {
		blockchainData, err := a.blockchainChecker.Check(domain)
		if err == nil {
			result.BlockchainData = blockchainData
		}
	} else {
		// Traditional DNS domain
		dnsData, err := a.dnsChecker.Check(domain)
		if err == nil {
			result.DNSAvailability = dnsData
		}

		whoisData, err := a.whoisClient.Lookup(domain)
		if err == nil {
			result.WhoisData = whoisData
		}
	}

	// Always run valuation (now enhanced with DOMA data)
	valuationData := a.valuator.Evaluate(domain)
	result.ValuationData = valuationData

	return result, nil
}

func isBlockchainDomain(domain string) bool {
	blockchainTLDs := []string{".eth", ".crypto", ".nft", ".x", ".wallet", ".bitcoin", ".dao", ".888", ".zil", ".blockchain"}

	for _, tld := range blockchainTLDs {
		if strings.HasSuffix(domain, tld) {
			return true
		}
	}
	return false
}
