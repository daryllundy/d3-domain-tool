package checker

import (
	"net"
	"strings"
	"time"
)

type DNSChecker struct {
	timeout time.Duration
}

type DNSResult struct {
	Available    bool              `json:"available"`
	TLD          string            `json:"tld"`
	HasRecords   bool              `json:"has_records"`
	RecordTypes  []string          `json:"record_types"`
	CheckedAt    time.Time         `json:"checked_at"`
	Error        string            `json:"error,omitempty"`
}

func NewDNSChecker() *DNSChecker {
	return &DNSChecker{
		timeout: 5 * time.Second,
	}
}

func (c *DNSChecker) Check(domain string) (*DNSResult, error) {
	result := &DNSResult{
		TLD:       extractTLD(domain),
		CheckedAt: time.Now(),
	}

	// Check for A records
	aRecords, err := net.LookupHost(domain)
	if err == nil && len(aRecords) > 0 {
		result.HasRecords = true
		result.RecordTypes = append(result.RecordTypes, "A")
		result.Available = false
	}

	// Check for MX records
	mxRecords, err := net.LookupMX(domain)
	if err == nil && len(mxRecords) > 0 {
		result.HasRecords = true
		result.RecordTypes = append(result.RecordTypes, "MX")
		result.Available = false
	}

	// Check for NS records
	nsRecords, err := net.LookupNS(domain)
	if err == nil && len(nsRecords) > 0 {
		result.HasRecords = true
		result.RecordTypes = append(result.RecordTypes, "NS")
		result.Available = false
	}

	// Check for TXT records
	txtRecords, err := net.LookupTXT(domain)
	if err == nil && len(txtRecords) > 0 {
		result.HasRecords = true
		result.RecordTypes = append(result.RecordTypes, "TXT")
		result.Available = false
	}

	// If no records found, likely available
	if !result.HasRecords {
		result.Available = true
	}

	return result, nil
}

func extractTLD(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return ""
	}
	return "." + parts[len(parts)-1]
}