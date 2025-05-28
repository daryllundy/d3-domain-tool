package blockchain

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Checker struct {
	client  *http.Client
	timeout time.Duration
}

type Result struct {
	Available     bool              `json:"available"`
	Type          string            `json:"type"`
	Owner         string            `json:"owner,omitempty"`
	Resolver      string            `json:"resolver,omitempty"`
	Records       map[string]string `json:"records,omitempty"`
	ExpiryDate    *time.Time        `json:"expiry_date,omitempty"`
	CheckedAt     time.Time         `json:"checked_at"`
	Error         string            `json:"error,omitempty"`
}

func NewChecker() *Checker {
	return &Checker{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		timeout: 10 * time.Second,
	}
}

func (c *Checker) Check(domain string) (*Result, error) {
	result := &Result{
		CheckedAt: time.Now(),
		Records:   make(map[string]string),
	}

	if strings.HasSuffix(domain, ".eth") {
		return c.checkENS(domain, result)
	} else if strings.HasSuffix(domain, ".crypto") || strings.HasSuffix(domain, ".nft") || 
		strings.HasSuffix(domain, ".x") || strings.HasSuffix(domain, ".wallet") ||
		strings.HasSuffix(domain, ".bitcoin") || strings.HasSuffix(domain, ".dao") ||
		strings.HasSuffix(domain, ".888") || strings.HasSuffix(domain, ".zil") {
		return c.checkUnstoppableDomains(domain, result)
	}

	return result, fmt.Errorf("unsupported blockchain domain type")
}

func (c *Checker) checkENS(domain string, result *Result) (*Result, error) {
	result.Type = "ENS"
	
	// Simulate ENS lookup - in a real implementation, you'd use web3 libraries
	// or call Ethereum nodes directly
	result.Available = c.simulateENSLookup(domain)
	
	if !result.Available {
		result.Owner = "0x" + strings.Repeat("a", 40) // Simulated address
		result.Resolver = "0x" + strings.Repeat("b", 40)
		result.Records["ETH"] = "0x" + strings.Repeat("c", 40)
		result.Records["BTC"] = "bc1" + strings.Repeat("d", 39)
	}

	return result, nil
}

func (c *Checker) checkUnstoppableDomains(domain string, result *Result) (*Result, error) {
	result.Type = "Unstoppable Domains"
	
	// Simulate Unstoppable Domains lookup
	result.Available = c.simulateUDLookup(domain)
	
	if !result.Available {
		result.Owner = "0x" + strings.Repeat("e", 40)
		result.Records["crypto.ETH.address"] = "0x" + strings.Repeat("f", 40)
		result.Records["crypto.BTC.address"] = "bc1" + strings.Repeat("g", 39)
	}

	return result, nil
}

// Simulate blockchain lookups - in production, these would make actual blockchain calls
func (c *Checker) simulateENSLookup(domain string) bool {
	// Simulate some domains being taken
	commonDomains := []string{"test.eth", "example.eth", "hello.eth", "world.eth"}
	for _, taken := range commonDomains {
		if domain == taken {
			return false
		}
	}
	// Most other domains are available in simulation
	return len(strings.Split(domain, ".")[0]) > 3
}

func (c *Checker) simulateUDLookup(domain string) bool {
	// Similar simulation for Unstoppable Domains
	commonDomains := []string{"test.crypto", "example.nft", "hello.x"}
	for _, taken := range commonDomains {
		if domain == taken {
			return false
		}
	}
	return len(strings.Split(domain, ".")[0]) > 3
}