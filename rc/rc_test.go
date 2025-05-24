package rc

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/lonelyday/rsync/config"
)

func TestIsValidPath(t *testing.T) {
	tmpDir := t.TempDir()
	if err := isValidPath(tmpDir); err != nil {
		t.Errorf("expected valid dir, got error: %v", err)
	}
	tmpFile := filepath.Join(tmpDir, "file")
	os.WriteFile(tmpFile, []byte("data"), config.FilePerm)
	if err := isValidPath(tmpFile); err == nil {
		t.Errorf("expected error for file, got nil")
	}
	nonExistent := filepath.Join(tmpDir, "nope")
	if err := isValidPath(nonExistent); err == nil {
		t.Errorf("expected error for non-existent path, got nil")
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "src.txt")
	dst := filepath.Join(tmpDir, "dst.txt")
	content := []byte("hello world")
	if err := os.WriteFile(src, content, config.FilePerm); err != nil {
		t.Fatalf("failed to write src: %v", err)
	}
	fi, err := os.Stat(src)
	if err != nil {
		t.Fatalf("stat src: %v", err)
	}
	if err := copyFile(src, dst, fi); err != nil {
		t.Errorf("copyFile failed: %v", err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read dst: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("copyFile content mismatch: got %q, want %q", got, content)
	}
}

func TestDelMissing(t *testing.T) {
	tmpDir := t.TempDir()
	// Simulate config.DstF
	oldDstF := config.DstF
	defer func() { config.DstF = oldDstF }()
	config.DstF = &tmpDir

	// Create files in dst
	os.MkdirAll(filepath.Join(tmpDir, "keep"), config.FolderPerm)
	os.MkdirAll(filepath.Join(tmpDir, "remove"), config.FolderPerm)
	dst := map[string]bool{"keep": true, "remove": true}
	src := map[string]bool{"keep": true}
	err := delMissing(src, dst)
	if err != nil {
		t.Errorf("delMissing error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "remove")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected 'remove' to be deleted, but it exists")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "keep")); err != nil {
		t.Errorf("expected 'keep' to exist, got error: %v", err)
	}
}

func TestGetPaths(t *testing.T) {
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, "a/b"), config.FolderPerm)
	os.WriteFile(filepath.Join(tmpDir, "a/file.txt"), []byte("x"), config.FilePerm)
	paths, err := getPaths(tmpDir)
	if err != nil {
		t.Fatalf("getPaths error: %v", err)
	}
	want := map[string]bool{
		"a":          true,
		"a/b":        true,
		"a/file.txt": true,
	}
	for k := range want {
		if !paths[k] {
			t.Errorf("expected path %q in result", k)
		}
	}
}
