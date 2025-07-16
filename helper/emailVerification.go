package helper

import (
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

func BuildVerificationEmail(recipientEmail string, token string) (subject string, bodyPlain string, bodyHTML string) {
	subject = "Kode Verifikasi Email Anda - Loka Kasir"

	bodyPlain = fmt.Sprintf(`Halo,

Terima kasih telah mendaftar di Loka Kasir.

Berikut adalah kode verifikasi email Anda:

KODE VERIFIKASI: %s

Silakan masukkan kode ini di aplikasi Loka Kasir untuk menyelesaikan proses verifikasi email Anda.

Jika Anda tidak merasa melakukan pendaftaran, abaikan email ini.

Hormat kami,
Tim Loka Kasir`, token)

	bodyHTML = fmt.Sprintf(`
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <title>Kode Verifikasi</title>
  </head>
  <body style="font-family: sans-serif; line-height: 1.6; color: #333;">
    <p>Halo,</p>
    <p>Terima kasih telah mendaftar di <strong>Loka Kasir</strong>.</p>
    <p>Berikut adalah kode verifikasi email Anda:</p>
    <h2 style="color: #007BFF;">%s</h2>
    <p>Silakan masukkan kode ini di aplikasi Loka Kasir untuk menyelesaikan proses verifikasi email Anda.</p>
    <p>Jika Anda tidak merasa melakukan pendaftaran, abaikan email ini.</p>
    <br />
    <p>Hormat kami,<br />Tim Loka Kasir</p>
  </body>
</html>
`, token)

	return subject, bodyPlain, bodyHTML
}

type EmailHelper struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	FromName string
}

func NewEmailHelper() *EmailHelper {
	return &EmailHelper{
		SMTPHost: "smtp.zoho.com",
		SMTPPort: 587,
		Username: "support@lokakasir.id",
		Password: "Kreativita#123",
		FromName: "Support Loka Kasir <support@lokakasir.id>",
	}
}

func (e *EmailHelper) Send(to string, subject string, bodyPlain string, bodyHTML string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(e.Username, e.FromName))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", bodyPlain)
	m.AddAlternative("text/html", bodyHTML)

	d := gomail.NewDialer(e.SMTPHost, e.SMTPPort, e.Username, e.Password)

	// Optional: Skip TLS verify (not recommended in production)
	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Println("Gagal mengirim email:", err)
		return err
	}

	log.Println("Email terkirim ke:", to)
	return nil
}
