clean:
	rm -r bin/
build:
	mkdir -p bin/; go build -tags=local -o bin/scalemetric main.go
local:
	mkdir -p bin/; go build -tags=incluster -o bin/scalemetric main.go
alpine:
	mkdir -p bin; CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=local -o bin/scalemetric main.go
run:
	./bin/scalemetric
debug:
	LOGLEVEL=DEBUG ./bin/scalemetric
docker:
	docker build -t kubestar/scalemetric:$(TAG) .
lambda:
	GOOS=linux go build main.go; zip function.zip main
lambda-clean:
	rm -rf main; rm -rf function.zip
# lambda:
# 	GOOS=linux go build main.go; zip function.zip main; aws lambda create-function --function-name my-function --runtime go1.x \
#   --zip-file fileb://function.zip --handler main \
#   --role arn:aws:iam::xx:role/execution_role
