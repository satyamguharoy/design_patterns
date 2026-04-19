// Package dependency_injection demonstrates the Dependency Injection pattern in Go.
//
// WITHOUT DI: a service creates its own dependencies — tightly coupled, untestable.
//
//	type BadNotificationService struct{}
//	func (s *BadNotificationService) Notify(msg string) {
//	    sender := smtp.NewSMTPSender("smtp.example.com") // hardcoded — can't swap in tests
//	    sender.Send(msg)
//	}
//
// WITH DI: dependencies are declared as interfaces and injected via the constructor.
// The service only knows about the interface, not the concrete type.
package dependency_injection

import "fmt"

// MessageSender abstracts the transport layer.
type MessageSender interface {
	Send(to, message string) error
}

// Logger abstracts logging so the service isn't tied to a specific logger.
type Logger interface {
	Log(msg string)
}

// NotificationService sends notifications using injected dependencies.
type NotificationService struct {
	sender MessageSender
	logger Logger
}

// NewNotificationService is the constructor — dependencies are injected here.
func NewNotificationService(sender MessageSender, logger Logger) *NotificationService {
	return &NotificationService{sender: sender, logger: logger}
}

func (s *NotificationService) Notify(to, message string) error {
	s.logger.Log(fmt.Sprintf("sending notification to %s", to))
	if err := s.sender.Send(to, message); err != nil {
		s.logger.Log(fmt.Sprintf("failed to send: %v", err))
		return err
	}
	s.logger.Log("notification sent successfully")
	return nil
}

// SMTPSender is a real production implementation of MessageSender.
type SMTPSender struct {
	host string
}

func NewSMTPSender(host string) *SMTPSender {
	return &SMTPSender{host: host}
}

func (s *SMTPSender) Send(to, message string) error {
	// real SMTP logic would go here
	fmt.Printf("[SMTP via %s] To: %s | %s\n", s.host, to, message)
	return nil
}

// ConsoleLogger is a simple production logger that writes to stdout.
type ConsoleLogger struct{}

func (l *ConsoleLogger) Log(msg string) {
	fmt.Println("[LOG]", msg)
}
