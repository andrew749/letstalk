package main

import (
	"fmt"
	"sort"

	"letstalk/server/core/api"
	"letstalk/server/core/query"
	"letstalk/server/core/survey"
	"letstalk/server/data"
	"letstalk/server/utility"

	"github.com/namsral/flag"
	"github.com/romana/rlog"
)

var (
	emailsFilename = flag.String("emails_file", "", "File containing emails sperated by newlines")
	mentorGradYear = flag.Int("mentor_grad_year", -1, "Youngest graduating class for mentor")
)

const (
	MAX_MENTEES_PER_MENTOR = 2
	EMAIL_FIELD            = "email"
)

// Also removes duplicates and users without cohorts
// Returns (mentor surveys, mentee surveys)
func separateMentorMentees(users []data.User, mentorGradYear uint) ([]data.User, []data.User) {
	usersById := make(map[data.TUserID]data.User)
	for _, user := range users {
		if user.Cohort != nil && user.Cohort.Cohort != nil {
			usersById[user.UserId] = user
		}
	}
	mentorUsers := make([]data.User, 0)
	menteeUsers := make([]data.User, 0)
	for _, userSurvey := range usersById {
		if userSurvey.Cohort.Cohort.GradYear <= mentorGradYear {
			mentorUsers = append(mentorUsers, userSurvey)
		} else {
			menteeUsers = append(menteeUsers, userSurvey)
		}
	}
	return mentorUsers, menteeUsers
}

type versionedSurveyGroup struct {
	group   data.SurveyGroup
	version int
}

type userWithSurveys struct {
	user    data.User
	surveys map[versionedSurveyGroup]data.UserSurvey
}

// Groups user surveys by survey group and version (only for currently supported surveys)
func groupUserSurveys(users []data.User, surveys []api.Survey) []userWithSurveys {
	withSurveys := make([]userWithSurveys, len(users))

	for i, user := range users {
		withSurveys[i] = userWithSurveys{
			user:    user,
			surveys: make(map[versionedSurveyGroup]data.UserSurvey),
		}

		for _, sur := range surveys {
			for _, userSurvey := range user.UserSurveys {
				if userSurvey.Group == sur.Group && userSurvey.Version == sur.Version {
					withSurveys[i].surveys[versionedSurveyGroup{sur.Group, sur.Version}] = userSurvey
				}
			}
		}
	}
	return withSurveys
}

type surveyAlgoMatch struct {
	mentor          data.User
	mentee          data.User
	matchingAnswers uint
}

type byMatchingAnswers []surveyAlgoMatch

func (a byMatchingAnswers) Len() int {
	return len(a)
}

func (a byMatchingAnswers) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a byMatchingAnswers) Less(i, j int) bool {
	// Sort by decreasing
	return a[i].matchingAnswers > a[j].matchingAnswers
}

func numMatchingSurveyResponse(
	theSurvey api.Survey,
	userSurveyOne data.UserSurvey,
	userSurveyTwo data.UserSurvey,
) uint {
	matchingAnswers := uint(0)
	for _, question := range theSurvey.Questions {
		mentorAnswer, mentorOk := userSurveyOne.Responses[question.Key]
		menteeAnswer, menteeOk := userSurveyTwo.Responses[question.Key]
		if menteeOk && mentorOk && mentorAnswer == menteeAnswer {
			matchingAnswers++
		}
	}
	return matchingAnswers
}

func computeMatches(
	surveys []api.Survey,
	mentors []userWithSurveys,
	mentees []userWithSurveys,
) []surveyAlgoMatch {
	allMatches := make([]surveyAlgoMatch, 0)
	// Generate all pairwise matchings of mentors and mentees - calculate overlapping answers
	for _, mentor := range mentors {
		for _, mentee := range mentees {
			matchingAnswers := uint(0)
			for _, sur := range surveys {
				surveyKey := versionedSurveyGroup{sur.Group, sur.Version}
				survey1, ok1 := mentor.surveys[surveyKey]
				survey2, ok2 := mentee.surveys[surveyKey]
				if ok1 && ok2 {
					matchingAnswers += numMatchingSurveyResponse(sur, survey1, survey2)
				}
			}
			allMatches = append(allMatches, surveyAlgoMatch{mentor.user, mentee.user, matchingAnswers})
		}
	}
	sort.Sort(byMatchingAnswers(allMatches))

	hasMatch := make(map[data.TUserID]interface{})
	menteeCount := make(map[data.TUserID]uint)
	for _, mentor := range mentors {
		menteeCount[mentor.user.UserId] = 0
	}

	matches := make([]surveyAlgoMatch, 0)
	for _, match := range allMatches {
		if len(hasMatch) == len(mentees) {
			// All mentees already have a mentor so we are done
			rlog.Debugf("Found all %d matches, breaking...", len(hasMatch))
			break
		}
		if _, ok := hasMatch[match.mentee.UserId]; ok {
			// Mentee already matched so we continue
			rlog.Debugf("Mentee %d already has a match, continuing...", match.mentee.UserId)
			continue
		}
		if count := menteeCount[match.mentor.UserId]; count >= MAX_MENTEES_PER_MENTOR {
			// Mentor already has MAX_MENTEES_PER_MENTOR mentees so don't assign another one
			rlog.Debugf(
				"Mentor %d already has %d matches, continuing...",
				match.mentor.UserId,
				MAX_MENTEES_PER_MENTOR,
			)
			continue
		}
		rlog.Debugf(
			"Adding match: mentor: %d, mentee: %d, count: %d",
			match.mentor.UserId,
			match.mentee.UserId,
			match.matchingAnswers,
		)
		menteeCount[match.mentor.UserId]++
		hasMatch[match.mentee.UserId] = nil
		matches = append(matches, match)
	}

	rlog.Debugf("Found %d matches", len(matches))
	return matches
}

func main() {
	flag.Parse()

	if mentorGradYear == nil || *mentorGradYear < 0 {
		panic("Must provide -mentor_grad_year")
	}

	emails, err := utility.GetEmailsFromFile(*emailsFilename)
	if err != nil {
		panic(err)
	}

	db, err := utility.GetDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	surveys := survey.GetAllSurveyDefinitions()
	users, err := query.GetUsersWithCohortAndSurveysByEmail(db, emails)
	if err != nil {
		panic(err)
	}
	mentors, mentees := separateMentorMentees(users, uint(*mentorGradYear))
	rlog.Debugf("Num mentors: %d", len(mentors))
	rlog.Debugf("Num mentees: %d", len(mentees))

	mentorsWithSurveys := groupUserSurveys(mentors, surveys)
	menteesWithSurveys := groupUserSurveys(mentees, surveys)

	matches := computeMatches(surveys, mentorsWithSurveys, menteesWithSurveys)
	fmt.Println("mentor_first_name,mentor_last_name,mentor_email,mentee_first_name," +
		"mentee_last_name,mentee_email,num_matches")
	for _, match := range matches {
		fmt.Printf("%s,%s,%s,%s,%s,%s,%d\n",
			match.mentor.FirstName,
			match.mentor.LastName,
			match.mentor.Email,
			match.mentee.FirstName,
			match.mentee.LastName,
			match.mentee.Email,
			match.matchingAnswers,
		)
	}
}
