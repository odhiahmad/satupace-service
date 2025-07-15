package helper

import "regexp"

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func IsEmail(input string) bool {
	return emailRegex.MatchString(input)
}
