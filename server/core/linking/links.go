package linking

import (
	"fmt"
	"letstalk/server/data"
	"letstalk/server/utility"
)

// NOTE: must be kept in sync with index.tsx in root of letstalk
const (
	QR_SCANNER_URL        = "QrScanner"
	MATCH_PROFILE_URL     = "MatchProfile?userId=%d"
	NOTIFICATION_VIEW_URL = "NotificationView"
	ADHOC_URL             = "NotificationContent?notificationId=%d"
	QR_CODE_URL           = "QrCode"
	ADD_POSITION_URL      = "AddPosition"
	ADD_GROUP_URL         = "AddGroup"
	ADD_SIMPLE_TRAIT_URL  = "AddSimpleTrait"
	UPDATE_PERSONAL_URL   = "UpdatePersonal"
)

func wrapWithUrlBase(url string) string {
	return fmt.Sprintf("%s://%s", utility.GetDeeplinkPrefix(), url)
}

func GetQrScannerUrl() string {
	return wrapWithUrlBase(QR_SCANNER_URL)
}

func GetAddPositionUrl() string {
	return wrapWithUrlBase(ADD_POSITION_URL)
}

func GetAddGroupUrl() string {
	return wrapWithUrlBase(ADD_GROUP_URL)
}

func GetAddSimpleTraitUrl() string {
	return wrapWithUrlBase(ADD_SIMPLE_TRAIT_URL)
}

func GetUpdatePersonalUrl() string {
	return wrapWithUrlBase(UPDATE_PERSONAL_URL)
}

func GetMatchProfileUrl(userId data.TUserID) string {
	return wrapWithUrlBase(fmt.Sprintf(MATCH_PROFILE_URL, userId))
}

func GetMatchProfileWithButtonUrl(userId data.TUserID) string {
	return fmt.Sprintf("%s&showRequestButton=true", GetMatchProfileUrl(userId))
}

func GetNotificationViewUrl() string {
	return wrapWithUrlBase(NOTIFICATION_VIEW_URL)
}

func GetAdhocLink(notificationId uint) string {
	return wrapWithUrlBase(fmt.Sprintf(ADHOC_URL, notificationId))
}

func GetQrCodeUrl() string {
	return wrapWithUrlBase(QR_CODE_URL)
}
