package api

import "time"

type User struct {
	Nickname       string
	FullName       string
	Email          string
	Gender         string
	Birthday       time.Time
	Program        string
	Sequence       string
	GraduatingYear int
}
