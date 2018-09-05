package linking

import "fmt"

// NOTE: must be kept in sync with index.tsx in root of letstalk
const (
	QR_SCANNER_URL        = "QrScanner"
	MATCH_PROFILE_URL     = "MatchProfile/%d"
	NOTIFICATION_VIEW_URL = "NotificationView"
	ADHOC_URL             = "NotificationContent/%d"
	QR_CODE_URL           = "QrCode"
)

func GetQrScannerUrl() string {
	return QR_SCANNER_URL
}

func GetMatchProfileUrl(userId uint) string {
	return fmt.Sprintf(MATCH_PROFILE_URL, userId)
}

func GetNotificationViewUrl() string {
	return NOTIFICATION_VIEW_URL
}

func GetAdhocLink(notificationId uint) string {
	return fmt.Sprintf(ADHOC_URL, notificationId)
}

func GetQrCodeUrl() string {
	return QR_CODE_URL
}
