package ctx

import (
	"letstalk/server/core/errs"
	"letstalk/server/core/search"
	"letstalk/server/core/sessions"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic"
)

type Context struct {
	GinContext     *gin.Context
	Db             *gorm.DB
	Es             *elastic.Client
	SessionData    *sessions.SessionData
	SessionManager *sessions.ISessionManagerBase
	Result         interface{}
}

func NewContext(
	g *gin.Context,
	db *gorm.DB,
	es *elastic.Client,
	sessionData *sessions.SessionData,
	sm *sessions.ISessionManagerBase,
) *Context {
	return &Context{
		GinContext:     g,
		Db:             db,
		Es:             es,
		SessionData:    sessionData,
		SessionManager: sm,
	}
}

func WithinTx(db *gorm.DB, f func(*gorm.DB) error) error {
	tx := db.Begin()
	if err := f(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func WithinTxRequestErr(db *gorm.DB, f func(*gorm.DB) errs.Error) errs.Error {
	tx := db.Begin()
	if err := f(tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return errs.NewDbError(err)
	}
	return nil
}

// WithinTx provides a transaction object to the given function and automatically performs rollback
// if an error is returned.
func (c *Context) WithinTx(f func(*gorm.DB) error) error {
	return WithinTx(c.Db, f)
}

func (c *Context) SearchClientWithContext() *search.ClientWithContext {
	return search.NewClientWithContext(c.Es, c.GinContext.Request.Context())
}
