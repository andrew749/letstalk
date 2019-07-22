package match_round_commit_job

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/romana/rlog"

	"letstalk/server/core/api"
	"letstalk/server/core/connection"
	"letstalk/server/core/query"
	"letstalk/server/data"
	"letstalk/server/jobmine"
	"letstalk/server/jobmine_utility"
)

const MATCH_ROUND_COMMIT_JOB jobmine.JobType = "MatchRoundCommitJob"

// TaskRecord keys
const (
	MENTEE_USER_ID_METADATA_KEY = "menteeUserId"
	MENTOR_USER_ID_METADATA_KEY = "mentorUserId"
)

// JobRecord keys
const (
	MATCH_ROUND_ID_METADATA_KEY = "matchRoundId"
)

type userMatch struct {
	menteeId data.TUserID
	mentorId data.TUserID
}

func packageTaskRecordMetadata(match data.MatchRoundMatch) map[string]interface{} {
	return map[string]interface{}{
		MENTEE_USER_ID_METADATA_KEY: match.MenteeUserId,
		MENTOR_USER_ID_METADATA_KEY: match.MentorUserId,
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
	// TODO(match-api): Annotate connection with group that was used to create it.

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

func getMatchRoundIdFromJobRecord(jobRecord jobmine.JobRecord) (data.TMatchRoundID, error) {
	if matchRoundIdIntf, ok := jobRecord.Metadata[MATCH_ROUND_ID_METADATA_KEY]; ok {
		if matchRoundId, ok := matchRoundIdIntf.(data.TMatchRoundID); ok {
			return matchRoundId, nil
		} else {
			return 0, errors.New("matchRoundId is not an int")
		}
	} else {
		return 0, errors.New("matchRoundId not provided in job record")
	}
}

func getTasksToCreate(db *gorm.DB, jobRecord jobmine.JobRecord) ([]jobmine.Metadata, error) {
	matchRoundId, err := getMatchRoundIdFromJobRecord(jobRecord)
	if err != nil {
		return nil, err
	}

	var matches []data.MatchRoundMatch
	err = db.Where(&data.MatchRoundMatch{MatchRoundId: matchRoundId}).Find(matches).Error
	if err != nil {
		return nil, err
	}

	metadatas := make([]jobmine.Metadata, 0, len(matches))
	for _, match := range matches {
		metadatas = append(metadatas, packageTaskRecordMetadata(match))
	}

	return metadatas, nil
}

var CommitJobSpec jobmine.JobSpec = jobmine.JobSpec{
	JobType:          MATCH_ROUND_COMMIT_JOB,
	TaskSpec:         seedTaskSpec,
	GetTasksToCreate: getTasksToCreate,
}

func createRunId(matchRoundId data.TMatchRoundID) string {
	return fmt.Sprintf("match-round-commit-%d", matchRoundId)
}

// CreateMatchJob creates a job to commit match round matches, meaning that it creates connections
// for those matches and sends emails to users about their new matches.
//
// NOTE: It is expected that this runs in a transaction to avoid race conditions.
func CreateCommitJob(
	db *gorm.DB,
	matchRoundId data.TMatchRoundID,
) (*string, error) {
	runId := createRunId(matchRoundId)
	err := db.Where(&jobmine.JobRecord{RunId: runId}).Find(&data.JobRecord{}).Error
	if err == nil {
		return nil, errors.New(
			fmt.Sprintf("Job record for match round %d already exists", matchRoundId),
		)
	} else if !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	metadata := map[string]interface{}{
		MATCH_ROUND_ID_METADATA_KEY: matchRoundId,
	}
	if err := db.Create(&jobmine.JobRecord{
		JobType:  MATCH_ROUND_COMMIT_JOB,
		RunId:    runId,
		Metadata: metadata,
		Status:   jobmine.STATUS_CREATED,
	}).Error; err != nil {
		return nil, err
	}

	return &runId, nil
}
