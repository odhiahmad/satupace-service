package helper

import "fmt"

func BuildVerificationEmail(recipient string, token string) (subject string, body string) {
	subject = "Kode Verifikasi Email Anda - Loka Kasir"
	body = fmt.Sprintf(`Halo %s,

Terima kasih telah mendaftar di Loka Kasir.

Berikut adalah kode verifikasi email Anda:

🔐 KODE VERIFIKASI: %s

Silakan masukkan kode ini di aplikasi Loka Kasir untuk menyelesaikan proses verifikasi email Anda.

Jika Anda tidak merasa melakukan pendaftaran, abaikan email ini. Kami menjaga keamanan data Anda.

Hormat kami,
Tim Loka Kasir
`, recipient, token)

	return subject, body
}
