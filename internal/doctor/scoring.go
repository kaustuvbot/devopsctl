package doctor

import (
	"github.com/kaustuvprajapati/devopsctl/internal/severity"
)

// Summary holds aggregated scoring information.
type Summary struct {
	TotalFindings int            `json:"total_findings"`
	Critical      int            `json:"critical"`
	High          int            `json:"high"`
	Medium        int            `json:"medium"`
	Low           int            `json:"low"`
	Score         int            `json:"score"`
	ModulesFailed int            `json:"modules_failed"`
	ModuleErrors  map[string]string `json:"module_errors,omitempty"`
}

// ComputeSummary calculates aggregate statistics from module reports.
func ComputeSummary(reports []ModuleReport) Summary {
	summary := Summary{
		ModuleErrors: make(map[string]string),
	}

	for _, report := range reports {
		if report.Error != "" {
			summary.ModulesFailed++
			summary.ModuleErrors[report.Module] = report.Error
			continue
		}

		for _, result := range report.Results {
			summary.TotalFindings++

			switch severity.Level(result.Severity) {
			case severity.Critical:
				summary.Critical++
				summary.Score += severity.Critical.Weight()
			case severity.High:
				summary.High++
				summary.Score += severity.High.Weight()
			case severity.Medium:
				summary.Medium++
				summary.Score += severity.Medium.Weight()
			case severity.Low:
				summary.Low++
				summary.Score += severity.Low.Weight()
			}
		}
	}

	return summary
}

// HighestSeverity returns the highest severity level from all results.
func HighestSeverity(reports []ModuleReport) severity.Level {
	levels := make([]severity.Level, 0)

	for _, report := range reports {
		for _, result := range report.Results {
			levels = append(levels, severity.Level(result.Severity))
		}
	}

	if len(levels) == 0 {
		return ""
	}

	return severity.Highest(levels)
}

// ExitCode returns the exit code based on highest severity.
func ExitCode(reports []ModuleReport) int {
	highest := HighestSeverity(reports)
	if highest == "" {
		return 0
	}
	return highest.ExitCode()
}
