package uw_email

import (
	"fmt"
	"regexp"
	"strings"
)

type UwEmail struct {
	raw string
	normalized string
}

var uwEmailRegex = regexp.MustCompile("(?i)^[A-Z0-9._%+-]+@(edu\\.)?uwaterloo\\.ca$")

func Validate(email string) bool { return uwEmailRegex.MatchString(email) }

func OfString(raw string) UwEmail {
	if !Validate(raw) {
		return UwEmail{}
	}
	user := strings.SplitN(raw, "@", 2)[0]
	normalized := fmt.Sprintf("%s@edu.uwaterloo.ca", strings.ToLower(user))
	return UwEmail{raw, normalized}
}

func (email *UwEmail) ToStringNormalized() string {
	return email.normalized
}

func (email *UwEmail) ToStringRaw() string {
	return email.raw
}
