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
	"letstalk/server/core/query"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	mentorEmailField     = "mentor_email"
	menteeEmailField     = "mentee_email"
	mentorFirstNameField = "mentor_first_name"
	menteeFirstNameField = "mentee_first_name"
	mentorLastNameField  = "mentor_last_name"
	menteeLastNameField  = "mentee_last_name"
)

func printUsage() {
	fmt.Println("Usage: go run add_mentorships_csv.go -file=mentorships.csv [-server=server] [-not-dry-run]")
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
	fmt.Println()
	// Log in with credentials.
	loginRequest := api.LoginRequestData{
		Email:             strings.TrimSpace(email),
		Password:          strings.TrimSpace(string(password)),
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

func addMentorship(server string, sessionId string, mentorEmail string, menteeEmail string, notDryRun bool) error {
	requestType := api.CREATE_MENTORSHIP_TYPE_DRY_RUN
	if notDryRun {
		requestType = api.CREATE_MENTORSHIP_TYPE_NOT_DRY_RUN
	}
	createRequest := api.CreateMentorshipByEmail{
		MentorEmail: mentorEmail,
		MenteeEmail: menteeEmail,
		RequestType: requestType,
	}
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
		var err struct{ Error query.Error }
		json.NewDecoder(response.Body).Decode(&err)
		return errors.New(fmt.Sprintf("%v", err))
	}
	return nil
}

func main() {
	fileFlag := flag.String("file", "", "csv file to parse")
	serverFlag := flag.String("server", "", "hive server")
	notDryRun := flag.Bool("not-dry-run", false, "whether to actually insert connections in database")
	flag.Parse()

	if len(*fileFlag) == 0 {
		printUsage()
		os.Exit(1)
	}

	if *notDryRun {
		fmt.Println("STARTING ACTUAL (NON-DRY) RUN!!")
	} else {
		fmt.Println("dry run, no changes will be made")
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
	if _, ok := fields[mentorEmailField]; !ok {
		fmt.Printf("header missing mentor email field '%s'\n", mentorEmailField)
		os.Exit(1)
	}
	if _, ok := fields[menteeEmailField]; !ok {
		fmt.Printf("header missing mentee email field '%s'\n", menteeEmailField)
		os.Exit(1)
	}

	for _, line := range fullCsv[1:] {
		mentorEmail := line[fields[mentorEmailField]]
		menteeEmail := line[fields[menteeEmailField]]
		mentorFirstName := line[fields[mentorFirstNameField]]
		menteeFirstName := line[fields[menteeFirstNameField]]
		mentorLastName := line[fields[mentorLastNameField]]
		menteeLastName := line[fields[menteeLastNameField]]
		fmt.Printf("trying mentor: '%s', mentee: '%s'", mentorEmail, menteeEmail)
		if len(*serverFlag) > 0 {
			if err := addMentorship(server, sessionId, mentorEmail, menteeEmail, *notDryRun); err != nil {
				fmt.Printf(
					"\nfailed to add (%s, %s, %s, %s, %s, %s): %v\n",
					mentorFirstName,
					mentorLastName,
					mentorEmail,
					menteeFirstName,
					menteeLastName,
					menteeEmail,
					err,
				)
			} else {
				fmt.Println(" ok!")
			}
		} else {
			fmt.Println()
		}
	}
}
