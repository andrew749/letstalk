#!/bin/bash

SECRETS_PATH=server/secrets.json go run server/email/test/test_email.go --subscribe --forgotPass
