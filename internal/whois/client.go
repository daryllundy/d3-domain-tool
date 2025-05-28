package whois

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type Client struct {
	timeout time.Duration
}

type Result struct {
	Available       bool       `json:"available"`
	Registrar       string     `json:"registrar,omitempty"`
	RegistrationDate *time.Time `json:"registration_date,omitempty"`
	ExpiryDate      *time.Time `json:"expiry_date,omitempty"`
	NameServers     []string   `json:"name_servers,omitempty"`
	Status          []string   `json:"status,omitempty"`
	UpdatedDate     *time.Time `json:"updated_date,omitempty"`
	CheckedAt       time.Time  `json:"checked_at"`
	RawData         string     `json:"raw_data,omitempty"`
	Error           string     `json:"error,omitempty"`
}

func NewClient() *Client {
	return &Client{
		timeout: 10 * time.Second,
	}
}

func (c *Client) Lookup(domain string) (*Result, error) {
	result := &Result{
		CheckedAt: time.Now(),
	}

	whoisServer := c.getWhoisServer(domain)
	if whoisServer == "" {
		result.Error = "No WHOIS server found for domain"
		return result, nil
	}

	rawData, err := c.queryWhoisServer(whoisServer, domain)
	if err != nil {
		result.Error = err.Error()
		return result, nil
	}

	result.RawData = rawData
	c.parseWhoisData(rawData, result)

	return result, nil
}

func (c *Client) getWhoisServer(domain string) string {
	tld := extractTLD(domain)
	
	whoisServers := map[string]string{
		".com":  "whois.verisign-grs.com",
		".net":  "whois.verisign-grs.com",
		".org":  "whois.pir.org",
		".info": "whois.afilias.net",
		".biz":  "whois.neulevel.biz",
		".name": "whois.nic.name",
		".io":   "whois.nic.io",
		".co":   "whois.nic.co",
		".me":   "whois.nic.me",
		".tv":   "whois.nic.tv",
		".cc":   "ccwhois.verisign-grs.com",
		".ws":   "whois.website.ws",
	}

	return whoisServers[tld]
}

func (c *Client) queryWhoisServer(server, domain string) (string, error) {
	conn, err := net.DialTimeout("tcp", server+":43", c.timeout)
	if err != nil {
		return "", fmt.Errorf("failed to connect to WHOIS server: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(domain + "\r\n"))
	if err != nil {
		return "", fmt.Errorf("failed to send query: %v", err)
	}

	var response strings.Builder
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		response.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	return response.String(), nil
}

func (c *Client) parseWhoisData(rawData string, result *Result) {
	lines := strings.Split(rawData, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check for "No match" or similar indicators of availability
		if strings.Contains(strings.ToLower(line), "no match") ||
		   strings.Contains(strings.ToLower(line), "not found") ||
		   strings.Contains(strings.ToLower(line), "no data found") {
			result.Available = true
			return
		}

		// Parse common WHOIS fields
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(strings.ToLower(parts[0]))
			value := strings.TrimSpace(parts[1])

			switch key {
			case "registrar":
				result.Registrar = value
			case "creation date", "created", "registration time":
				if date, err := parseDate(value); err == nil {
					result.RegistrationDate = &date
				}
			case "expiry date", "expires", "expiration time":
				if date, err := parseDate(value); err == nil {
					result.ExpiryDate = &date
				}
			case "updated date", "last modified", "last updated":
				if date, err := parseDate(value); err == nil {
					result.UpdatedDate = &date
				}
			case "name server":
				result.NameServers = append(result.NameServers, value)
			case "status", "domain status":
				result.Status = append(result.Status, value)
			}
		}
	}

	// If we parsed data, domain is not available
	if result.Registrar != "" || result.RegistrationDate != nil {
		result.Available = false
	}
}

func parseDate(dateStr string) (time.Time, error) {
	dateFormats := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"02-Jan-2006",
		"2006/01/02",
	}

	for _, format := range dateFormats {
		if date, err := time.Parse(format, dateStr); err == nil {
			return date, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func extractTLD(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return ""
	}
	return "." + parts[len(parts)-1]
}