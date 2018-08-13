#!/bin/bash

case "$1" in
	-u|--update)
		UPDATE=1;
		shift
	;;
	-n|--new)
		NEW=1;
		shift
	;;
	*)
	break
	;;
esac

FUNCTION_NAME=$1

# build the function
compile() {
	cd $FUNCTION_NAME && GOOS=linux go build -o main
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
	  --role arn:aws:iam::016267150191:role/LambdaRole \
	  --runtime go1.x \
	  --zip-file fileb://main.zip \
	  --handler main
}

update() {
	aws lambda update-function-code \
		--function-name $1 \
		--zip-file fileb://main.zip
}

echo "Compiling"
compile $FUNCTION_NAME

echo "Compressing"
compress $FUNCTION_NAME

if ! [[ -z $UPDATE ]]; then
	echo "Updating"
	update $FUNCTION_NAME
elif ! [[ -z $NEW ]]; then
	echo "Uploading"
	upload $FUNCTION_NAME
else
	echo "NOOP"
fi

if [[ "$?" = "1" ]]; then
	echo "Maybe try updating the function with flag --update"
	exit 1
fi

echo "Cleaning up"
rm main
rm main.zip
