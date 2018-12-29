package survey

import "letstalk/server/core/api"

var Generic_v1 api.Survey = api.Survey{
	Group:   "GENERIC",
	Version: 1,
	Questions: []api.SurveyQuestion{
		{
			Key:    "free_time",
			Prompt: "You have some free time, what would you most likely be doing?",
			Options: []api.SurveyOption{
				{Key: "reading", Text: "Some reading, listening to music or watching a movie"},
				{Key: "artistic", Text: "Playing music, drawing or painting"},
				{Key: "lowkey", Text: "Playing games, grabbing dinner with friends or playing some sport"},
				{Key: "party", Text: "Going clubbing or to a music festival"},
			},
		},
		{
			Key:    "group_size",
			Prompt: "Do you prefer spending time with lots of people or chilling with a smaller group?",
			Options: []api.SurveyOption{
				{Key: "larger", Text: "More is better!"},
				{Key: "both", Text: "I can live with both"},
				{Key: "smaller", Text: "I'm more comfortable with smaller groups"},
			},
		},
		{
			Key:    "exercise",
			Prompt: "How often do you play sports/exercise?",
			Options: []api.SurveyOption{
				{Key: "daily", Text: "Every. Day."},
				{Key: "weekly", Text: "Weekly, I still have a routine"},
				{Key: "sometimes", Text: "Somewhat regularly, when I find the time"},
				{Key: "rarely", Text: "When the wifi goes out, I prefer not to"},
			},
		},
		{
			Key:    "school_work",
			Prompt: "How important is school work to you?",
			Options: []api.SurveyOption{
				{Key: "minimally", Text: "We have school work? I prefer not to do any."},
				{Key: "moderately", Text: "I meet requirements. No more, no less."},
				{Key: "maximally", Text: "I study and work hard on school work for my success."},
			},
		},
		{
			Key:    "working_on",
			Prompt: "What are you currently working the most on during your time at university?",
			Options: []api.SurveyOption{
				{Key: "school", Text: "Courses and grades"},
				{Key: "responsibilities", Text: "Getting my life together, managing responsibilities"},
				{Key: "career", Text: "Planning out my career path"},
				{Key: "relationships", Text: "Building relationships"},
			},
		},
	},
}
