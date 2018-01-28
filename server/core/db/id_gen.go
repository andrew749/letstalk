package db

import (
	"letstalk/server/core/ctx"
	"letstalk/server/data"
	"sync"

	"github.com/mijia/modelq/gmq"
)

var idMutex = sync.Mutex{}

// NumId safely generates a unique numerical id.
func NumId(c *ctx.Context) int {
	idMutex.Lock()
	defer idMutex.Unlock()
	var nextId int
	err := gmq.WithinTx(c.Db, func(tx *gmq.Tx) error {
		idGen, err := data.IdGenObjs.Select().OrderBy("-Id").One(tx)
		if err != nil {
			return err
		}
		idGen.NumId++
		nextId = idGen.NumId
		_, err = idGen.Insert(tx)
		return err
	})
	if err != nil {
		c.GinContext.Error(err)
		return 0
	}
	return nextId
}
