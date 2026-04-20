// Package adapter demonstrates the Adapter pattern in Go.
//
// An adapter wraps an incompatible type and translates its interface into the
// one your application expects. Neither the target interface nor the adaptee
// need to change.
//
// Real-world fit: you depend on a third-party SDK or a legacy system whose
// method signatures don't match your domain interface. Instead of forking it
// or scattering type-conversion logic across callers, you write one adapter.
package adapter

import (
	"errors"
	"fmt"
)

// CloudStorage is the interface our application works with.
// All business logic depends only on this — it never imports the legacy package.
type CloudStorage interface {
	Upload(key, data string) error
	Download(key string) (string, error)
	Remove(key string) error
}

// --- Adaptee (third-party / legacy) ---

// LegacyFTPClient simulates a third-party FTP library with an incompatible API:
// different method names, a combined error signal via a bool, and no Remove.
type LegacyFTPClient struct {
	files map[string]string
}

func NewLegacyFTPClient() *LegacyFTPClient {
	return &LegacyFTPClient{files: make(map[string]string)}
}

// PutFile stores a file. Returns false on failure — no error type.
func (f *LegacyFTPClient) PutFile(filename, content string) bool {
	if filename == "" {
		return false
	}
	f.files[filename] = content
	fmt.Printf("[FTP] PUT %s\n", filename)
	return true
}

// FetchFile retrieves a file. Returns ("", false) when not found.
func (f *LegacyFTPClient) FetchFile(filename string) (string, bool) {
	content, ok := f.files[filename]
	if ok {
		fmt.Printf("[FTP] FETCH %s\n", filename)
	}
	return content, ok
}

// DeleteFile removes a file. No return value.
func (f *LegacyFTPClient) DeleteFile(filename string) {
	delete(f.files, filename)
	fmt.Printf("[FTP] DELETE %s\n", filename)
}

// --- Adapter ---

// FTPAdapter translates CloudStorage calls into LegacyFTPClient calls.
// It is the only place that knows about both interfaces.
type FTPAdapter struct {
	client *LegacyFTPClient
}

func NewFTPAdapter(client *LegacyFTPClient) *FTPAdapter {
	return &FTPAdapter{client: client}
}

func (a *FTPAdapter) Upload(key, data string) error {
	if ok := a.client.PutFile(key, data); !ok {
		return errors.New("ftp: upload failed")
	}
	return nil
}

func (a *FTPAdapter) Download(key string) (string, error) {
	content, ok := a.client.FetchFile(key)
	if !ok {
		return "", fmt.Errorf("ftp: file %q not found", key)
	}
	return content, nil
}

func (a *FTPAdapter) Remove(key string) error {
	// LegacyFTPClient.DeleteFile has no return value — adapter absorbs that gap.
	a.client.DeleteFile(key)
	return nil
}
