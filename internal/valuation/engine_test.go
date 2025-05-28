package valuation

import (
	"testing"
)

func TestEngine_Evaluate(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name           string
		domain         string
		expectPositive bool
		expectBrandable bool
	}{
		{
			name:           "short domain",
			domain:         "app.com",
			expectPositive: true,
			expectBrandable: true,
		},
		{
			name:           "long domain",
			domain:         "verylongdomainnamethatishard.com",
			expectPositive: true,
			expectBrandable: false,
		},
		{
			name:           "numbers in domain",
			domain:         "test123.com",
			expectPositive: true,
			expectBrandable: false,
		},
		{
			name:           "hyphens in domain",
			domain:         "test-domain.com",
			expectPositive: true,
			expectBrandable: false,
		},
		{
			name:           "blockchain domain",
			domain:         "myname.eth",
			expectPositive: true,
			expectBrandable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.Evaluate(tt.domain)
			
			if result == nil {
				t.Fatal("Expected result, got nil")
			}

			if tt.expectPositive && result.EstimatedValue <= 0 {
				t.Errorf("Expected positive value, got %d", result.EstimatedValue)
			}

			if result.Factors.Brandable != tt.expectBrandable {
				t.Errorf("Expected brandable=%v, got %v", tt.expectBrandable, result.Factors.Brandable)
			}

			if result.Currency != "USD" {
				t.Errorf("Expected currency USD, got %s", result.Currency)
			}

			if result.Confidence == "" {
				t.Error("Expected confidence level, got empty string")
			}
		})
	}
}

func TestEngine_calculateLengthScore(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		length   int
		expected float64
	}{
		{1, 10.0},
		{3, 10.0},
		{5, 8.0},
		{7, 6.0},
		{10, 4.0},
		{15, 2.0},
		{20, 1.0},
	}

	for _, tt := range tests {
		score := engine.calculateLengthScore(tt.length)
		if score != tt.expected {
			t.Errorf("For length %d, expected score %f, got %f", tt.length, tt.expected, score)
		}
	}
}

func TestContainsNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"hello", false},
		{"hello123", true},
		{"123hello", true},
		{"hel3lo", true},
		{"", false},
	}

	for _, tt := range tests {
		result := containsNumbers(tt.input)
		if result != tt.expected {
			t.Errorf("For input %s, expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestIsAllLetters(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"hello", true},
		{"hello123", false},
		{"hello-world", false},
		{"", true},
		{"Hello", true},
	}

	for _, tt := range tests {
		result := isAllLetters(tt.input)
		if result != tt.expected {
			t.Errorf("For input %s, expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}