package main

import (
	"fmt"
	"sort"

	"letstalk/server/core/api"
	"letstalk/server/core/survey"
	"letstalk/server/data"
	"letstalk/server/utility"

	"github.com/jinzhu/gorm"
	"github.com/namsral/flag"
)

var (
	outFile        = flag.String("out", "", "Output csv file with matchings")
	surveyGroupStr = flag.String("survey_group", "", "Survey group to do matchings based on")
	surveyVersion  = flag.Int("survey_version", -1, "Version of survey")
	mentorGradYear = flag.Int("mentor_grad_year", -1, "Youngest graduating class for mentor")
)

func getUserSurveys(
	db *gorm.DB,
	surveyGroup data.SurveyGroup,
	surveyVersion int,
) ([]data.UserSurvey, error) {
	var userSurveys []data.UserSurvey
	if err := db.Where(
		&data.UserSurvey{Group: surveyGroup, Version: surveyVersion},
	).Preload("User.Cohort.Cohort").Find(&userSurveys).Error; err != nil {
		return nil, err
	}
	return userSurveys, nil
}

// Also removes duplicates and users without cohorts
// Returns (mentor surveys, mentee surveys)
func separateMentorMentees(
	userSurveys []data.UserSurvey,
	mentorGradYear uint,
) ([]data.UserSurvey, []data.UserSurvey) {
	userSurveysById := make(map[data.TUserID]data.UserSurvey)
	for _, userSurvey := range userSurveys {
		if userSurvey.User != nil &&
			userSurvey.User.Cohort != nil &&
			userSurvey.User.Cohort.Cohort != nil {
			userSurveysById[userSurvey.UserId] = userSurvey
		}
	}
	mentorUserSurveys := make([]data.UserSurvey, 0)
	menteeUserSurveys := make([]data.UserSurvey, 0)
	for _, userSurvey := range userSurveysById {
		if userSurvey.User.Cohort.Cohort.GradYear <= mentorGradYear {
			mentorUserSurveys = append(mentorUserSurveys, userSurvey)
		} else {
			menteeUserSurveys = append(menteeUserSurveys, userSurvey)
		}
	}
	return mentorUserSurveys, menteeUserSurveys
}

type withGenericSurvey struct {
	data.UserSurvey
	genericSurvey *data.UserSurvey
}

func getGenericSurveys(
	db *gorm.DB,
	genericSurvey api.Survey,
	userSurveys []data.UserSurvey,
) ([]withGenericSurvey, error) {
	withGeneric := make([]withGenericSurvey, len(userSurveys))
	for i, userSurvey := range userSurveys {
		genericUserSurvey := &data.UserSurvey{}
		res := db.Where(&data.UserSurvey{
			UserId:  userSurvey.UserId,
			Group:   genericSurvey.Group,
			Version: genericSurvey.Version,
		}).Find(genericUserSurvey)
		if res.RecordNotFound() {
			genericUserSurvey = nil
		} else if res.Error != nil {
			return nil, res.Error
		}
		withGeneric[i] = withGenericSurvey{userSurvey, genericUserSurvey}
	}
	return withGeneric, nil
}

