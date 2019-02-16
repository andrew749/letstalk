package seed_mentorships_job

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/romana/rlog"

	"letstalk/server/core/api"
	"letstalk/server/core/connection"
	"letstalk/server/core/query"
	"letstalk/server/core/verify_link"
	"letstalk/server/data"
	"letstalk/server/jobmine"
	"letstalk/server/jobmine_utility"
	"letstalk/server/recommendations"
)

const SEED_MENTORSHIPS_JOB jobmine.JobType = "SeedMentorshipsJob"

// TaskRecord keys
const (
	MENTEE_USER_ID_METADATA_KEY = "menteeUserId"
	MENTOR_USER_ID_METADATA_KEY = "mentorUserId"
)

// JobRecord keys
const (
	PROGRAM_IDS_METADATA_KEY                    = "programIds"
	YOUNGEST_UPPER_GRAD_YEAR_METADATA_KEY       = "youngestUpperGradYear"
	TERM_START_TIME_METADATA_KEY                = "termStartTime"
	TERM_END_TIME_METADATA_KEY                  = "termEndTime"
	MAX_UPPER_YEARS_PER_LOWER_YEAR_METADATA_KEY = "maxUpperYearsPerLowerYear"
	MAX_LOWER_YEARS_PER_UPPER_YEAR_METADATA_KEY = "maxLowerYearsPerUpperYear"
	IS_DRY_RUN_METADATA_KEY                     = "isDryRun"
)

type userMatch struct {
	menteeId data.TUserID
	mentorId data.TUserID
}

func packageTaskRecordMetadata(match userMatch) map[string]interface{} {
	return map[string]interface{}{
		MENTEE_USER_ID_METADATA_KEY: match.menteeId,
		MENTOR_USER_ID_METADATA_KEY: match.mentorId,
	}
}

func parseUserInfo(taskRecord jobmine.TaskRecord) (*userMatch, error) {
	menteeId, err := jobmine_utility.UserIdFromTaskRecord(taskRecord, MENTEE_USER_ID_METADATA_KEY)
	if err != nil {
		return nil, err
	} else if menteeId == nil {
		return nil, errors.New("menteeId not provided in task record")
	}
	mentorId, err := jobmine_utility.UserIdFromTaskRecord(taskRecord, MENTOR_USER_ID_METADATA_KEY)
	if err != nil {
		return nil, err
	} else if mentorId == nil {
		return nil, errors.New("mentorId not provided in task record")
	}
	return &userMatch{menteeId: *menteeId, mentorId: *mentorId}, nil
}

func userStr(db *gorm.DB, user data.User, userStr string) string {
	if user.Cohort != nil && user.Cohort.Cohort != nil {
		return fmt.Sprintf("%s(%d, %d) - program(%s, %d)", userStr, user.UserId, user.Gender,
			user.Cohort.Cohort.ProgramId, user.Cohort.Cohort.GradYear)
	} else {
		return fmt.Sprintf("%s(%d, %d)", userStr, user.UserId, user.Gender)
	}
}

func printMatches(db *gorm.DB, matches []recommendations.UserMatch) error {
	allUserIds := make([]data.TUserID, 0)
	for _, match := range matches {
		allUserIds = append(allUserIds, match.UserOneId, match.UserTwoId)
	}

	var users []data.User
	err := db.Where("user_id IN (?)", allUserIds).Preload("Cohort.Cohort").Find(&users).Error
	if err != nil {
		return err
	}

	usersById := make(map[data.TUserID]data.User)
	for _, user := range users {
		usersById[user.UserId] = user
	}

	for i, match := range matches {
		menteeStr := userStr(db, usersById[match.UserOneId], "mentee")
		mentorStr := userStr(db, usersById[match.UserTwoId], "mentor")

		rlog.Infof(
			"match(%d), %s, %s, score(%f)",
			i, menteeStr, mentorStr, match.Score)
	}
	return nil
}

func execute(
	db *gorm.DB,
	jobRecord jobmine.JobRecord,
	taskRecord jobmine.TaskRecord,
) (interface{}, error) {
	userMatch, err := parseUserInfo(taskRecord)
	if err != nil {
		return nil, err
	}
	err = connection.AddMentorship(
		db, userMatch.mentorId, userMatch.menteeId, api.CREATE_MENTORSHIP_TYPE_NOT_DRY_RUN)
	if err != nil {
		return nil, err
	}

	mentor, err := query.GetUserById(db, userMatch.mentorId)
	if err != nil {
		return nil, err
	}
	mentee, err := query.GetUserById(db, userMatch.menteeId)
	if err != nil {
		return nil, err
	}
	err = connection.SendMentorshipNotifications(db, mentor, mentee)
	if err != nil {
		return nil, err
	}

	return "Success", nil
}

