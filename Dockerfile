FROM alpine:latest
#FROM golang:1.12.10
COPY bin/scalemetric /scalemetric
WORKDIR /
EXPOSE 8080
ENTRYPOINT ["/scalemetric"]