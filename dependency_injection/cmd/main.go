package main

import (
	di "design_patterns/dependency_injection"
)

func main() {
	// Wire up real implementations at the top — this is the only place
	// that knows which concrete types are in use.
	sender := di.NewSMTPSender("smtp.example.com")
	logger := &di.ConsoleLogger{}
	svc := di.NewNotificationService(sender, logger)

	svc.Notify("alice@example.com", "Your report is ready.")
}
