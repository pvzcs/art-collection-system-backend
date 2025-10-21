package utils

import (
	"art-collection-system/internal/config"
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

// EmailService handles email sending operations
type EmailService struct {
	config *config.EmailConfig
}

// NewEmailService creates a new email service instance
func NewEmailService(cfg *config.EmailConfig) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

// SendVerificationCode sends a verification code email
func (s *EmailService) SendVerificationCode(to, code string) error {
	subject := "美术作品投稿系统 - 邮箱验证码"
	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 5px;">
				<h2 style="color: #4CAF50;">美术作品投稿系统</h2>
				<p>您好，</p>
				<p>您正在注册美术作品投稿系统账号，您的验证码是：</p>
				<div style="background-color: #f4f4f4; padding: 15px; text-align: center; font-size: 24px; font-weight: bold; letter-spacing: 5px; margin: 20px 0;">
					%s
				</div>
				<p style="color: #666;">验证码有效期为 <strong>5分钟</strong>，请尽快完成验证。</p>
				<p style="color: #999; font-size: 12px; margin-top: 30px;">
					如果这不是您的操作，请忽略此邮件。
				</p>
			</div>
		</body>
		</html>
	`, code)

	return s.sendEmail(to, subject, body)
}

// sendEmail sends an email using SMTP
func (s *EmailService) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		s.config.SMTPHost,
		s.config.SMTPPort,
		s.config.Username,
		s.config.Password,
	)

	// Use TLS for port 465 (SMTPS)
	if s.config.SMTPPort == 465 {
		d.SSL = true
	} else {
		// Use STARTTLS for other ports (like 587)
		d.TLSConfig = &tls.Config{InsecureSkipVerify: false}
	}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
