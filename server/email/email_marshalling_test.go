package email

import (
	"testing"

	"github.com/romana/rlog"
	"github.com/stretchr/testify/assert"
)

func TestMarshallEmail(t *testing.T) {
	testStruct := struct {
		A string `email_sub:":Test"`
	}{"test"}

	res := MarshallEmailSubstitutions(testStruct)
	rlog.Debug(res)
	assert.Contains(t, res, ":Test")
}
