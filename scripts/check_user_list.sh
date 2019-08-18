#!/bin/bash

set -e
filename=$1
sessionId=$2
while IFS= read -r line;
do
user=${line%$'\r'}
url="https://api.hiveapp.org/admin/user_exists?email=$user"
curl -s -X GET $url -H 'sessionId: $sessionId' | jq -r '.Result' | xargs -I{} echo "${user}={}";

done < $filename
