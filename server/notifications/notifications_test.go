package notifications

import (
	"encoding/json"
	"testing"

	"github.com/romana/rlog"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshallNotification(t *testing.T) {
	d := []byte(`{"data":[{"id":"3d55301b-1ab0-4657-a6c2-3d31398775b9","status":"ok"}]}`)
	var res ExpoNotificationSendResponse
	err := json.Unmarshal(d, &res)
	assert.NoError(t, err)
	rlog.Info(res)
	assert.Equal(t, "3d55301b-1ab0-4657-a6c2-3d31398775b9", res.Data[0].Id)
}
