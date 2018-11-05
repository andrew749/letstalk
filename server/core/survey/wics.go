package survey

import "letstalk/server/core/api"

var Wics_v1 api.Survey = api.Survey{
	Group: "WICS",
	Version: 1,
	Questions: []api.SurveyQuestion{
		{
			Key: "interests",
			Prompt: "If you had to pick one, which area interests you the most currently?",
			Options: []api.SurveyOption{
				{ Key: "distributed", Text: "Distributed Systems" },
				{ Key: "web_dev", Text: "Web Development" },
				{ Key: "mobile_dev", Text: "Mobile Development" },
				{ Key: "security", Text: "Security & Cryptography" },
				{ Key: "machine_learning", Text: "Machine Learning/Artificial Intelligence" },
				{ Key: "research", Text: "Research" },
				{ Key: "unknown", Text: "I'm still figuring it out" },
			},
		},
		{
			Key: "goals",
			Prompt: "What would you like to get out of this mentorship program?",
			Options: []api.SurveyOption{
				{ Key: "career", Text: "Career guidance" },
				{ Key: "academics", Text: "Academic guidance" },
				{ Key: "networking", Text: "Expand your network, meet some interesting people" },
				{ Key: "leadership", Text: "Leadership skills" },
			},
		},
		{
			Key: "meetup_frequency",
			Prompt: "How often do you want to meet up with your match?",
			Options: []api.SurveyOption{
				{ Key: "often", Text: "Multiple times per week" },
				{ Key: "weekly", Text: "Once a week" },
				{ Key: "biweekly", Text: "Biweekly" },
				{ Key: "monthly", Text: "Monthly" },
				{ Key: "request", Text: "On Request" },
			},
		},
	},
}