func onError(db *gorm.DB, jobRecord jobmine.JobRecord, taskRecord jobmine.TaskRecord, err error) {
	match, parseErr := parseUserInfo(taskRecord)
	if parseErr != nil {
		rlog.Infof("Unable to create mentorship (%+v): %+v", parseErr, err)
	} else {
		rlog.Infof("Unable to create mentorship for match %v: %+v", *match, err)
	}
}

func onSuccess(
	db *gorm.DB,
	jobRecord jobmine.JobRecord,
	taskRecord jobmine.TaskRecord,
	res interface{},
) {
	match, parseErr := parseUserInfo(taskRecord)
	if parseErr != nil {
		rlog.Infof("Successfully created mentorship (%+v)", parseErr)
	} else {
		rlog.Infof("Successfully created mentorship for match %v", *match)
	}
}

var seedTaskSpec = jobmine.TaskSpec{
	Execute:   execute,
	OnError:   onError,
	OnSuccess: onSuccess,
}

type jobMetadata struct {
	isDryRun                  bool
	programIds                []string
	youngestUpperGradYear     uint
	maxLowerYearsPerUpperYear uint
	maxUpperYearsPerLowerYear uint
	termStartTime             *time.Time
	termEndTime               *time.Time
}

func parseJobMetadata(jobRecord jobmine.JobRecord) (*jobMetadata, error) {
	// default dry run to true
	meta := jobMetadata{isDryRun: true}

	if isDryRunIntf, exists := jobRecord.Metadata[IS_DRY_RUN_METADATA_KEY]; exists {
		var isBool bool
		if meta.isDryRun, isBool = isDryRunIntf.(bool); !isBool {
			return nil, errors.New(fmt.Sprintf("isDryRun must be a bool, got %v", isDryRunIntf))
		}
	}

	if programIdsIntf, exists := jobRecord.Metadata[PROGRAM_IDS_METADATA_KEY]; exists {
		var (
			programIdIntfs []interface{}
			isStringArr    bool
		)
		if programIdIntfs, isStringArr = programIdsIntf.([]interface{}); !isStringArr {
			return nil, errors.New(fmt.Sprintf(
				"programIds must be an array of strings, got %v", programIdsIntf))
		}
		meta.programIds = make([]string, len(programIdIntfs))
		for i, programIdIntf := range programIdIntfs {
			if programId, isString := programIdIntf.(string); !isString {
				return nil, errors.New(fmt.Sprintf(
					"programIds must be an array of strings, got %v", programIdsIntf))
			} else {
				meta.programIds[i] = programId
			}
		}
	} else {
		return nil, errors.New("jobRecord missing programIds")
	}

	youngestUpperGradYearPtr, err := jobmine_utility.UIntFromJobRecord(
		jobRecord, YOUNGEST_UPPER_GRAD_YEAR_METADATA_KEY)
	if err != nil {
		return nil, err
	} else if youngestUpperGradYearPtr == nil {
		return nil, errors.New("jobRecord missing youngestUpperGradYear")
	}
	meta.youngestUpperGradYear = *youngestUpperGradYearPtr

	maxLowerYearsPerUpperYearPtr, err := jobmine_utility.UIntFromJobRecord(
		jobRecord, MAX_LOWER_YEARS_PER_UPPER_YEAR_METADATA_KEY)
	if err != nil {
		return nil, err
	} else if maxLowerYearsPerUpperYearPtr == nil {
		return nil, errors.New("jobRecord missing maxLowerYearsPerUpperYear")
	}
	meta.maxLowerYearsPerUpperYear = *maxLowerYearsPerUpperYearPtr

	maxUpperYearsPerLowerYearPtr, err := jobmine_utility.UIntFromJobRecord(
		jobRecord, MAX_UPPER_YEARS_PER_LOWER_YEAR_METADATA_KEY)
	if err != nil {
		return nil, err
	} else if maxUpperYearsPerLowerYearPtr == nil {
		return nil, errors.New("jobRecord missing maxUpperYearsPerLowerYear")
	}
	meta.maxUpperYearsPerLowerYear = *maxUpperYearsPerLowerYearPtr

	meta.termStartTime, err = jobmine_utility.TimeFromJobRecord(
		jobRecord, TERM_START_TIME_METADATA_KEY)
	if err != nil {
		return nil, err
	}
	meta.termEndTime, err = jobmine_utility.TimeFromJobRecord(
		jobRecord, TERM_END_TIME_METADATA_KEY)
	if err != nil {
		return nil, err
	}
	return &meta, nil
}

