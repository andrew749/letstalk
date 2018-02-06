package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenderIdByNameMale(t *testing.T) {
	assert.Equal(t, GenderIdByName("MALE"), 2)
}

func TestGenderIdByNameFemale(t *testing.T) {
	assert.Equal(t, GenderIdByName("FEMALE"), 1)
}

func TestGenderIdByNameOther(t *testing.T) {
	assert.Equal(t, GenderIdByName("test"), 0)
}
