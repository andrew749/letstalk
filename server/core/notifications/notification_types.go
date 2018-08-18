package notifications

import "letstalk/server/data"

const (
	NOTIF_TYPE_NEW_CREDENTIAL_MATCH data.NotifType = "NEW_CREDENTIAL_MATCH"
	NOTIF_TYPE_ADHOC                data.NotifType = "ADHOC_NOTIFICATION"
	NOTIF_TYPE_REQUEST_TO_MATCH     data.NotifType = "REQUEST_TO_MATCH"
	NOTIF_TYPE_NEW_MATCH            data.NotifType = "NEW_MATCH"
	NOTIF_TYPE_MATCH_VERIFIED       data.NotifType = "MATCH_VERIFIED"
)
