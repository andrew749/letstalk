package converters

import (
	"fmt"
	"letstalk/server/data"
	"letstalk/server/utility"
)

// GetManagedGroupReferralEmail Get the link to refer a user to a group
func GetManagedGroupReferralLink(groupUUID data.TGroupID) string {
	return fmt.Sprintf("%s/registerWithGroup/%s", utility.GetWebappUrl(), groupUUID)
}
