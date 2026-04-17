package utils

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func writeZip(t *testing.T, path string, entries map[string]string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	for name, content := range entries {
		w, err := zw.Create(name)
		if err != nil {
			_ = zw.Close()
			t.Fatalf("create entry %q: %v", name, err)
		}
		if _, err := w.Write([]byte(content)); err != nil {
			_ = zw.Close()
			t.Fatalf("write entry %q: %v", name, err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
}

func TestUnzip_ExtractsFilesUnderDest(t *testing.T) {
	tmp := t.TempDir()
	zipPath := filepath.Join(tmp, "ok.zip")
	dest := filepath.Join(tmp, "out")
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatalf("mkdir dest: %v", err)
	}

	writeZip(t, zipPath, map[string]string{
		"hello.txt":     "world",
		"dir/nested.md": "ok",
	})

	paths, err := Unzip(zipPath, dest)
	if err != nil {
		t.Fatalf("Unzip: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(paths))
	}
	if _, err := os.Stat(filepath.Join(dest, "hello.txt")); err != nil {
		t.Fatalf("missing extracted file: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "dir", "nested.md")); err != nil {
		t.Fatalf("missing extracted nested file: %v", err)
	}
}

func TestUnzip_RejectsPathTraversal(t *testing.T) {
	tmp := t.TempDir()
	zipPath := filepath.Join(tmp, "traversal.zip")
	dest := filepath.Join(tmp, "out")
	_ = os.MkdirAll(dest, 0o755)

	writeZip(t, zipPath, map[string]string{
		"../evil.txt": "nope",
	})

	_, err := Unzip(zipPath, dest)
	if err == nil {
		t.Fatalf("expected error for path traversal, got nil")
	}
}

func TestUnzip_RejectsAbsolutePath(t *testing.T) {
	tmp := t.TempDir()
	zipPath := filepath.Join(tmp, "abs.zip")
	dest := filepath.Join(tmp, "out")
	_ = os.MkdirAll(dest, 0o755)

	// Zip spec uses forward slashes. This should still be treated as an absolute path on unix.
	writeZip(t, zipPath, map[string]string{
		"/abs.txt": "nope",
	})

	_, err := Unzip(zipPath, dest)
	if err == nil {
		t.Fatalf("expected error for absolute path, got nil")
	}
}

