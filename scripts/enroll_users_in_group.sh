#!/bin/bash

set -e
sessionId="$1"
groupUUID="$2"
filename="$3"
while IFS= read -r line;
do
email=${line%$'\r'}
curl -X POST \
  https://api.hiveapp.org/admin/enroll_user_in_group_by_email \
  -H 'Accept: */*' \
  -H 'Accept-Encoding: gzip, deflate' \
  -H 'Cache-Control: no-cache' \
  -H 'Connection: keep-alive' \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -H "sessionId: $sessionId" \
  -d "{
	\"groupUUID\": \"$groupUUID\",
	\"email\": \"$email\"
}"

done < $filename
