package connection

import (
	"fmt"
	"time"

	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/notifications"
	"letstalk/server/core/query"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"letstalk/server/core/meetup_reminder"
)

/**
 * PostRequestConnection creates a new unaccepted connection between two users.
 */
func PostRequestConnection(c *ctx.Context) errs.Error {
	var input api.ConnectionRequest
	if err := c.GinContext.BindJSON(&input); err != nil {
		return errs.NewRequestError("Failed to parse input")
	}

	if newConnection, err := HandleRequestConnection(c, input); err != nil {
		return err
	} else {
		c.Result = newConnection
	}
	return nil
}

// TODO(wojtechnology): Give this a more explicit public interface.
func HandleRequestConnection(
	c *ctx.Context,
	request api.ConnectionRequest,
) (*api.ConnectionRequest, errs.Error) {
	// Assert users exist and are not equal.
	authUser, _ := query.GetUserById(c.Db, c.SessionData.UserId)
	if c.SessionData.UserId == request.UserId {
		return nil, errs.NewRequestError("Cannot connect with self")
	}
	connUser, err := query.GetUserById(c.Db, request.UserId)
	if err != nil {
		return nil, errs.NewRequestError("Invalid user id")
	}
	// Assert request does not already exist.
	existing, dbErr := query.GetConnectionDetailsUndirected(c.Db, authUser.UserId, connUser.UserId)
	if dbErr != nil {
		return nil, errs.NewDbError(dbErr)
	}
	if existing != nil {
		return nil, errs.NewRequestError("Connection already exists")
	}
	// Save new connection and intent.
	connection := data.Connection{
		UserOneId: authUser.UserId,
		UserTwoId: connUser.UserId,
		CreatedAt: time.Now(),
	}
	intent := data.ConnectionIntent{
		Type:          request.IntentType,
		SearchedTrait: request.SearchedTrait,
		Message:       request.Message,
	}
	dbErr = c.WithinTx(func(tx *gorm.DB) error {
		if err := tx.Create(&connection).Error; err != nil {
			return err
		}
		intent.ConnectionId = connection.ConnectionId
		if err := tx.Create(&intent).Error; err != nil {
			return err
		}
		if err := notifications.ConnectionRequestedNotification(
			tx,
			connUser.UserId,
			authUser.UserId,
			fmt.Sprintf("%s %s", authUser.FirstName, authUser.LastName),
		); err != nil {
			return err
		}
		return nil
	})
	if dbErr != nil {
		return nil, errs.NewDbError(dbErr)
	}
	request.CreatedAt = connection.CreatedAt
	return &request, nil
}

/**
 * PostAcceptConnection accepts an existing connection request
 */
func PostAcceptConnection(c *ctx.Context) errs.Error {
	var input api.AcceptConnectionRequest
	if err := c.GinContext.BindJSON(&input); err != nil {
		return errs.NewRequestError("Failed to parse input")
	}

	if newConnection, err := HandleAcceptConnection(c, input); err != nil {
		return err
	} else {
		c.Result = newConnection
	}
	return nil
}

// TODO(wojtechnology): Give this a more explicit public interface.
func HandleAcceptConnection(
	c *ctx.Context,
	request api.AcceptConnectionRequest,
) (*api.ConnectionRequest, errs.Error) {
	// Assert request exists from request user to auth user.
	connection, err := query.GetConnectionDetails(c.Db, request.UserId, c.SessionData.UserId)
	if err != nil {
		return nil, errs.NewDbError(err)
	}
	if connection == nil || connection.DeletedAt != nil {
		return nil, errs.NewRequestError("No such connection request exists")
	}
	result := api.ConnectionRequest{
		UserId:        request.UserId,
		IntentType:    connection.Intent.Type,
		CreatedAt:     connection.CreatedAt,
		SearchedTrait: connection.Intent.SearchedTrait,
	}
	if connection.AcceptedAt != nil {
		// Already accepted, do nothing.
		result.AcceptedAt = connection.AcceptedAt
		return &result, nil
	}
	authUser, err := query.GetUserById(c.Db, c.SessionData.UserId)
	if err != nil {
		return nil, errs.NewRequestError("Cannot find myself")
	}
	connUser, err := query.GetUserById(c.Db, request.UserId)
	if err != nil {
		return nil, errs.NewRequestError("Invalid user id")
	}
	now := time.Now()
	connection.AcceptedAt = &now
	dbErr := c.WithinTx(func(tx *gorm.DB) error {
		if err := tx.Save(connection).Error; err != nil {
			return err
		}
		if err := notifications.ConnectionAcceptedNotification(
			tx,
			connUser.UserId,
			authUser.UserId,
			fmt.Sprintf("%s %s", authUser.FirstName, authUser.LastName),
		); err != nil {
			return err
		}
		// Schedule first meetup reminder for both users in the match.
		return meetup_reminder.ScheduleInitialReminder(tx, authUser.UserId, connUser.UserId)
	})
	if dbErr != nil {
		fmt.Printf("DEBUG, WTF got dbError: %v\n", dbErr)
		return nil, errs.NewDbError(err)
	}
	result.AcceptedAt = connection.AcceptedAt
	return &result, nil
}

/**
 * RemoveConnection removes an existing connection
 */
func RemoveConnection(c *ctx.Context) errs.Error {
	var input api.RemoveConnection
	if err := c.GinContext.BindJSON(&input); err != nil {
		return errs.NewRequestError("Failed to parse input")
	}

	return removeConnection(c, input)
}

func removeConnection(c *ctx.Context, request api.RemoveConnection) errs.Error {
	meUserId := c.SessionData.UserId
	youUserId := request.UserId

	connection, err := query.GetConnectionDetailsUndirected(c.Db, meUserId, youUserId)
	if err != nil {
		return errs.NewDbError(err)
	}
	if connection == nil {
		return nil
	}

	if err := c.Db.Delete(&connection).Error; err != nil {
		return errs.NewDbError(err)
	}

	return nil
}

// DataToApi converts a data.Connection to an api.ConnectionRequest.
// otherUserId: Id of the non-auth user involved in the connection.
// data: Must have non-nil Intent.
func DataToApi(otherUserId data.TUserID, data data.Connection) api.ConnectionRequest {
	return api.ConnectionRequest{
		UserId:        otherUserId,
		SearchedTrait: data.Intent.SearchedTrait,
		IntentType:    data.Intent.Type,
		Message:       data.Intent.Message,
		CreatedAt:     data.CreatedAt,
		AcceptedAt:    data.AcceptedAt,
	}
}
