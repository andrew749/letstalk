package utility

import (
	"bufio"
	"os"
)

func GetEmailsFromFile(emailsFilename string) ([]string, error) {
	emailsFile, err := os.Open(emailsFilename)
	if err != nil {
		return nil, err
	}
	defer emailsFile.Close()
	emails := make([]string, 0)
	scanner := bufio.NewScanner(emailsFile)
	for scanner.Scan() {
		emails = append(emails, scanner.Text())
	}
	err = scanner.Err()
	if err != nil {
		return nil, err
	}
	return emails, nil
}
