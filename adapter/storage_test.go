package adapter_test

import (
	"testing"

	a "design_patterns/adapter"
)

func newAdapter() a.CloudStorage {
	return a.NewFTPAdapter(a.NewLegacyFTPClient())
}

func TestUpload_Success(t *testing.T) {
	store := newAdapter()
	if err := store.Upload("report.csv", "col1,col2"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUpload_EmptyKey(t *testing.T) {
	store := newAdapter()
	if err := store.Upload("", "data"); err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestDownload_Success(t *testing.T) {
	store := newAdapter()
	store.Upload("notes.txt", "hello world")

	content, err := store.Download("notes.txt")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if content != "hello world" {
		t.Errorf("expected 'hello world', got %q", content)
	}
}

func TestDownload_NotFound(t *testing.T) {
	store := newAdapter()
	_, err := store.Download("missing.txt")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestRemove(t *testing.T) {
	store := newAdapter()
	store.Upload("temp.txt", "data")
	store.Remove("temp.txt")

	_, err := store.Download("temp.txt")
	if err == nil {
		t.Fatal("expected error after remove, got nil")
	}
}

// TestAdapterSatisfiesInterface verifies at compile time that FTPAdapter
// implements CloudStorage — the whole point of the pattern.
func TestAdapterSatisfiesInterface(t *testing.T) {
	var _ a.CloudStorage = a.NewFTPAdapter(a.NewLegacyFTPClient())
}
