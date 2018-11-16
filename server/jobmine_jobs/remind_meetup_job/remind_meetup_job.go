package remind_meetup_job

import (
	"errors"
	"fmt"

	"letstalk/server/core/query"
	"letstalk/server/data"
	"letstalk/server/jobmine"

	"github.com/jinzhu/gorm"
)

const RemindMeetupJob jobmine.JobType = "RemindMeetupJob"

type UserType string

const (
	USER_TYPE_MENTOR UserType = "MENTOR"
	USER_TYPE_MENTEE UserType = "MENTEE"
)

const (
	UserTypeMetadataKey             = "userType"
	UserIdMetadataKey               = "userId"
	ConnectionUserIdMetadataKey     = "connectionUserId"
	ConnectionFirstNameMetadataKey  = "connectionFirstName"
	ConnectionLastNameMetadataKey   = "connectionLastName"
	ConnectionProfilePicMetadataKey = "connectionProfilePic"
)

func packageTaskRecordMetadata(
	userId data.TUserID,
	userType UserType,
	connectionUserId data.TUserID,
	connectionFirstName string,
	connectionLastName string,
	connectionProfilePic *string,
) map[string]interface{} {
	return map[string]interface{}{
		UserIdMetadataKey:               userId,
		UserTypeMetadataKey:             userType,
		ConnectionUserIdMetadataKey:     connectionUserId,
		ConnectionFirstNameMetadataKey:  connectionFirstName,
		ConnectionLastNameMetadataKey:   connectionLastName,
		ConnectionProfilePicMetadataKey: connectionProfilePic,
	}
}

var reminderTaskSpec = jobmine.TaskSpec{}

func getUserType(userId data.TUserID, mentorUserId data.TUserID) UserType {
	if userId == mentorUserId {
		return USER_TYPE_MENTOR
	} else {
		return USER_TYPE_MENTEE
	}
}

func getTasksToCreate(db *gorm.DB, jobRecord jobmine.JobRecord) (*[]jobmine.Metadata, error) {
	connections, err := query.GetAllMentorshipConnections(db)
	if err != nil {
		return nil, err
	}

	// TODO(wojtek): Filter out connections that have already received this notification

	metadata := make([]jobmine.Metadata, len(connections)*2)
	for _, connection := range connections {
		if connection.Mentorship == nil || connection.UserOne == nil || connection.UserTwo == nil {
			return nil, errors.New(fmt.Sprintf("Connection %d is missing data", connection.ConnectionId))
		}

		metadata1 := jobmine.Metadata(packageTaskRecordMetadata(
			connection.UserOneId,
			getUserType(connection.UserOneId, connection.Mentorship.MentorUserId),
			connection.UserTwoId,
			connection.UserTwo.FirstName,
			connection.UserTwo.LastName,
			connection.UserTwo.ProfilePic,
		))
		metadata2 := jobmine.Metadata(packageTaskRecordMetadata(
			connection.UserTwoId,
			getUserType(connection.UserTwoId, connection.Mentorship.MentorUserId),
			connection.UserOneId,
			connection.UserOne.FirstName,
			connection.UserOne.LastName,
			connection.UserOne.ProfilePic,
		))
		metadata = append(metadata, metadata1, metadata2)
	}
	return &metadata, nil
}

var ReminderJobSpec jobmine.JobSpec = jobmine.JobSpec{
	JobType:          RemindMeetupJob,
	TaskSpec:         reminderTaskSpec,
	GetTasksToCreate: getTasksToCreate,
}
