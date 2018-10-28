package api

type NotificationCampaignSendRequest struct {
	Title             string                 `json:"title" binding:"required"`
	Message           string                 `json:"message" binding:"required"`
	GroupId           string                 `json:"groupId" binding:"required"`
	RunId             string                 `json:"runId" binding:"required"`
	Thumbnail         *string                `json:"thumbnail"`
	TemplatePath      string                 `json:"templatePath" binding:"required"`
	TemplatedMetadata map[string]interface{} `json:"templateMetadata" binding:"required"`
}
