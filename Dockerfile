FROM alpine:latest
#FROM golang:1.12.10
COPY bin/cronjobwatch /cronjobwatch
WORKDIR /
ENTRYPOINT ["/cronjobwatch"]
