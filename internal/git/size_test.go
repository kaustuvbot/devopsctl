package git

import (
	"testing"

	appconfig "github.com/kaustuvbot/devopsctl/internal/config"
)

// TestConvertToMB tests all unit conversions
func TestConvertToMB(t *testing.T) {
	tests := []struct {
		value    float64
		unit     string
		expected float64
	}{
		{100, "B", 100.0 / (1024 * 1024)},
		{100, "KB", 100.0 / 1024},
		{100, "MB", 100.0},
		{1, "GB", 1024.0},
		{1, "TB", 1024.0 * 1024},
		{100, "MiB", 100.0},
		{1, "GiB", 1024.0},
	}

	for _, tt := range tests {
		result := convertToMB(tt.value, tt.unit)
		if result != tt.expected {
			t.Errorf("convertToMB(%v, %s) = %v; want %v", tt.value, tt.unit, result, tt.expected)
		}
	}
}

// TestRepoSizeThreshold tests threshold comparison logic
func TestRepoSizeThreshold(t *testing.T) {
	cfg := appconfig.GitConfig{
		RepoSizeMB:    500,
		BranchAgeDays: 90,
		LargeFileMB:   50,
	}

	tests := []struct {
		size     float64
		unit     string
		above    bool
	}{
		{10, "MiB", false},
		{100, "MiB", false},
		{500, "MiB", false},
		{501, "MiB", true},
		{1, "GiB", true},
	}

	for _, tt := range tests {
		sizeMB := convertToMB(tt.size, tt.unit)
		isAbove := sizeMB > float64(cfg.RepoSizeMB)
		if isAbove != tt.above {
			t.Errorf("convertToMB(%v, %s) = %v MB, expected above=%v", tt.size, tt.unit, sizeMB, tt.above)
		}
	}
}
