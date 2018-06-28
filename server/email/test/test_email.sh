#!/bin/bash

# SECRETS_PATH=server/secrets.json go run server/email/test/test_email.go --subscribe --forgot_password
SECRETS_PATH=server/secrets.json go run server/email/test/test_email.go --forgotPass
