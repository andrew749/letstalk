package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPasswordInvertible(t *testing.T) {
	hashedPassword, err := HashPassword("1")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.True(t, CheckPasswordHash("1", hashedPassword))
}
