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

// Validate checks if the given email string is a valid UW email.
func Validate(email string) bool { return uwEmailRegex.MatchString(email) }

// FromString must be used to create a UwEmail from a raw email address string.
func FromString(raw string) UwEmail {
	if !Validate(raw) {
		return UwEmail{}
	}
	user := strings.SplitN(raw, "@", 2)[0]
	normalized := fmt.Sprintf("%s@edu.uwaterloo.ca", strings.ToLower(user))
	return UwEmail{raw, normalized}
}

// ToStringNormalized returns a UW email normalized to a lower-case "@edu.uwaterloo.ca" address.
func (email *UwEmail) ToStringNormalized() string {
	return email.normalized
}

// ToStringRaw returns the raw email string that was used to create this UwEmail.
func (email *UwEmail) ToStringRaw() string {
	return email.raw
}
