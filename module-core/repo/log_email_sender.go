package repo

import (
	"context"
	"fmt"

	"github.com/bing127/enterprise-agent/module-core/domain/repo"
)

// LogEmailSender is a stub EmailSender that prints emails to stdout.
// Replace with a real SMTP/SES implementation in production.
type LogEmailSender struct{}

// NewLogEmailSender creates a new LogEmailSender.
func NewLogEmailSender() repo.EmailSender {
	return &LogEmailSender{}
}

func (s *LogEmailSender) Send(_ context.Context, to, subject, body string) error {
	fmt.Printf("[EMAIL] To: %s | Subject: %s | Body: %s\n", to, subject, body)
	return nil
}
