#!/bin/bash

# go over a secrets file and echo them
for s in $(cat $1 | jq -r "to_entries|map(\"\(.key|ascii_upcase)=\(.value|tostring)\")|.[]" ); do
    export $s
done
