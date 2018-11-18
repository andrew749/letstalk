# Manual Task Scheduler

Helps with the creation of adhoc jobmine tasks so that we can test stuff oneoff.

## Example commands:

### Create a new reminder job for a specific user
`./run_in_env.sh './task_scheduler -jobMetadata "{}"  -taskMetadata="{\"userId\":1,\"reminderType\":\"REMINDER_TYPE_TRAIT\"}" -jobType "RemindOnboardJob" -runId "adhoc_task_4"'`
