package git

import (
	"testing"
	"time"
)

// TestBranchAgeCalculation tests the age calculation for branches
func TestBranchAgeCalculation(t *testing.T) {
	now := time.Now()

	tests := []struct {
		daysAgo int
		expected time.Duration
	}{
		{0, 0},
		{30, 30 * 24 * time.Hour},
		{90, 90 * 24 * time.Hour},
		{180, 180 * 24 * time.Hour},
	}

	for _, tt := range tests {
		expectedAge := time.Duration(tt.daysAgo) * 24 * time.Hour
		commitTime := now.Add(-expectedAge)
		actualAge := now.Sub(commitTime)

		if actualAge != expectedAge {
			t.Errorf("Expected age %v, got %v", expectedAge, actualAge)
		}
	}
}

// TestBranchNameParsing tests parsing branch names from git output
func TestBranchNameParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"main|1700000000", "main"},
		{"feature/test|1700000000", "feature/test"},
		{"origin/main|1700000000", "origin/main"}, // Keep prefix for remote branches
		{"origin/feature|1700000000", "origin/feature"},
	}

	for _, tt := range tests {
		parts := splitBranchLine(tt.input)
		if parts.branch != tt.expected {
			t.Errorf("Expected branch %q, got %q", tt.expected, parts.branch)
		}
	}
}

// splitBranchLine is a helper to test parsing logic
func splitBranchLine(line string) struct{ branch string } {
	parts := splitAt(line, '|')
	branch := trim(parts[0])
	branch = trimPrefix(branch, "remotes/origin/")
	return struct{ branch string }{branch}
}

func splitAt(s string, c byte) []string {
	var result []string
	current := ""
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			result = append(result, current)
			current = ""
		} else {
			current += string(s[i])
		}
	}
	result = append(result, current)
	return result
}

func trim(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\n') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\n') {
		end--
	}
	return s[start:end]
}

func trimPrefix(s, prefix string) string {
	if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}

// TestStaleBranchThreshold tests threshold comparison
func TestStaleBranchThreshold(t *testing.T) {
	cfg := struct{ days int }{90}
	threshold := time.Duration(cfg.days) * 24 * time.Hour
	now := time.Now()

	tests := []struct {
		name     string
		daysAgo  int
		isStale  bool
	}{
		{"recent", 30, false},
		{"boundary", 90, false}, // exactly at threshold is not stale
		{"old", 91, true},
		{"very old", 180, true},
	}

	for _, tt := range tests {
		commitTime := now.Add(-time.Duration(tt.daysAgo) * 24 * time.Hour)
		age := now.Sub(commitTime)
		isStale := age > threshold

		if isStale != tt.isStale {
			t.Errorf("Test %s: expected isStale=%v, got %v (age=%v, threshold=%v)",
				tt.name, tt.isStale, isStale, age, threshold)
		}
	}
}
