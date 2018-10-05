package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"letstalk/server/core/api"

	"golang.org/x/crypto/ssh/terminal"
	"letstalk/server/core/query"
)

const mentorField = "Mentor Email"
const menteeField = "Mentee Email"

func printUsage() {
	fmt.Println("Usage: go run add_mentorships_csv.go -file=mentorships.csv [-server=server]")
}

// Return session token.
func login(server string) string {
	// Get login credentials
	var email string
	var password []byte
	var err error
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("email: ")
	if email, err = reader.ReadString('\n'); err != nil {
		panic(err)
	}
	fmt.Print("pwd:")
	if password, err = terminal.ReadPassword(0); err != nil {
		panic(err)
	}
	// Log in with credentials.
	loginRequest := api.LoginRequestData{
		Email: strings.TrimSpace(email),
		Password: strings.TrimSpace(string(password)),
		NotificationToken: nil,
	}
	var buf = &bytes.Buffer{}
	json.NewEncoder(buf).Encode(loginRequest)
	res, err := http.Post(fmt.Sprintf("%s/v1/login", server), "application/json", buf)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != http.StatusOK {
		res.Body.Read(buf.Bytes())
		fmt.Println("failed to log in", res.Status, buf.String())
		os.Exit(1)
	}
	type loginResult struct {
		Result api.LoginResponse
	}
	loginResponse := loginResult{}
	json.NewDecoder(res.Body).Decode(&loginResponse)
	return loginResponse.Result.SessionId
}

func addMentorship(server string, sessionId string, mentorEmail string, menteeEmail string) error {
	createRequest := api.CreateMentorshipByEmail{MentorEmail: mentorEmail, MenteeEmail: menteeEmail}
	var buf = &bytes.Buffer{}
	json.NewEncoder(buf).Encode(createRequest)
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/admin/mentorship", server), buf)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("sessionId", sessionId)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		var err struct {Error query.Error}
		json.NewDecoder(response.Body).Decode(&err)
		return errors.New(fmt.Sprintf("%v", err))
	}
	return nil
}

func main() {
	fileFlag := flag.String("file", "", "csv file to parse")
	serverFlag := flag.String("server", "", "hive server")

	// Once all flags are declared, call `flag.Parse()`
	// to execute the command-line parsing.
	flag.Parse()

	if len(*fileFlag) == 0 {
		printUsage()
		os.Exit(1)
	}

	var server, sessionId string
	if len(*serverFlag) > 0 {
		server = strings.TrimSpace(*serverFlag)
		sessionId = login(server)
	}

	file, err := os.Open(*fileFlag)
	if err != nil {
		fmt.Println("no such file", *fileFlag, err)
		os.Exit(1)
	}
	fullCsv, err := csv.NewReader(bufio.NewReader(file)).ReadAll()
	if err != nil {
		fmt.Println("failed to parse file", err)
		os.Exit(1)
	}
	fields := make(map[string]int)
	for i, field := range fullCsv[0] {
		fields[field] = i
	}
	for _, line := range fullCsv[1:] {
		mentorEmail := line[fields[mentorField]]
		menteeEmail := line[fields[menteeField]]
		fmt.Printf("mentor '%s', mentee '%s'\n", mentorEmail, menteeEmail)
		if len(*serverFlag) > 0 {
			if err := addMentorship(server, sessionId, mentorEmail, menteeEmail); err != nil {
				fmt.Println("failed to add", err)
			}
		}
	}
}


