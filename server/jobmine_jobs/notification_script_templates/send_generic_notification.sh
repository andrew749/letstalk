#!/bin/sh

# check number of arguments
if [ "$#" -ne "9" ]; then
echo "
Usage:
    $0 RUN_ID DRY_RUN TITLE MESSAGE CAPTION BODY USER_SELECTOR SEND_EMAIL SEND_PUSH
    - RUN_ID runId for this job
    - DRY_RUN Whether to actually create tasks to send.
    - TITLE The title that a push notification will get.
    - MESSAGE The message a push notification will get.
    - CAPTION The caption that will be printed in both the push and email.
    - BODY The content of the actual message
    - USER_SELECTOR An sql query defining which users to send to.
    - SEND_EMAIL 1 if we want to send an email, otherwise anything else
    - SEND_PUSH 1 if we want to send a push, otherwise anything else
"
exit 1;
fi

# default dry run
RUN_ID=${1}
DRY_RUN=${2}
TITLE=${3}
MESSAGE=${4}
CAPTION=${5}
BODY=${6}
USER_SELECTOR=${7}
SEND_EMAIL=${8}
SEND_PUSH=${9}

# determine which platforms to send to
PLATFORM_STRING=""
if [[ "$SEND_PUSH" -eq "1" ]]; then
 PLATFORM_STRING+=",\\\"notificationTemplate\\\":\\\"generic_notification.html\\\"";
fi

if [[ "$SEND_EMAIL" -eq "1" ]]; then
 PLATFORM_STRING+=", \\\"emailTemplate\\\": \\\"d-e54d5bdf1c7b4155a42adc379a82369b\\\""
fi
echo "$PLATFORM_STRING"

# not dry run, send a notification to all users. With appropriate templating.
DB_NET=letstalk_db_net ./run_in_env.sh \
"RLOG_LOG_LEVEL=DEBUG build/manual_job_scheduler \
 -runId $RUN_ID \
 -jobType GenericNotificationJob \
 -metadata \"{\\\"dryRun\\\":$DRY_RUN,\\\"message\\\":\\\"$MESSAGE\\\",\\\"title\\\":\\\"$TITLE\\\" $PLATFORM_STRING, \\\"userSelectorQuery\\\":\\\"$USER_SELECTOR\\\", \\\"data\\\": {\\\"caption\\\": \\\"$CAPTION\\\",\\\"body\\\": \\\"$BODY\\\"}}\"";
((JOB_ID+=1));