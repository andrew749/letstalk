package data

type UserVector struct {
	UserId         int  `json:"userId"`
	User           User `gorm:"foreignkey:UserId"`
	PreferenceType int  `json:"preferenceType"`
	Sociable       int  `json:"sociable"`
	HardWorking    int  `json:"hardworking"`
	Ambitious      int  `json:"ambitious"`
	Energetic      int  `json:"energetic"`
	Carefree       int  `json:"carefree"`
	Confident      int  `json:"confident"`
}
