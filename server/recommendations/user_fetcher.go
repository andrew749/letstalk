package recommendations

import (
	"time"

	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

// Options to be used with FetchUsers
// CreateAtStart, if provided, filters out users created before this time
// CreateAtEnd, if provided, filters out users created after this time
// UserIds, if provided, only fetches users with these ids
type UserFetcherOptions struct {
	CreatedAtStart *time.Time
	CreatedAtEnd   *time.Time
	UserIds        []data.TUserID
}

// Fetches users from the database.
func FetchUsers(
	db *gorm.DB,
	options UserFetcherOptions, // options related to which users to load
	preloadObjects []string, // objects to preload for each user (e.g. cohort, positions)
) ([]data.User, error) {
	var users []data.User

	query := db
	if options.CreatedAtStart != nil {
		query = query.Where("created_at >= ?", *options.CreatedAtStart)
	}
	if options.CreatedAtEnd != nil {
		query = query.Where("created_at <= ?", *options.CreatedAtEnd)
	}
	if options.UserIds != nil {
		query = query.Where("user_id IN (?)", options.UserIds)
	}
	if preloadObjects != nil {
		for _, object := range preloadObjects {
			query = query.Preload(object)
		}
	}
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
