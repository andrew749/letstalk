package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNormalizedProgramMapping(t *testing.T) {
	mapping := GetNormalizedProgramMapping()
	assert.Contains(t, mapping, "SOFTWARE_ENGINEERING")
	assert.Contains(t, mapping, "COMPUTER_ENGINEERING")
}

func TestGetReverseNormalizedProgramMapping(t *testing.T) {
	mapping := GetReverseNormalizedProgramMapping()
	for _, program := range stream4Programs {
		assert.Contains(t, mapping, program)
	}
	for _, program := range stream8Programs {
		assert.Contains(t, mapping, program)
	}
	assert.Contains(t, mapping, "Software Engineering")
	assert.Contains(t, mapping, "Computer Engineering")
}
