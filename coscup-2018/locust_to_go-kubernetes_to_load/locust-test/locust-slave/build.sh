#!/bin/sh

go build -v -o locust-slave locust-slave.go
docker build -t locust-slave .
