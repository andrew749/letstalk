package survey

import "letstalk/server/core/api"

var Generic_v1 api.Survey = api.Survey{
	Version: "generic_v1",
	Questions: []api.SurveyQuestion{
		{
			Key:    "friday_night",
			Prompt: "What activity would you most likely be doing on a Friday night?",
			Options: []api.SurveyOption{
				{ Key: "reading", Text: "Reading" },
				{ Key: "studying", Text: "Studying" },
				{ Key: "board_games", Text: "Board games" },
				{ Key: "partying", Text: "Partying" },
			},
		},
		{
			Key:    "social_gatherings",
			Prompt: "Do you prefer large social gatherings or smaller more intimate settings?",
			Options: []api.SurveyOption{
				{ Key: "large", Text: "Large gatherings" },
				{ Key: "small", Text: "Smaller settings" },
			},
		},
		{
			Key:    "athletics",
			Prompt: "How often do you play sports/exercise?",
			Options: []api.SurveyOption{
				{ Key: "daily", Text: "Daily" },
				{ Key: "weekly", Text: "Weekly" },
				{ Key: "sometimes", Text: "Somewhat regularly" },
				{ Key: "rarely", Text: "When the wifi goes out" },
			},
		},
		{
			Key:    "studying",
			Prompt: "How much of your time do you spend on schoolwork?",
			Options: []api.SurveyOption{
				{ Key:"minimum", Text: "As little as possible"},
				{ Key:"half", Text: "About half" },
				{ Key:"majority", Text: "The majority" },
				{ Key:"all", Text: "Almost all" },
			},
		},
		{
			Key:    "stress",
			Prompt: "What have you been most stressed about in University?",
			Options: []api.SurveyOption{
				{ Key:"courses", Text: "Courses and grades" },
				{ Key:"adulting", Text: "Adult responsibilities (bills, laundry, etc.)" },
				{ Key:"coop", Text: "Co-op and career" },
				{ Key:"socializing", Text: "Socializing and relationships" },
			},
		},
	},
}

