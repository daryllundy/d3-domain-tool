package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"d3-domain-tool/internal/analyzer"
)

type Formatter struct {
	format string
}

func NewFormatter(format string) *Formatter {
	return &Formatter{
		format: format,
	}
}

func (f *Formatter) Display(result *analyzer.Result) error {
	switch f.format {
	case "json":
		return f.displayJSON(result)
	case "table":
		return f.displayTable(result)
	default:
		return fmt.Errorf("unsupported format: %s", f.format)
	}
}

func (f *Formatter) displayJSON(result *analyzer.Result) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func (f *Formatter) displayTable(result *analyzer.Result) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Header
	fmt.Fprintf(w, "\n🔍 D3 DOMAIN ANALYSIS REPORT\n")
	fmt.Fprintf(w, "═══════════════════════════════════════════════════════════════\n\n")

	// Basic Info
	fmt.Fprintf(w, "Domain:\t%s\n", result.Domain)
	fmt.Fprintf(w, "Analyzed:\t%s\n\n", result.Timestamp.Format("2006-01-02 15:04:05 MST"))

	// DNS Availability Section
	if result.DNSAvailability != nil {
		fmt.Fprintf(w, "📡 DNS AVAILABILITY\n")
		fmt.Fprintf(w, "───────────────────\n")

		status := "❌ Taken"
		if result.DNSAvailability.Available {
			status = "✅ Available"
		}
		fmt.Fprintf(w, "Status:\t%s\n", status)
		fmt.Fprintf(w, "TLD:\t%s\n", result.DNSAvailability.TLD)

		if result.DNSAvailability.HasRecords {
			fmt.Fprintf(w, "Records:\t%s\n", strings.Join(result.DNSAvailability.RecordTypes, ", "))
		}

		if result.DNSAvailability.Error != "" {
			fmt.Fprintf(w, "Error:\t%s\n", result.DNSAvailability.Error)
		}
		fmt.Fprintf(w, "\n")
	}

	// DOMA Protocol Section
	if result.DomaData != nil {
		fmt.Fprintf(w, "🔶 DOMA PROTOCOL INTEGRATION\n")
		fmt.Fprintf(w, "───────────────────────────\n")

		tokenizedIcon := "❌"
		if result.DomaData.IsTokenized {
			tokenizedIcon = "✅"
		}
		fmt.Fprintf(w, "Tokenized:\t%s\n", tokenizedIcon)

		if result.DomaData.IsTokenized {
			if result.DomaData.TokenizationChain != "" {
				fmt.Fprintf(w, "Chain:\t%s\n", result.DomaData.TokenizationChain)
			}

			// DOMA Record Information
			if result.DomaData.DomaRecord != nil {
				record := result.DomaData.DomaRecord
				fmt.Fprintf(w, "Token ID:\t%s\n", record.TokenId)
				fmt.Fprintf(w, "Owner:\t%s\n", record.Owner)

				if record.RegistrationDate != nil {
					fmt.Fprintf(w, "Registered:\t%s\n", record.RegistrationDate.Format("2006-01-02"))
				}

				if record.ExpirationDate != nil {
					fmt.Fprintf(w, "Expires:\t%s\n", record.ExpirationDate.Format("2006-01-02"))
				}

				fmt.Fprintf(w, "Sync Status:\t%s\n", record.SyncStatus)
			}

			// Token Rights Information
			if result.DomaData.TokenRights != nil {
				rights := result.DomaData.TokenRights
				fmt.Fprintf(w, "\n🪙 Token Rights:\n")
				fmt.Fprintf(w, "  Total Tokens:\t%d\n", rights.Total)
				fmt.Fprintf(w, "  Available:\t%d\n", rights.Available)
				fmt.Fprintf(w, "  Locked:\t%d\n", rights.Locked)

				if len(rights.FractionalOwners) > 0 {
					fmt.Fprintf(w, "  Owners:\t%d\n", len(rights.FractionalOwners))
				}
			}

			// DeFi Status
			if result.DomaData.DeFiStatus != nil {
				defi := result.DomaData.DeFiStatus
				fmt.Fprintf(w, "\n💎 DeFi Integration:\n")

				collateralIcon := "❌"
				if defi.IsCollateral {
					collateralIcon = "✅"
				}
				fmt.Fprintf(w, "  Used as Collateral:\t%s\n", collateralIcon)

				if defi.LendingPlatform != "" {
					fmt.Fprintf(w, "  Platform:\t%s\n", defi.LendingPlatform)
					fmt.Fprintf(w, "  Collateral Value:\t$%.2f\n", defi.CollateralValue)
					fmt.Fprintf(w, "  Borrowed:\t$%.2f\n", defi.BorrowedAmount)
				}

				yieldIcon := "❌"
				if defi.YieldGeneration {
					yieldIcon = "✅"
				}
				fmt.Fprintf(w, "  Yield Generation:\t%s\n", yieldIcon)

				if defi.StakingRewards > 0 {
					fmt.Fprintf(w, "  Staking Rewards:\t$%.2f\n", defi.StakingRewards)
				}
			}

			// Cross-Chain Data
			if len(result.DomaData.CrossChainData) > 0 {
				fmt.Fprintf(w, "\n🌐 Cross-Chain Presence:\n")
				for chain := range result.DomaData.CrossChainData {
					fmt.Fprintf(w, "  %s:\t✅ Deployed\n", strings.Title(chain))
				}
			}
		} else {
			// Check eligibility for non-tokenized domains
			fmt.Fprintf(w, "Eligibility:\t⚠️ Not currently tokenized\n")
			fmt.Fprintf(w, "Note:\tTraditional domains can be tokenized via DOMA Protocol\n")
		}

		if result.DomaData.Error != "" {
			fmt.Fprintf(w, "Error:\t%s\n", result.DomaData.Error)
		}
		fmt.Fprintf(w, "\n")
	}

	// Blockchain Section
	if result.BlockchainData != nil {
		fmt.Fprintf(w, "⛓️ BLOCKCHAIN DATA\n")
		fmt.Fprintf(w, "──────────────────\n")

		status := "❌ Taken"
		if result.BlockchainData.Available {
			status = "✅ Available"
		}
		fmt.Fprintf(w, "Status:\t%s\n", status)
		fmt.Fprintf(w, "Type:\t%s\n", result.BlockchainData.Type)

		if result.BlockchainData.Owner != "" {
			fmt.Fprintf(w, "Owner:\t%s\n", result.BlockchainData.Owner)
		}

		if result.BlockchainData.Resolver != "" {
			fmt.Fprintf(w, "Resolver:\t%s\n", result.BlockchainData.Resolver)
		}

		if len(result.BlockchainData.Records) > 0 {
			fmt.Fprintf(w, "Records:\n")
			for key, value := range result.BlockchainData.Records {
				fmt.Fprintf(w, "  %s:\t%s\n", key, value)
			}
		}

		if result.BlockchainData.ExpiryDate != nil {
			fmt.Fprintf(w, "Expires:\t%s\n", result.BlockchainData.ExpiryDate.Format("2006-01-02"))
		}
		fmt.Fprintf(w, "\n")
	}

	// WHOIS Section
	if result.WhoisData != nil {
		fmt.Fprintf(w, "📋 WHOIS DATA\n")
		fmt.Fprintf(w, "─────────────\n")

		status := "❌ Taken"
		if result.WhoisData.Available {
			status = "✅ Available"
		}
		fmt.Fprintf(w, "Status:\t%s\n", status)

		if result.WhoisData.Registrar != "" {
			fmt.Fprintf(w, "Registrar:\t%s\n", result.WhoisData.Registrar)
		}

		if result.WhoisData.RegistrationDate != nil {
			fmt.Fprintf(w, "Created:\t%s\n", result.WhoisData.RegistrationDate.Format("2006-01-02"))
		}

		if result.WhoisData.ExpiryDate != nil {
			fmt.Fprintf(w, "Expires:\t%s\n", result.WhoisData.ExpiryDate.Format("2006-01-02"))
		}

		if result.WhoisData.UpdatedDate != nil {
			fmt.Fprintf(w, "Updated:\t%s\n", result.WhoisData.UpdatedDate.Format("2006-01-02"))
		}

		if len(result.WhoisData.NameServers) > 0 {
			fmt.Fprintf(w, "Name Servers:\t%s\n", strings.Join(result.WhoisData.NameServers, ", "))
		}

		if len(result.WhoisData.Status) > 0 {
			fmt.Fprintf(w, "Status:\t%s\n", strings.Join(result.WhoisData.Status, ", "))
		}

		if result.WhoisData.Error != "" {
			fmt.Fprintf(w, "Error:\t%s\n", result.WhoisData.Error)
		}
		fmt.Fprintf(w, "\n")
	}

	// Valuation Section
	if result.ValuationData != nil {
		fmt.Fprintf(w, "💰 DOMAIN VALUATION\n")
		fmt.Fprintf(w, "───────────────────\n")

		fmt.Fprintf(w, "Estimated Value:\t$%d %s\n",
			result.ValuationData.EstimatedValue,
			result.ValuationData.Currency)

		confidence := result.ValuationData.Confidence
		confidenceIcon := "🟡"
		switch confidence {
		case "high":
			confidenceIcon = "🟢"
		case "low":
			confidenceIcon = "🔴"
		}
		fmt.Fprintf(w, "Confidence:\t%s %s\n", confidenceIcon, strings.Title(confidence))

		fmt.Fprintf(w, "Reasoning:\t%s\n", result.ValuationData.Reasoning)

		fmt.Fprintf(w, "\nValuation Factors:\n")
		factors := result.ValuationData.Factors
		fmt.Fprintf(w, "  Length:\t%d chars (Score: %.1f/10)\n", factors.Length, factors.LengthScore)
		fmt.Fprintf(w, "  Character Quality:\t%.1f/5\n", factors.CharacterScore)
		fmt.Fprintf(w, "  Word Value:\t%.1f/10\n", factors.WordScore)
		fmt.Fprintf(w, "  TLD Value:\t%.1f/5\n", factors.TLDScore)

		brandableIcon := "❌"
		if factors.Brandable {
			brandableIcon = "✅"
		}
		fmt.Fprintf(w, "  Brandable:\t%s\n", brandableIcon)

		pronounceableIcon := "❌"
		if factors.Pronounceable {
			pronounceableIcon = "✅"
		}
		fmt.Fprintf(w, "  Pronounceable:\t%s\n", pronounceableIcon)

		if factors.HasNumbers {
			fmt.Fprintf(w, "  Contains Numbers:\t❌ (reduces value)\n")
		}

		if factors.HasHyphens {
			fmt.Fprintf(w, "  Contains Hyphens:\t❌ (reduces value)\n")
		}
	}

	fmt.Fprintf(w, "\n")
	return w.Flush()
}
