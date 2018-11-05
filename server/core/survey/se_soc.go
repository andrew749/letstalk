package survey

import "letstalk/server/core/api"

var Se_soc_v1 api.Survey = api.Survey{
	Group:   "SE_SOC",
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
			Key:    "interests",
			Prompt: "If you had to pick one, which area interests you the most currently?",
			Options: []api.SurveyOption{
				{Key: "distributed", Text: "Distributed Systems"},
				{Key: "web_dev", Text: "Web Development"},
				{Key: "mobile_dev", Text: "Mobile Development"},
				{Key: "security", Text: "Security & Cryptography"},
				{Key: "machine_learning", Text: "Machine Learning/Artificial Intelligence"},
				{Key: "research", Text: "Research"},
				{Key: "unknown", Text: "I'm still figuring it out"},
			},
		},
		{
			Key:    "role_explore",
			Prompt: "What role do you most want to explore during your time at university?",
			Options: []api.SurveyOption{
				{Key: "software_engineering", Text: "Software Engineering"},
				{Key: "design", Text: "Design"},
				{Key: "product_management", Text: "Product Management"},
				{Key: "data_science", Text: "Data Science"},
			},
		},
		{
			Key:    "extracurriculars",
			Prompt: "What extracurricular activity/hobby do you enjoy the most?",
			Options: []api.SurveyOption{
				{Key: "arts", Text: "Drawing or painting"},
				{Key: "media", Text: "Movies or books"},
				{Key: "visual_art", Text: "Photography or videography"},
				{Key: "video_games", Text: "Video games"},
				{Key: "playing_music", Text: "Playing music"},
				{Key: "sports", Text: "Sports and exercise"},
				{Key: "partying", Text: "Getting lit and partying"},
				{Key: "cooking", Text: "Cooking"},
				{Key: "coding", Text: "Hackathons and side projects"},
				{Key: "other", Text: "Other"},
			},
		},
	},
}
