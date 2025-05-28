package doma

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	timeout    time.Duration
}

type Result struct {
	Domain            string                 `json:"domain"`
	IsTokenized       bool                   `json:"is_tokenized"`
	TokenizationChain string                 `json:"tokenization_chain,omitempty"`
	DomaRecord        *DomaRecord            `json:"doma_record,omitempty"`
	TokenRights       *TokenRights           `json:"token_rights,omitempty"`
	DeFiStatus        *DeFiStatus            `json:"defi_status,omitempty"`
	CrossChainData    map[string]interface{} `json:"cross_chain_data,omitempty"`
	CheckedAt         time.Time              `json:"checked_at"`
	Error             string                 `json:"error,omitempty"`
}

type DomaRecord struct {
	TokenId          string            `json:"token_id"`
	Owner            string            `json:"owner"`
	Resolver         string            `json:"resolver"`
	Records          map[string]string `json:"records"`
	RegistrationDate *time.Time        `json:"registration_date"`
	ExpirationDate   *time.Time        `json:"expiration_date"`
	LastUpdated      *time.Time        `json:"last_updated"`
	SyncStatus       string            `json:"sync_status"`
}

type TokenRights struct {
	Total            int                    `json:"total_tokens"`
	Available        int                    `json:"available_tokens"`
	Locked           int                    `json:"locked_tokens"`
	RightsBreakdown  map[string]interface{} `json:"rights_breakdown"`
	FractionalOwners []string               `json:"fractional_owners"`
}

type DeFiStatus struct {
	IsCollateral    bool    `json:"is_collateral"`
	LendingPlatform string  `json:"lending_platform,omitempty"`
	CollateralValue float64 `json:"collateral_value,omitempty"`
	BorrowedAmount  float64 `json:"borrowed_amount,omitempty"`
	YieldGeneration bool    `json:"yield_generation"`
	StakingRewards  float64 `json:"staking_rewards,omitempty"`
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		baseURL: "https://api.doma.xyz",
		timeout: 15 * time.Second,
	}
}

func (c *Client) CheckDomain(domain string) (*Result, error) {
	result := &Result{
		Domain:         domain,
		CheckedAt:      time.Now(),
		CrossChainData: make(map[string]interface{}),
	}

	// Check if domain is tokenized on DOMA Protocol
	tokenized, err := c.isTokenized(domain)
	if err != nil {
		result.Error = err.Error()
		return result, nil
	}

	result.IsTokenized = tokenized

	if tokenized {
		// Get detailed DOMA record data
		record, err := c.getDomaRecord(domain)
		if err == nil {
			result.DomaRecord = record
		}

		// Get token rights information
		rights, err := c.getTokenRights(domain)
		if err == nil {
			result.TokenRights = rights
		}

		// Get DeFi status
		defiStatus, err := c.getDeFiStatus(domain)
		if err == nil {
			result.DeFiStatus = defiStatus
		}

		// Get cross-chain data
		crossChain, err := c.getCrossChainData(domain)
		if err == nil {
			result.CrossChainData = crossChain
		}

		// Determine tokenization chain
		result.TokenizationChain = c.getTokenizationChain(domain)
	}

	return result, nil
}

func (c *Client) isTokenized(domain string) (bool, error) {
	// In a real implementation, this would call the DOMA API
	// For now, simulate based on domain characteristics

	// Check if it's a traditional domain that could be tokenized
	if strings.Contains(domain, ".com") || strings.Contains(domain, ".net") ||
		strings.Contains(domain, ".org") || strings.Contains(domain, ".io") {
		// Simulate some domains being tokenized
		return c.simulateTokenizationStatus(domain), nil
	}

	// Blockchain domains might also be tokenized through DOMA
	if strings.Contains(domain, ".eth") || strings.Contains(domain, ".crypto") {
		return c.simulateTokenizationStatus(domain), nil
	}

	return false, nil
}

