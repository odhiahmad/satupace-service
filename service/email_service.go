package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type EmailService interface {
	Send(to string, subject string, body string) error
}

type emailService struct {
	from    string
	apiKey  string
	baseURL string
}

func NewEmailService(from, apiKey string) EmailService {
	return &emailService{
		from:    from,
		apiKey:  apiKey,
		baseURL: "https://api.resend.com/emails",
	}
}

func (s *emailService) Send(to string, subject string, body string) error {
	type resendEmailRequest struct {
		From    string   `json:"from"`
		To      []string `json:"to"`
		Subject string   `json:"subject"`
		Html    string   `json:"html"`
	}

	payload := resendEmailRequest{
		From:    s.from,
		To:      []string{to},
		Subject: subject,
		Html:    body,
	}

	log.Printf("📨 [Resend] Preparing to send email to: %s, subject: %s", to, subject)

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		log.Printf("❌ [Resend] Failed to marshal email payload: %v", err)
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	req, err := http.NewRequest("POST", s.baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("❌ [Resend] Failed to create HTTP request: %v", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Set("Content-Type", "application/json")

	log.Printf("🌐 [Resend] Sending email request to %s", s.baseURL)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ [Resend] Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("📬 [Resend] Response Status: %s", resp.Status)

	if resp.StatusCode >= 400 {
		var respBody bytes.Buffer
		_, _ = respBody.ReadFrom(resp.Body)
		log.Printf("❌ [Resend] API Error Response: %s", respBody.String())
		return fmt.Errorf("resend API returned status: %s", resp.Status)
	}

	log.Printf("✅ [Resend] Email sent successfully to: %s", to)
	return nil
}
