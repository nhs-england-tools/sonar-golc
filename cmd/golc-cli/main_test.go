package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI_CurrentDir(t *testing.T) {
	binary := buildCLI(t)

	cmd := exec.Command(binary, ".")
	cmd.Dir = t.TempDir()

	// Create a sample Go file in the temp dir
	goFile := filepath.Join(cmd.Dir, "hello.go")
	if err := os.WriteFile(goFile, []byte("package main\n\nfunc main() {}\n"), 0644); err != nil {
		t.Fatalf("failed to write sample file: %v", err)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI exited with error: %v\nOutput: %s", err, out)
	}

	if !strings.Contains(string(out), "Golang") {
		t.Errorf("expected output to contain 'Golang', got:\n%s", out)
	}
	if !strings.Contains(string(out), "Total") {
		t.Errorf("expected output to contain 'Total', got:\n%s", out)
	}
}

func TestCLI_NonExistentDir(t *testing.T) {
	binary := buildCLI(t)

	cmd := exec.Command(binary, "/nonexistent-path-xyz")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit for nonexistent directory")
	}
	if !strings.Contains(string(out), "Error") {
		t.Errorf("expected error message, got:\n%s", out)
	}
}

func TestCLI_DefaultsToCurrentDir(t *testing.T) {
	binary := buildCLI(t)

	dir := t.TempDir()
	pyFile := filepath.Join(dir, "app.py")
	if err := os.WriteFile(pyFile, []byte("# comment\nprint('hello')\n"), 0644); err != nil {
		t.Fatalf("failed to write sample file: %v", err)
	}

	cmd := exec.Command(binary)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI exited with error: %v\nOutput: %s", err, out)
	}

	if !strings.Contains(string(out), "Python") {
		t.Errorf("expected output to contain 'Python', got:\n%s", out)
	}
}

func buildCLI(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "golc-cli-test")
	cmd := exec.Command("go", "build", "-o", binary, ".")
	cmd.Dir = filepath.Join(projectRoot(t), "cmd", "golc-cli")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build CLI: %v\n%s", err, out)
	}
	return binary
}

func projectRoot(t *testing.T) string {
	t.Helper()
	// Walk up from this test file's location to find go.mod
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root (go.mod)")
		}
		dir = parent
	}
}
