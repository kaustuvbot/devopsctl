package docker

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "Dockerfile")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return path
}

func TestParseDockerfile_Basic(t *testing.T) {
	content := `FROM ubuntu:22.04
RUN apt-get update
USER 1001
EXPOSE 8080
HEALTHCHECK CMD curl -f http://localhost/
`
	path := writeTemp(t, content)
	df, err := ParseDockerfile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(df.Instructions) != 5 {
		t.Errorf("expected 5 instructions, got %d: %v", len(df.Instructions), df.Instructions)
	}
	if df.Instructions[0].Command != "FROM" {
		t.Errorf("expected first instruction FROM, got %s", df.Instructions[0].Command)
	}
	if df.Instructions[2].Command != "USER" {
		t.Errorf("expected third instruction USER, got %s", df.Instructions[2].Command)
	}
}

func TestParseDockerfile_LineContinuation(t *testing.T) {
	content := `FROM ubuntu:22.04
RUN apt-get update && \
    apt-get install -y curl && \
    rm -rf /var/lib/apt/lists/*
`
	path := writeTemp(t, content)
	df, err := ParseDockerfile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(df.Instructions) != 2 {
		t.Errorf("expected 2 instructions (continuation joined), got %d: %v", len(df.Instructions), df.Instructions)
	}
	if df.Instructions[1].Command != "RUN" {
		t.Errorf("expected RUN, got %s", df.Instructions[1].Command)
	}
}

func TestParseDockerfile_Comments(t *testing.T) {
	content := `# Base image
FROM alpine:3.18
# Install deps
RUN apk add curl
`
	path := writeTemp(t, content)
	df, err := ParseDockerfile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(df.Instructions) != 2 {
		t.Errorf("expected 2 instructions (comments skipped), got %d", len(df.Instructions))
	}
}

func TestParseDockerfile_CaseInsensitive(t *testing.T) {
	content := "from ubuntu:22.04\nrun echo hi\n"
	path := writeTemp(t, content)
	df, err := ParseDockerfile(path)
	if err != nil {
		t.Fatal(err)
	}
	if df.Instructions[0].Command != "FROM" {
		t.Errorf("expected FROM (uppercased), got %s", df.Instructions[0].Command)
	}
}

func TestParseDockerfile_NotFound(t *testing.T) {
	_, err := ParseDockerfile("/nonexistent/Dockerfile")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
