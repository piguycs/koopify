package internal

import (
	"github.com/resend/resend-go/v3"
)

type EmailHandler struct {
	client *resend.Client
}

func NewEmailHandler(resendApiKey string) EmailHandler {
	client := resend.NewClient(resendApiKey)

	return EmailHandler{
		client: client,
	}
}

func (ec *EmailHandler) SendEmail(from, to, subject, body string) (*resend.SendEmailResponse, error) {
	params := resend.SendEmailRequest{
		From:    from,
		To:      []string{to},
		Text:    body,
		Subject: subject,
	}

	sent, err := ec.client.Emails.Send(&params)
	if err != nil {
		return nil, err
	}

	return sent, nil
}
