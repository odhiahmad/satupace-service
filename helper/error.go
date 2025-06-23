package helper

import (
	"log"
)

// ErrorPanic akan memunculkan panic jika terjadi error, serta mencatatnya ke log
func ErrorPanic(err error) {
	if err != nil {
		log.Printf("Fatal error: %v\n", err)
		panic(err)
	}
}
