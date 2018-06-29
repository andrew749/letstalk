package auth

import "letstalk/server/data"

func HasAdminAccess(authUser *data.User) bool {
	return authUser != nil && authUser.Role == data.USER_ROLE_ADMIN
}
