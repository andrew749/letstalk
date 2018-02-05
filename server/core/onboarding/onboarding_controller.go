package onboarding

import (
	"letstalk/server/core/ctx"
)
// TODO(acod): make this flow dynamic in nature
// i.e. fetch a onboarding type and the posible options

```
{
  "program": string,
  "sequence": string (i.e. 4STREAM, 8STREAM),
  "graduating_year": int,
}
```
// Update a user with new information for their school
func UpdateUserSchoolInfo(c *ctx.Context) {
  program := c.GinContext.Query("program")
  sequence := c.GinContext.Query("sequence")
  graduatingYear := c.GinContext.Query("graduating_year")
}

func RegisterUserInCohort(user api.User) errs.Error {
  // TODO(acod): add cohort data for the user.
}
