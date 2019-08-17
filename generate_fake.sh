#!/bin/sh

first=(
  Alton
  Brett
  Jeffery
  Eleanor
  Miriam
  Mamie
  Edna
  Roberto
  Hector
  Jose
  Ramon
  Lola
  Douglas
  Darlene
  Whitney
  Ida
  Aaron
  Van
  Luis
  Ashley
  Dale
  Bryan
  Mercedes
  Courtney
  Lionel
)
last=(
  Nelson
  Conner
  Joseph
  Owen
  Patton
  Simmons
  Webb
  Leonard
  Stone
  Watson
  Castro
  Gomez
  Rodgers
  Harper
  Hansen
  Hampton
  Ramirez
  Ball
  Cobb
)
emails=(
  EdnaNelson@uwaterloo.ca
  RobertoConner@uwaterloo.ca
  HectorJoseph@uwaterloo.ca
  JoseOwen@uwaterloo.ca
  RamonPatton@uwaterloo.ca
  LolaSimmons@uwaterloo.ca
  DouglasWebb@uwaterloo.ca
  DarleneLeonard@uwaterloo.ca
  WhitneyStone@uwaterloo.ca
  IdaWatson@uwaterloo.ca
  AaronCastro@uwaterloo.ca
  VanGomez@uwaterloo.ca
  LuisRodgers@uwaterloo.ca
  AshleyHarper@uwaterloo.ca
  DaleHansen@uwaterloo.ca
  BryanHampton@uwaterloo.ca
  MercedesRamirez@uwaterloo.ca
  CourtneyBall@uwaterloo.ca
  LionelCobb@uwaterloo.ca
)

createUser() {
    local firstname=$1
    local lastname=$2
    local email=$3
    local cohort=$4
    echo $firstname $lastname $email
#signup
curl -X POST \
  http://localhost/v1/signup \
  -H 'Accept: */*' \
  -H 'Accept-Encoding: gzip, deflate' \
  -H 'Cache-Control: no-cache' \
  -H 'Connection: keep-alive' \
  -H 'Content-Type: application/json' \
  -H 'Host: localhost' \
  -H 'Postman-Token: d51c641c-1c75-4cb3-b5ec-4d77bcc54121,af88169d-cbf2-447f-a107-b54d52759b4e' \
  -H 'User-Agent: PostmanRuntime/7.15.2' \
  -H 'cache-control: no-cache' \
  -d "{
	\"firstName\": \"$firstname\",
	\"lastName\": \"$lastname\",
	\"email\": \"$email\",
	\"gender\": 1,
	\"birthdate\": \"1996-10-07\",
	\"phoneNumber\":\"5555555555\",
	\"password\": \"test\"
}"

echo "Done signup"
#login
sessionId=$(curl -X POST \
  http://localhost/v1/login \
  -H 'Accept: */*' \
  -H 'Accept-Encoding: gzip, deflate' \
  -H 'Cache-Control: no-cache' \
  -H 'Connection: keep-alive' \
  -H 'Content-Type: application/json' \
  -H 'Host: localhost' \
  -H 'Postman-Token: cd594573-2621-4604-bc58-a8ab99192ebb,2e2b6c1e-2d4c-40f1-a1f0-f740ead984c1' \
  -H 'User-Agent: PostmanRuntime/7.15.2' \
  -H 'cache-control: no-cache' \
  -d "{\"email\":\"$email\",\"password\":\"test\"}" | jq -r '.Result.sessionId')

echo $sessionId
echo "Done login"

# set the cohort of the user
curl -X POST \
  http://localhost/v1/cohort \
  -H 'Accept: */*' \
  -H 'Accept-Encoding: gzip, deflate' \
  -H 'Cache-Control: no-cache' \
  -H 'Connection: keep-alive' \
  -H 'Content-Type: application/json' \
  -H 'Host: localhost' \
  -H 'Postman-Token: 50d71728-114f-4d08-9934-4af94bcc403f,db2dfb1e-cce9-47e8-87d5-c9cee69b3b5e' \
  -H 'User-Agent: PostmanRuntime/7.15.2' \
  -H 'cache-control: no-cache' \
  -H "sessionId: $sessionId" \
  -d "{\"cohortId\": $cohort,    \"mentorshipPreference\": 0,    \"bio\": \"\",    \"hometown\": \"\"}"
  echo "Done creating"
}

for i in `seq 12`;
do
    firstnamementor=${first[$i]}
    firstnamementee=${first[$i + 12]}
    lastnamementor=${last[$i]}
    lastnamementee=${last[$i + 12]}
    emailmentor=${emails[$i]}
    emailmentee=${emails[$i + 12]}
    createUser $firstnamementor $lastnamementor $emailmentor 25
    createUser $firstnamementee $lastnamementee $emailmentee 44
done