type surveyAlgoMatch struct {
	mentorSurvey    withGenericSurvey
	menteeSurvey    withGenericSurvey
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

const MAX_MENTEES_PER_MENTOR = 2

func numMatches(
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
	theSurvey api.Survey,
	genericSurvey api.Survey,
	mentorUserSurveys []withGenericSurvey,
	menteeUserSurveys []withGenericSurvey,
) []surveyAlgoMatch {
	allMatches := make([]surveyAlgoMatch, 0)
	for _, mentorSurvey := range mentorUserSurveys {
		for _, menteeSurvey := range menteeUserSurveys {
			matchingAnswers := numMatches(theSurvey, mentorSurvey.UserSurvey, menteeSurvey.UserSurvey)
			if mentorSurvey.genericSurvey != nil && menteeSurvey.genericSurvey != nil {
				matchingAnswers += numMatches(
					genericSurvey,
					*mentorSurvey.genericSurvey,
					*menteeSurvey.genericSurvey,
				)
			}
			allMatches = append(allMatches, surveyAlgoMatch{mentorSurvey, menteeSurvey, matchingAnswers})
		}
	}
	sort.Sort(byMatchingAnswers(allMatches))

	hasMatch := make(map[data.TUserID]interface{})
	menteeCount := make(map[data.TUserID]uint)
	for _, mentorSurvey := range mentorUserSurveys {
		menteeCount[mentorSurvey.UserId] = 0
	}

	matches := make([]surveyAlgoMatch, 0)
	for _, match := range allMatches {
		if len(hasMatch) == len(menteeUserSurveys) {
			// All mentees already have a mentor so we are done
			break
		}
		if _, ok := hasMatch[match.menteeSurvey.UserId]; ok {
			// Mentee already matched so we continue
			continue
		}
		if count := menteeCount[match.mentorSurvey.UserId]; count >= MAX_MENTEES_PER_MENTOR {
			// Mentor already has MAX_MENTEES_PER_MENTOR mentees so don't assign another one
			continue
		}
		menteeCount[match.mentorSurvey.UserId]++
		hasMatch[match.menteeSurvey.UserId] = nil
		matches = append(matches, match)
	}

	return matches
}

func main() {
	flag.Parse()

	if surveyGroupStr == nil || *surveyGroupStr == "" {
		panic("Must provide -survey_group")
	}
	if mentorGradYear == nil || *mentorGradYear < 0 {
		panic("Must provide -mentor_grad_year")
	}
	if surveyVersion == nil || *surveyVersion < 0 {
		panic("Must provide -survey_version")
	}
	surveyGroup := data.SurveyGroup(*surveyGroupStr)
	theSurvey := survey.GetSurveyDefinitionByGroup(surveyGroup)
	if theSurvey == nil {
		panic(fmt.Sprintf("Cannot find survey %s", surveyGroup))
	} else if theSurvey.Version != *surveyVersion {
		panic(fmt.Sprintf("Invalid survey version %d. Expected %d", *surveyVersion, theSurvey.Version))
	}
	genericSurvey := survey.GetSurveyDefinitionByGroup("GENERIC")
	if genericSurvey == nil {
		panic("Must have a GENERIC survey")
	}

	db, err := utility.GetDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userSurveys, err := getUserSurveys(db, surveyGroup, *surveyVersion)
	if err != nil {
		panic(err)
	}

	mentorUserSurveys, menteeUserSurveys := separateMentorMentees(userSurveys, uint(*mentorGradYear))
	fmt.Printf("Num mentors: %d\n", len(mentorUserSurveys))
	fmt.Printf("Num mentees: %d\n", len(menteeUserSurveys))

	mentorUserSurveysWithGeneric, err := getGenericSurveys(db, *genericSurvey, mentorUserSurveys)
	if err != nil {
		panic(err)
	}
	menteeUserSurveysWithGeneric, err := getGenericSurveys(db, *genericSurvey, menteeUserSurveys)
	if err != nil {
		panic(err)
	}

	matches := computeMatches(
		*theSurvey,
		*genericSurvey,
		mentorUserSurveysWithGeneric,
		menteeUserSurveysWithGeneric,
	)
	fmt.Println("mentor_first_name,mentor_last_name,mentor_email,mentee_first_name," +
		"mentee_last_name,mentee_email,num_matches")
	for _, match := range matches {
		fmt.Printf("%s,%s,%s,%s,%s,%s,%d\n",
			match.mentorSurvey.User.FirstName,
			match.mentorSurvey.User.LastName,
			match.mentorSurvey.User.Email,
			match.menteeSurvey.User.FirstName,
			match.menteeSurvey.User.LastName,
			match.menteeSurvey.User.Email,
			match.matchingAnswers,
		)
	}
}
