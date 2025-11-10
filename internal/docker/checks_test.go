package docker

import (
	"testing"
)

func makeDF(t *testing.T, content string) *ParsedDockerfile {
	t.Helper()
	path := writeTemp(t, content)
	df, err := ParseDockerfile(path)
	if err != nil {
		t.Fatalf("makeDF: %v", err)
	}
	return df
}

func TestCheckLatestTag_Latest(t *testing.T) {
	df := makeDF(t, "FROM ubuntu:latest\nRUN echo hi\n")
	results := CheckLatestTag(df)
	if len(results) != 1 || results[0].Severity != "MEDIUM" {
		t.Errorf("expected 1 MEDIUM result for :latest, got %v", results)
	}
}

func TestCheckLatestTag_Untagged(t *testing.T) {
	df := makeDF(t, "FROM ubuntu\nRUN echo hi\n")
	results := CheckLatestTag(df)
	if len(results) != 1 {
		t.Errorf("expected 1 result for untagged image, got %d", len(results))
	}
}

func TestCheckLatestTag_Pinned(t *testing.T) {
	df := makeDF(t, "FROM ubuntu:22.04\nRUN echo hi\n")
	results := CheckLatestTag(df)
	if len(results) != 0 {
		t.Errorf("expected no results for pinned tag, got %v", results)
	}
}

func TestCheckLatestTag_Scratch(t *testing.T) {
	df := makeDF(t, "FROM scratch\nCOPY app /app\n")
	results := CheckLatestTag(df)
	if len(results) != 0 {
		t.Errorf("expected no results for scratch, got %v", results)
	}
}

func TestCheckNoUser_Missing(t *testing.T) {
	df := makeDF(t, "FROM ubuntu:22.04\nRUN echo hi\n")
	results := CheckNoUser(df)
	if len(results) != 1 || results[0].Severity != "HIGH" {
		t.Errorf("expected 1 HIGH result for missing USER, got %v", results)
	}
}

func TestCheckNoUser_NonRoot(t *testing.T) {
	df := makeDF(t, "FROM ubuntu:22.04\nUSER 1001\n")
	results := CheckNoUser(df)
	if len(results) != 0 {
		t.Errorf("expected no results with non-root USER, got %v", results)
	}
}

func TestCheckNoUser_ExplicitRoot(t *testing.T) {
	df := makeDF(t, "FROM ubuntu:22.04\nUSER root\n")
	results := CheckNoUser(df)
	if len(results) != 1 {
		t.Errorf("expected 1 result for explicit root USER, got %d", len(results))
	}
}

func TestCheckNoHealthcheck_Missing(t *testing.T) {
	df := makeDF(t, "FROM ubuntu:22.04\nRUN echo hi\n")
	results := CheckNoHealthcheck(df)
	if len(results) != 1 || results[0].Severity != "LOW" {
		t.Errorf("expected 1 LOW result for missing HEALTHCHECK, got %v", results)
	}
}

func TestCheckNoHealthcheck_Present(t *testing.T) {
	df := makeDF(t, "FROM ubuntu:22.04\nHEALTHCHECK CMD curl -f http://localhost/\n")
	results := CheckNoHealthcheck(df)
	if len(results) != 0 {
		t.Errorf("expected no results with HEALTHCHECK, got %v", results)
	}
}

func TestCheckNoMultiStage_Single(t *testing.T) {
	df := makeDF(t, "FROM ubuntu:22.04\nRUN echo hi\n")
	results := CheckNoMultiStage(df)
	if len(results) != 1 {
		t.Errorf("expected 1 result for single-stage build, got %d", len(results))
	}
}

func TestCheckNoMultiStage_Multi(t *testing.T) {
	df := makeDF(t, "FROM golang:1.21 AS builder\nFROM alpine:3.18\n")
	results := CheckNoMultiStage(df)
	if len(results) != 0 {
		t.Errorf("expected no results for multi-stage build, got %v", results)
	}
}

func TestCheckRiskyExpose_SSH(t *testing.T) {
	df := makeDF(t, "FROM ubuntu:22.04\nEXPOSE 22\n")
	results := CheckRiskyExpose(df)
	if len(results) != 1 || results[0].Severity != "MEDIUM" {
		t.Errorf("expected 1 MEDIUM result for SSH port, got %v", results)
	}
}

func TestCheckRiskyExpose_Safe(t *testing.T) {
	df := makeDF(t, "FROM ubuntu:22.04\nEXPOSE 8080\n")
	results := CheckRiskyExpose(df)
	if len(results) != 0 {
		t.Errorf("expected no results for port 8080, got %v", results)
	}
}

func TestCheckRiskyExpose_WithProtocol(t *testing.T) {
	df := makeDF(t, "FROM ubuntu:22.04\nEXPOSE 22/tcp\n")
	results := CheckRiskyExpose(df)
	if len(results) != 1 {
		t.Errorf("expected 1 result for 22/tcp, got %d", len(results))
	}
}