func (c *Client) getDomaRecord(domain string) (*DomaRecord, error) {
	// Simulate DOMA Record data
	now := time.Now()
	expiry := now.AddDate(1, 0, 0)        // 1 year from now
	registration := now.AddDate(-1, 0, 0) // 1 year ago

	return &DomaRecord{
		TokenId:  c.generateTokenId(domain),
		Owner:    "0x" + strings.Repeat("1", 40),
		Resolver: "0x" + strings.Repeat("2", 40),
		Records: map[string]string{
			"A":    "192.168.1.1",
			"AAAA": "2001:db8::1",
			"TXT":  "v=spf1 include:_spf.google.com ~all",
			"ETH":  "0x" + strings.Repeat("3", 40),
			"BTC":  "bc1" + strings.Repeat("4", 39),
		},
		RegistrationDate: &registration,
		ExpirationDate:   &expiry,
		LastUpdated:      &now,
		SyncStatus:       "synced",
	}, nil
}

func (c *Client) getTokenRights(domain string) (*TokenRights, error) {
	// Simulate token rights data
	return &TokenRights{
		Total:     1000,
		Available: 750,
		Locked:    250,
		RightsBreakdown: map[string]interface{}{
			"ownership":  500,
			"revenue":    300,
			"governance": 150,
			"utility":    50,
		},
		FractionalOwners: []string{
			"0x" + strings.Repeat("a", 40),
			"0x" + strings.Repeat("b", 40),
			"0x" + strings.Repeat("c", 40),
		},
	}, nil
}

func (c *Client) getDeFiStatus(domain string) (*DeFiStatus, error) {
	// Simulate DeFi integration status
	return &DeFiStatus{
		IsCollateral:    true,
		LendingPlatform: "DOMA Lending",
		CollateralValue: 50000.0,
		BorrowedAmount:  30000.0,
		YieldGeneration: true,
		StakingRewards:  125.50,
	}, nil
}

func (c *Client) getCrossChainData(domain string) (map[string]interface{}, error) {
	// Simulate cross-chain presence data
	return map[string]interface{}{
		"ethereum": map[string]interface{}{
			"contract_address": "0x" + strings.Repeat("e", 40),
			"token_id":         c.generateTokenId(domain),
			"last_update":      time.Now().Unix(),
		},
		"polygon": map[string]interface{}{
			"contract_address": "0x" + strings.Repeat("f", 40),
			"bridged":          true,
			"bridge_fee":       0.01,
		},
		"arbitrum": map[string]interface{}{
			"contract_address": "0x" + strings.Repeat("d", 40),
			"layer2_benefits":  true,
			"gas_savings":      "95%",
		},
	}, nil
}

func (c *Client) getTokenizationChain(domain string) string {
	// Determine primary tokenization chain
	// In practice, this would come from the API response
	return "ethereum"
}

func (c *Client) simulateTokenizationStatus(domain string) bool {
	// Simulate tokenization - premium/short domains more likely to be tokenized
	domainPart := strings.Split(domain, ".")[0]

	// Short domains (3 chars or less) are likely tokenized
	if len(domainPart) <= 3 {
		return true
	}

	// Common/premium domains might be tokenized
	premiumDomains := []string{"crypto", "defi", "nft", "web3", "blockchain", "ethereum", "bitcoin"}
	for _, premium := range premiumDomains {
		if strings.Contains(strings.ToLower(domainPart), premium) {
			return true
		}
	}

	// Dictionary words might be tokenized
	if len(domainPart) >= 4 && len(domainPart) <= 8 {
		return len(domainPart)%2 == 0 // Simulate 50% chance for medium length domains
	}

	return false
}

func (c *Client) generateTokenId(domain string) string {
	// Generate a simulated token ID based on domain
	hash := fmt.Sprintf("%x", []byte(domain))
	if len(hash) > 20 {
		hash = hash[:20]
	}
	return hash
}

// Helper function to check if domain could be eligible for DOMA tokenization
func (c *Client) IsEligibleForTokenization(domain string) (bool, string) {
	// Traditional domains are eligible
	traditionalTLDs := []string{".com", ".net", ".org", ".io", ".co", ".me", ".tv", ".cc", ".ws"}
	for _, tld := range traditionalTLDs {
		if strings.HasSuffix(domain, tld) {
			return true, "Traditional domain eligible for DOMA tokenization"
		}
	}

	// Some blockchain domains can also be bridged
	blockchainTLDs := []string{".eth", ".crypto"}
	for _, tld := range blockchainTLDs {
		if strings.HasSuffix(domain, tld) {
			return true, "Blockchain domain eligible for DOMA bridge"
		}
	}

	return false, "Domain type not supported for DOMA tokenization"
}
