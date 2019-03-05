#!/bin/bash

set -e
set -x

JOB_ID=${1}
MENTOR_BODY="It's been a little over a week since you've received your mentee matches. If you haven't met up yet, it's still not too late! Reach out to your mentee matches as soon as you can for some sweet sweet boba.<br><br>To claim the bubble tea credits, you can collect gift cards from either <a href="https://www.facebook.com/andrew.codispoti">Andrew Codispoti</a> or <a href="https://www.facebook.com/stevendjkong">Steven Kong</a> (and giving us a printed copy of <a href="https://uwaterloo.ca/finance/sites/ca.finance/files/uploads/files/gift_reporting_form_150225.pdf">this short form</a>). There is 5$ in credits from Sweet Dreams allocated per student enrolled in the mentorship program. Thank you to those who already took the initiative to meet up, we've been hearing great things about the experience so far! Keep it up 🙂!<br><br>Don't feel daunted about this first meeting! It's a casual way to get to know another SE students. Anticipating some questions  beforehand can definitely help drive the conversation! Let them know what you wish you knew when you were in their shoes.<br><br>The Hive team 🐝"

MENTOR_QUERY="select user_id, first_name from se_mentors;";

MENTEE_BODY="It's been a little over a week since you've received your mentor matches. If you haven't met up with your mentors yet, it's still not too late! Try reaching out to your mentor if they haven't contacted you yet. <br><br>Don't feel daunted about this first meeting! It's a casual way to get to know another SE student. If you're not sure what to ask, give <a href="https://medium.com/@uwhive/being-an-effective-mentee-6ef9e9177498">our mentee guide</a> a read on how you can make the most of your mentorship experience. Having some questions prepared beforehand can definitely help drive the conversation and make effective use of your mentor's time! <br><br>The Hive team 🐝"

MENTEE_QUERY="select user_id, first_name from se_mentees;";

PUSH_TITLE="{{.first_name}}, have you met up yet?" ;
PUSH_MESSAGE="Meet up with your matches!";

MESSAGE_CAPTION="Get your free bubble tea!";

EMAIL_TEMPLATE_MENTEE="d-6e861dff5bf64f7c851eea2f6dd4dc50";
EMAIL_TEMPLATE_MENTOR="d-7fc85f51e72f47908714260ebb37d15d";


echo "Sending mentor notifications";
# send mentor notifications
jobmine_jobs/notification_script_templates/send_generic_notification.sh "${JOB_ID}_mentor" "false" "$PUSH_TITLE" "$PUSH_MESSAGE" "$MESSAGE_CAPTION" "$MENTOR_BODY" "$MENTOR_QUERY" "true" "true" "true" "$EMAIL_TEMPLATE_MENTOR";

echo "Sending mentee notifications";
# send mentee notifications
jobmine_jobs/notification_script_templates/send_generic_notification.sh "${JOB_ID}_mentee" "false" "$PUSH_TITLE" "$PUSH_MESSAGE" "$MESSAGE_CAPTION" "$MENTEE_BODY" "$MENTEE_QUERY" "true" "true" "true" "$EMAIL_TEMPLATE_MENTEE";