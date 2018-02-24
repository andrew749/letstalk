#!/bin/bash

FUNCTION_NAME=$1

# build the function
compile() {
	cd $1 && GOOS=linux go build -o main
}

# package the function
compress() {
	zip main.zip main
	chmod u+x main.zip
}

# upload the function
upload() {
	aws lambda create-function \
	  --region us-east-1 \
	  --function-name $1 \
	  --memory 128 \
	  --role arn:aws:iam::947945882937:role/lambda_user \
	  --runtime go1.x \
	  --zip-file fileb://main.zip \
	  --handler main
}

echo "Compiling"
compile $FUNCTION_NAME

echo "Compressing"
compress $FUNCTION_NAME

echo "Uploading"
upload $FUNCTION_NAME

echo "Cleaning up"
rm main
rm main.zip
