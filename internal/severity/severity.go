package severity

// Level represents the severity of a finding.
type Level string

const (
	Low      Level = "LOW"
	Medium   Level = "MEDIUM"
	High     Level = "HIGH"
	Critical Level = "CRITICAL"
)

// ExitCode returns the exit code corresponding to a severity level.
func (l Level) ExitCode() int {
	switch l {
	case Low:
		return 1
	case Medium:
		return 2
	case High:
		return 3
	case Critical:
		return 4
	default:
		return 0
	}
}

// Weight returns a numeric weight for severity comparison.
func (l Level) Weight() int {
	switch l {
	case Low:
		return 1
	case Medium:
		return 2
	case High:
		return 3
	case Critical:
		return 4
	default:
		return 0
	}
}

// AllLevels returns all severity levels in ascending order.
func AllLevels() []Level {
	return []Level{Low, Medium, High, Critical}
}

// IsValid checks if a severity string is a recognized level.
func IsValid(s string) bool {
	switch Level(s) {
	case Low, Medium, High, Critical:
		return true
	default:
		return false
	}
}

// Highest returns the highest severity from a list of levels.
func Highest(levels []Level) Level {
	highest := Level("")
	for _, l := range levels {
		if l.Weight() > highest.Weight() {
			highest = l
		}
	}
	return highest
}
