package valuation

import (
	"math"
	"strings"
	"unicode"
)

type Engine struct {
	premiumWords []string
	commonTLDs   map[string]float64
}

type Result struct {
	EstimatedValue   int     `json:"estimated_value"`
	Currency         string  `json:"currency"`
	Confidence       string  `json:"confidence"`
	Factors          Factors `json:"factors"`
	Reasoning        string  `json:"reasoning"`
}

type Factors struct {
	Length           int     `json:"length"`
	LengthScore      float64 `json:"length_score"`
	CharacterScore   float64 `json:"character_score"`
	WordScore        float64 `json:"word_score"`
	TLDScore         float64 `json:"tld_score"`
	Pronounceable    bool    `json:"pronounceable"`
	Brandable        bool    `json:"brandable"`
	HasNumbers       bool    `json:"has_numbers"`
	HasHyphens       bool    `json:"has_hyphens"`
}

func NewEngine() *Engine {
	return &Engine{
		premiumWords: []string{
			"app", "web", "tech", "crypto", "blockchain", "ai", "ml", "data",
			"cloud", "api", "dev", "code", "digital", "online", "smart",
			"auto", "health", "finance", "bank", "pay", "shop", "store",
			"game", "play", "social", "network", "security", "privacy",
		},
		commonTLDs: map[string]float64{
			".com":  1.0,
			".net":  0.7,
			".org":  0.6,
			".io":   0.8,
			".co":   0.6,
			".app":  0.7,
			".dev":  0.6,
			".tech": 0.5,
			".eth":  0.9, // High value for blockchain domains
			".crypto": 0.8,
			".nft":  0.7,
		},
	}
}

func (e *Engine) Evaluate(domain string) *Result {
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return &Result{
			EstimatedValue: 0,
			Currency:       "USD",
			Confidence:     "low",
			Reasoning:      "Invalid domain format",
		}
	}

	name := parts[0]
	tld := "." + parts[len(parts)-1]

	factors := e.analyzeDomain(name, tld)
	value := e.calculateValue(factors)
	confidence := e.determineConfidence(factors)
	reasoning := e.generateReasoning(factors)

	return &Result{
		EstimatedValue: int(value),
		Currency:       "USD",
		Confidence:     confidence,
		Factors:        factors,
		Reasoning:      reasoning,
	}
}

func (e *Engine) analyzeDomain(name, tld string) Factors {
	factors := Factors{
		Length:     len(name),
		HasNumbers: containsNumbers(name),
		HasHyphens: strings.Contains(name, "-"),
	}

	// Length scoring (shorter is generally better)
	factors.LengthScore = e.calculateLengthScore(len(name))

	// Character composition scoring
	factors.CharacterScore = e.calculateCharacterScore(name)

	// Word/brandability scoring
	factors.WordScore = e.calculateWordScore(name)

	// TLD scoring
	factors.TLDScore = e.calculateTLDScore(tld)

	// Pronounceable check
	factors.Pronounceable = e.isPronounceableWord(name)

	// Brandable check
	factors.Brandable = e.isBrandable(name)

	return factors
}

func (e *Engine) calculateLengthScore(length int) float64 {
	switch {
	case length <= 3:
		return 10.0 // Premium short domains
	case length <= 5:
		return 8.0
	case length <= 7:
		return 6.0
	case length <= 10:
		return 4.0
	case length <= 15:
		return 2.0
	default:
		return 1.0
	}
}

func (e *Engine) calculateCharacterScore(name string) float64 {
	score := 5.0

	// Penalize numbers and hyphens
	if containsNumbers(name) {
		score -= 2.0
	}
	if strings.Contains(name, "-") {
		score -= 1.5
	}

	// Bonus for all letters
	if isAllLetters(name) {
		score += 1.0
	}

	// Penalize mixed case inconsistency
	if hasMixedCase(name) {
		score -= 0.5
	}

	return math.Max(0, score)
}

func (e *Engine) calculateWordScore(name string) float64 {
	score := 0.0
	nameLower := strings.ToLower(name)

	// Check for premium words
	for _, word := range e.premiumWords {
		if strings.Contains(nameLower, word) {
			score += 3.0
		}
	}

	// Check if it's a dictionary word (simplified check)
	if e.isLikelyDictionaryWord(nameLower) {
		score += 2.0
	}

	// Bonus for compound words
	if e.isLikelyCompoundWord(nameLower) {
		score += 1.0
	}

	return score
}

func (e *Engine) calculateTLDScore(tld string) float64 {
	if score, exists := e.commonTLDs[tld]; exists {
		return score * 5.0 // Scale to match other scoring
	}
	return 1.0 // Default for unknown TLDs
}

