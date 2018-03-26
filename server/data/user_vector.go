package data

type UserVector struct {
	User           User `gorm:"foreignkey:UserId"`
	UserId         int  `json:"user_id"`
	PreferenceType int  `json:"preference_type"`
	Sociable       int  `json:"sociable"`
	HardWorking    int  `json:"hard_working"`
	Ambitious      int  `json:"ambitious"`
	Energetic      int  `json:"energetic"`
	Carefree       int  `json:"carefree"`
	Confident      int  `json:"confident"`
}
