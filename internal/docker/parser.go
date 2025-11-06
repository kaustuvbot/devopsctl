package docker

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Instruction represents a single parsed Dockerfile directive.
type Instruction struct {
	Command string // Uppercased: FROM, RUN, USER, EXPOSE, HEALTHCHECK, etc.
	Args    string // Everything after the command on the logical line
	LineNum int    // Line number of the instruction's first line in the file
}

// ParsedDockerfile holds all instructions parsed from a Dockerfile.
type ParsedDockerfile struct {
	Instructions []Instruction
	Path         string
}

// ParseDockerfile reads and parses a Dockerfile into structured instructions.
// It handles comment lines, blank lines, and backslash line continuations.
// Parser directives (# syntax=...) at the top are treated as comments.
func ParseDockerfile(path string) (*ParsedDockerfile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open Dockerfile %q: %w", path, err)
	}
	defer f.Close()

	var instructions []Instruction
	scanner := bufio.NewScanner(f)
	lineNum := 0
	var (
		accumulated string
		startLine   int
	)

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Skip blank lines and comments when not in a continuation
		if accumulated == "" && (trimmed == "" || strings.HasPrefix(trimmed, "#")) {
			continue
		}

		// Line continuation: join next line
		if strings.HasSuffix(trimmed, "\\") {
			if accumulated == "" {
				startLine = lineNum
			}
			accumulated += strings.TrimSuffix(trimmed, "\\") + " "
			continue
		}

		// End of a logical line
		if accumulated != "" {
			accumulated += trimmed
			trimmed = strings.TrimSpace(accumulated)
			accumulated = ""
		} else {
			startLine = lineNum
		}

		if trimmed == "" {
			continue
		}

		parts := strings.SplitN(trimmed, " ", 2)
		cmd := strings.ToUpper(parts[0])
		args := ""
		if len(parts) == 2 {
			args = strings.TrimSpace(parts[1])
		}

		instructions = append(instructions, Instruction{
			Command: cmd,
			Args:    args,
			LineNum: startLine,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading %q: %w", path, err)
	}

	return &ParsedDockerfile{Instructions: instructions, Path: path}, nil
}
