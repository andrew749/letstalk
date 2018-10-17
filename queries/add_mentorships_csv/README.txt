# Add Mentorships CSV

Go executable to add new mentorship connections between existing users via
csv file. By default, performs a "dry run" which does not make actual database
changes and only reports errors.

Usage:
$ go run add_mentorships_csv.go -file=mentorships.csv [-server=http://hostname] [-not-dry-run]

## Dependencies

go get golang.org/x/crypto/ssh/terminal
