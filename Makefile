all: deploy.zip

.PHONY: aria-alexa
aria-alexa:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

deploy.zip: aria-alexa
	zip deploy.zip aria-alexa

.PHONY: clean
clean:
	rm -f aria-alexa deploy.zip

.PHONY: deploy
deploy: deploy.zip
	aws lambda update-function-code --function-name ${LAMBDA_FUNCTION_ARN} --zip-file fileb://deploy.zip
