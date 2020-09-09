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