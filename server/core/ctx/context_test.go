package ctx_test

import (
	"letstalk/server/core/ctx"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mijia/modelq/gmq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/http"
)

func TestNewContext(t *testing.T) {
	writer := http.TestResponseWriter{}
	g, _ := gin.CreateTestContext(&writer)
	db := &gmq.Db{}
	c := ctx.NewContext(g, db)
	assert.Equal(t, db, c.Db)
	assert.Equal(t, g, c.GinContext)
	assert.Nil(t, c.Result)
}
