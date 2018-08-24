REGION=us-east-1
LAMBDA=ExampleQuizSkill

SRC=intents.go lambda.go quiz.go

test: main.zip
	 sam local invoke --event input/tell_california.json

create: main.zip
	aws lambda create-function \
		--region $(REGION) \
		--function-name $(LAMBDA) \
		--role arn:aws:iam::110049094078:role/lambda_basic_execution \
		--runtime go1.x \
		--handler main \
		--zip-file fileb://`pwd`/main.zip

update: main.zip
	aws lambda update-function-code \
		--region $(REGION) \
		--function-name $(LAMBDA) \
		--zip-file fileb://`pwd`/main.zip

main.zip: main
	zip main.zip main

main: $(SRC)
	GOOS=linux go build -o main $(SRC)

clean:
	rm -f main main.zip
