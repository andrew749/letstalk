#!/bin/bash

set -e
set -x

JOB_ID=${1}
MENTOR_BODY="Hope you've all been having a great week. We've heard some very positive reactions to the meetups so far, keep up the great work! Now that coop is soon approaching, the 2023 SEs will soon experience their first work term and they're full of questions.<br><br>To help them out, we're hosting a <a href=\"https://www.facebook.com/events/1217688428378110/\">SE Coop Panel and Mentorship Mixer</a> from <b>6-8pm Thursday March 14</b> at <b>E7 4433</b>. Come support a panel of your peers to talk about topics ranging from how to crush your work term as an intern to effectively learning from your mentors. If you haven't met up with your mentees yet, invite them to come with you to the event! <br><br>Hope to see you there, <br> The Hive team 🐝"

MENTOR_QUERY="select user_id from se_mentors;";

MENTEE_BODY="Hope you've all been having a great week. We've heard some very positive reactions to the meetups so far, keep up the great work! Now that coop is soon approaching, it's worth starting a conversation on how to make the most of your next term. <br><br>From how to crush your work term as an intern to effectively learning from your mentors, let's hear from some of our upper year SE panelists at the <a href=\"https://www.facebook.com/events/1217688428378110/\">SE Coop Panel and Mentorship Mixer</a> from <b>6-8pm Thursday March 14</b> at <b>E7 4433</b>. This is a great way to chat with some upper years to see what insights they drew from their coops so far. If you haven't met up with your mentor yet, invite them to come with you to the event! <br><br>Hope to see you there, <br> The Hive team 🐝"

MENTEE_QUERY="select user_id from se_mentees;";

PUSH_TITLE="SE Coop Panel and Mentorship Mixer";
PUSH_MESSAGE="Come out this Thursday for some insightful talks!";

MESSAGE_CAPTION="Come out this Thursday for some insightful talks!";

echo "Sending mentor notifications";
# send mentor notifications
jobmine_jobs/notification_script_templates/send_generic_notification.sh "${JOB_ID}_mentor" "false" "$PUSH_TITLE" "$PUSH_MESSAGE" "$MESSAGE_CAPTION" "$MENTOR_BODY" "$MENTOR_QUERY" "true" "true" "true" "d-e54d5bdf1c7b4155a42adc379a82369b";

echo "Sending mentee notifications";
# send mentee notifications
jobmine_jobs/notification_script_templates/send_generic_notification.sh "${JOB_ID}_mentee" "false" "$PUSH_TITLE" "$PUSH_MESSAGE" "$MESSAGE_CAPTION" "$MENTEE_BODY" "$MENTEE_QUERY" "true" "true" "true" "d-e54d5bdf1c7b4155a42adc379a82369b";