func getTasksToCreate(db *gorm.DB, jobRecord jobmine.JobRecord) ([]jobmine.Metadata, error) {
	meta, err := parseJobMetadata(jobRecord)
	if err != nil {
		return nil, err
	}

	userIds, err := GetFilteredLowerAndAllUpperYears(db, meta.programIds, meta.youngestUpperGradYear,
		meta.termStartTime, meta.termEndTime)
	if err != nil {
		return nil, err
	}

	var strat recommendations.RecommendationStrategy

	if meta.termStartTime == nil {
		strat = getRecommendationStrategy(
			meta.maxLowerYearsPerUpperYear,
			meta.maxUpperYearsPerLowerYear,
			meta.youngestUpperGradYear,
		)
	} else {
		// Get users that verified the Winter 2019 whitelist and make them the blacklist for the
		// older downrank. (I know a little confusing, but the blacklist refers to users that we
		// omit from the downrank, so they will actually not be penalized during matching).
		blacklistUserIds, err := verify_link.GetVerifiedUserIds(
			db, verify_link.LINK_TYPE_WHITELIST_WINTER_2019)
		if err != nil {
			return nil, err
		}

		strat = getRecommendationStrategyWithOlderDownrank(
			meta.maxLowerYearsPerUpperYear,
			meta.maxUpperYearsPerLowerYear,
			meta.youngestUpperGradYear,
			*meta.termStartTime,
			blacklistUserIds,
		)
	}

	fetcherOptions := recommendations.UserFetcherOptions{UserIds: userIds}
	matches, err := recommendations.Recommend(db, fetcherOptions, strat)
	if err != nil {
		return nil, err
	}

	if meta.isDryRun {
		err = printMatches(db, matches)
		if err != nil {
			return nil, err
		}
		return []jobmine.Metadata{}, nil
	} else {
		metadatas := make([]jobmine.Metadata, len(matches))
		for i, match := range matches {
			metadatas[i] = packageTaskRecordMetadata(userMatch{
				menteeId: match.UserOneId,
				mentorId: match.UserTwoId,
			})
		}

		return metadatas, nil
	}
}

var SeedJobSpec jobmine.JobSpec = jobmine.JobSpec{
	JobType:          SEED_MENTORSHIPS_JOB,
	TaskSpec:         seedTaskSpec,
	GetTasksToCreate: getTasksToCreate,
}

// CreateSeedJob Creates a seed job record to get run at some point.
func CreateSeedJob(
	db *gorm.DB,
	runId string,
	isDryRun bool,
	programIds []string,
	youngestUpperGradYear uint,
	maxLowerYearsPerUpperYear uint,
	maxUpperYearsPerLowerYear uint,
	termStartTime *time.Time,
	termEndTime *time.Time,
) error {
	metadata := map[string]interface{}{
		PROGRAM_IDS_METADATA_KEY:                    programIds,
		YOUNGEST_UPPER_GRAD_YEAR_METADATA_KEY:       youngestUpperGradYear,
		MAX_LOWER_YEARS_PER_UPPER_YEAR_METADATA_KEY: maxLowerYearsPerUpperYear,
		MAX_UPPER_YEARS_PER_LOWER_YEAR_METADATA_KEY: maxUpperYearsPerLowerYear,
		IS_DRY_RUN_METADATA_KEY:                     isDryRun,
	}
	if termStartTime != nil {
		metadata[TERM_START_TIME_METADATA_KEY] = jobmine_utility.FormatTime(*termStartTime)
	}
	if termEndTime != nil {
		metadata[TERM_END_TIME_METADATA_KEY] = jobmine_utility.FormatTime(*termEndTime)
	}

	if err := db.Create(&jobmine.JobRecord{
		JobType:  SEED_MENTORSHIPS_JOB,
		RunId:    runId,
		Metadata: metadata,
		Status:   jobmine.STATUS_CREATED,
	}).Error; err != nil {
		return err
	}
	return nil
}
