#!/bin/sh

# Set up adb reverse tunneling to run the dev app over usb.

adb -d reverse tcp:19000 tcp:19000
adb -d reverse tcp:19001 tcp:19001
adb -d reverse tcp:8080 tcp:80

echo "Expo server at exp://localhost:19000"
echo "Hive server at localhost:8080"

