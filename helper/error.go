package helper

import (
	"errors"
	"log"
)

var (
	ErrInvalidPassword      = errors.New("invalid password")
	ErrMembershipInactive   = errors.New("inactive membership")
	ErrUserNotFound         = errors.New("user not found")
	ErrEmailNotVerified     = errors.New("no hp belum diverifikasi")
	ErrEmailAlreadyVerified = errors.New("no hp sudah diverifikasi")
	ErrProductFailedSaved   = errors.New("Gagal menyimpan variant")
)

// ErrorPanic akan memunculkan panic jika terjadi error, serta mencatatnya ke log
func ErrorPanic(err error) {
	if err != nil {
		log.Printf("Fatal error: %v\n", err)
		panic(err)
	}
}
