package connection
import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/data"

	"github.com/romana/rlog"
	"github.com/jinzhu/gorm"
	"time"
)

/**
 * PostRequestConnection creates a new unaccepted connection between two users.
 */
func PostRequestConnection(c *ctx.Context) errs.Error {
	var input api.Connection
	if err := c.GinContext.BindJSON(&input); err != nil {
		return errs.NewRequestError("Failed to parse input")
	}
	rlog.Info("Received input: ", input)

	if newConnection, err := handleRequestConnection(c, input); err != nil {
		return err
	} else {
		c.Result = newConnection
	}
	return nil
}

func handleRequestConnection(c *ctx.Context, request api.Connection) (*api.Connection, errs.Error) {
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
	existing, err := query.GetConnectionDetailsUndirected(c.Db, authUser.UserId, connUser.UserId)
	if err != nil {
		return nil, errs.NewDbError(err)
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
		Type: request.IntentType,
	}
	if len(request.SearchedTrait) > 0 {
		intent.SearchedTrait = &request.SearchedTrait
	}
	// TODO(aklen): send notification to requested user
	dbErr := c.WithinTx(func(tx *gorm.DB) error {
		if err := tx.Create(&connection).Error; err != nil {
			return err
		}
		intent.ConnectionId = connection.ConnectionId
		if err := tx.Create(&intent).Error; err != nil {
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
	var input api.Connection
	if err := c.GinContext.BindJSON(&input); err != nil {
		return errs.NewRequestError("Failed to parse input")
	}
	rlog.Info("Received input: ", input)

	if newConnection, err := handleAcceptConnection(c, input); err != nil {
		return err
	} else {
		c.Result = newConnection
	}
	return nil
}

func handleAcceptConnection(c *ctx.Context, request api.Connection) (*api.Connection, errs.Error) {
	// Assert request exists from request user to auth user.
	connection, err := query.GetConnectionDetails(c.Db, request.UserId, c.SessionData.UserId)
	if err != nil {
		return nil, errs.NewDbError(err)
	}
	if connection == nil || connection.DeletedAt != nil {
		return nil, errs.NewRequestError("No such connection request exists")
	}
	result := api.Connection{
		UserId: request.UserId,
		IntentType: connection.Intent.Type,
		CreatedAt: connection.CreatedAt,
	}
	if connection.Intent.SearchedTrait != nil {
		result.SearchedTrait = *connection.Intent.SearchedTrait
	}
	if connection.AcceptedAt != nil {
		// Already accepted, do nothing.
		result.AcceptedAt = connection.AcceptedAt
		return &result, nil
	}
	now := time.Now()
	connection.AcceptedAt = &now
	// TODO(aklen): send notification to accepted user
	if err := c.Db.Save(connection).Error; err != nil {
		return nil, errs.NewDbError(err)
	}
	result.AcceptedAt = connection.AcceptedAt
	return &result, nil
}
