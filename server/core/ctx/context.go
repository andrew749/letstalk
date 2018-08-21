package ctx

import (
	"net/http"

	"letstalk/server/core/search"
	"letstalk/server/core/sessions"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic"
)

type Context struct {
	GinContext     *gin.Context
	Db             *gorm.DB
	SearchClient   *search.RequestSearchClient
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
	var request *http.Request = nil
	if g != nil {
		request = g.Request
	}
	return &Context{
		GinContext:     g,
		Db:             db,
		SearchClient:   search.NewSearchClient(es, request),
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

// WithinTx provides a transaction object to the given function and automatically performs rollback
// if an error is returned.
func (c *Context) WithinTx(f func(*gorm.DB) error) error {
	return WithinTx(c.Db, f)
}
