# LetsTalkMessenger

The main service to handle user registration, message handling.

Note: I set this up as a intellij project so you should just be able to import it.

## Building
`sbt compile .`

## Running
`sbt run`

## Sample Request
curl -X POST http://localhost:8080/messages/send \
    -H "Content-Type: application/json" \
    -d"{\"from\":\"acod\", \"to\":\"andrew\", \"payload\":\"Hello World\"}"
