package survey

import (
	"testing"

	"letstalk/server/data"

	"github.com/stretchr/testify/assert"
)

func TestGenericV1SimilarityScoreWeightTotal(t *testing.T) {
	totalWeight := 0.0
	for _, question := range Generic_v1.Questions {
		totalWeight += genericV1Similarities[question.Key].weight
	}
	assert.Equal(t, 1.0, totalWeight)
}

func TestGenericV1SimilarityScoreNoResponses(t *testing.T) {
	noResponses := map[data.SurveyQuestionKey]data.SurveyOptionKey{}
	assert.Equal(t, 0.5, genericV1SimilarityScore(noResponses, noResponses))
}

func TestGenericV1SimilarityScore(t *testing.T) {
	r1 := map[data.SurveyQuestionKey]data.SurveyOptionKey{
		"free_time":   "reading",
		"group_size":  "smaller",
		"exercise":    "rarely",
		"school_work": "moderately",
		"working_on":  "responsibilities",
	}
	r2 := map[data.SurveyQuestionKey]data.SurveyOptionKey{
		"free_time":   "lowkey",
		"group_size":  "both",
		"exercise":    "sometimes",
		"school_work": "minimally",
		"working_on":  "career",
	}
	assert.Equal(t, 0.466, genericV1SimilarityScore(r1, r2))
}
