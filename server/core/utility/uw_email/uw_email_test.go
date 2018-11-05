package uw_email

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	type testcase struct {
		email string
		valid bool
	}
	tests := []testcase{
		{ email: "foo@uwaterloo.ca", valid: true },
		{ email: "foo@edu.uwaterloo.ca", valid: true },
		{ email: "foo", valid: false },
		{ email: "foo@gmail.com", valid: false },
		{ email: "foo@uwaterloo.com", valid: false },
		{ email: "^>>hello@edu.uwaterloo.ca", valid: false },
		{ email: "foo@bar.edu.uwaterloo.ca", valid: false },
	}
	for _, test := range tests {
		assert.Equal(t, test.valid, Validate(test.email), "Expected '%s' validation to be %v", test.email, test.valid)
	}
}

func TestNormalize(t *testing.T) {
	type testcase struct {
		raw string
		normalized string
	}
	tests := []testcase{
		{ raw: "foo@uwaterloo.ca", normalized: "foo@edu.uwaterloo.ca" },
		{ raw: "foo@edu.uwaterloo.ca", normalized: "foo@edu.uwaterloo.ca" },
		{ raw: "Foo@edu.uwaterloo.ca", normalized: "foo@edu.uwaterloo.ca" },
	}
	for _, test := range tests {
		uwEmail := FromString(test.raw)
		assert.Equal(t, test.normalized, uwEmail.ToStringNormalized())
		assert.Equal(t, test.raw, uwEmail.ToStringRaw())
	}
}