func (e *Engine) calculateValue(factors Factors) float64 {
	baseValue := 100.0 // Minimum base value

	// Apply multiplicative factors
	multiplier := 1.0
	multiplier *= (factors.LengthScore / 10.0) * 2.0     // Length is very important
	multiplier *= (factors.CharacterScore / 5.0)         // Character quality
	multiplier *= (factors.TLDScore / 5.0)               // TLD premium
	multiplier += (factors.WordScore / 10.0)             // Word bonus (additive)

	// Brandability bonuses
	if factors.Brandable {
		multiplier *= 1.5
	}
	if factors.Pronounceable {
		multiplier *= 1.2
	}

	// Penalties
	if factors.HasNumbers {
		multiplier *= 0.7
	}
	if factors.HasHyphens {
		multiplier *= 0.6
	}

	value := baseValue * multiplier

	// Apply some realistic bounds
	if value < 10 {
		value = 10
	}
	if value > 1000000 {
		value = 1000000
	}

	return value
}

func (e *Engine) determineConfidence(factors Factors) string {
	score := 0

	if factors.Length <= 5 {
		score += 2
	}
	if factors.Brandable {
		score += 2
	}
	if factors.Pronounceable {
		score += 1
	}
	if factors.TLDScore >= 4.0 {
		score += 2
	}
	if !factors.HasNumbers && !factors.HasHyphens {
		score += 1
	}

	switch {
	case score >= 6:
		return "high"
	case score >= 3:
		return "medium"
	default:
		return "low"
	}
}

func (e *Engine) generateReasoning(factors Factors) string {
	var reasons []string

	if factors.Length <= 3 {
		reasons = append(reasons, "Very short domain (premium)")
	} else if factors.Length <= 5 {
		reasons = append(reasons, "Short and memorable")
	} else if factors.Length > 15 {
		reasons = append(reasons, "Long domain name")
	}

	if factors.Brandable {
		reasons = append(reasons, "Brandable name")
	}

	if factors.Pronounceable {
		reasons = append(reasons, "Easy to pronounce")
	}

	if factors.WordScore > 2 {
		reasons = append(reasons, "Contains valuable keywords")
	}

	if factors.HasNumbers {
		reasons = append(reasons, "Contains numbers (reduces value)")
	}

	if factors.HasHyphens {
		reasons = append(reasons, "Contains hyphens (reduces value)")
	}

	if len(reasons) == 0 {
		return "Standard domain name"
	}

	return strings.Join(reasons, "; ")
}

// Helper functions
func containsNumbers(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func isAllLetters(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func hasMixedCase(s string) bool {
	hasUpper := false
	hasLower := false
	for _, r := range s {
		if unicode.IsUpper(r) {
			hasUpper = true
		}
		if unicode.IsLower(r) {
			hasLower = true
		}
	}
	return hasUpper && hasLower
}

func (e *Engine) isPronounceableWord(name string) bool {
	// Simple heuristic: check vowel distribution
	vowels := "aeiou"
	vowelCount := 0
	consonantCount := 0

	for _, r := range strings.ToLower(name) {
		if strings.ContainsRune(vowels, r) {
			vowelCount++
		} else if unicode.IsLetter(r) {
			consonantCount++
		}
	}

	if vowelCount == 0 {
		return false
	}

	ratio := float64(consonantCount) / float64(vowelCount)
	return ratio >= 0.5 && ratio <= 4.0
}

func (e *Engine) isBrandable(name string) bool {
	// Simple brandability heuristics
	if len(name) < 3 || len(name) > 12 {
		return false
	}

	if containsNumbers(name) || strings.Contains(name, "-") {
		return false
	}

	if !e.isPronounceableWord(name) {
		return false
	}

	return true
}

func (e *Engine) isLikelyDictionaryWord(name string) bool {
	// Very simplified dictionary word detection
	commonWords := []string{
		"app", "web", "net", "tech", "data", "info", "news", "shop", "store",
		"game", "play", "work", "home", "life", "love", "time", "world",
		"best", "new", "top", "first", "last", "good", "great", "super",
	}

	for _, word := range commonWords {
		if name == word {
			return true
		}
	}

	return false
}

func (e *Engine) isLikelyCompoundWord(name string) bool {
	// Simple compound word detection
	if len(name) < 6 {
		return false
	}

	// Look for common prefixes/suffixes
	prefixes := []string{"web", "app", "my", "get", "the", "new", "top", "best"}
	suffixes := []string{"app", "web", "net", "tech", "hub", "lab", "pro", "max"}

	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) && len(name) > len(prefix)+2 {
			return true
		}
	}

	for _, suffix := range suffixes {
		if strings.HasSuffix(name, suffix) && len(name) > len(suffix)+2 {
			return true
		}
	}

	return false
}