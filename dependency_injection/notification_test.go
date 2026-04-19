package dependency_injection_test

import (
	"errors"
	"testing"

	di "design_patterns/dependency_injection"
)

// mockSender and mockLogger are test doubles — only possible because we injected interfaces.

type mockSender struct {
	sentTo  string
	sentMsg string
	err     error
}

func (m *mockSender) Send(to, message string) error {
	m.sentTo = to
	m.sentMsg = message
	return m.err
}

type mockLogger struct {
	logs []string
}

func (l *mockLogger) Log(msg string) {
	l.logs = append(l.logs, msg)
}

func TestNotify_Success(t *testing.T) {
	sender := &mockSender{}
	logger := &mockLogger{}
	svc := di.NewNotificationService(sender, logger)

	if err := svc.Notify("alice@example.com", "hello"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sender.sentTo != "alice@example.com" {
		t.Errorf("expected sentTo alice@example.com, got %s", sender.sentTo)
	}
	if sender.sentMsg != "hello" {
		t.Errorf("expected sentMsg hello, got %s", sender.sentMsg)
	}
}

func TestNotify_SenderError(t *testing.T) {
	sender := &mockSender{err: errors.New("connection refused")}
	logger := &mockLogger{}
	svc := di.NewNotificationService(sender, logger)

	if err := svc.Notify("bob@example.com", "hi"); err == nil {
		t.Fatal("expected error, got nil")
	}

	// verify the failure was logged
	logged := false
	for _, l := range logger.logs {
		if l == "failed to send: connection refused" {
			logged = true
		}
	}
	if !logged {
		t.Error("expected failure to be logged")
	}
}
