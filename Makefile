.DEFAULT_GOAL := build

APP?=test-tcp
RELEASE?=0.0.1
# COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?= $(shell date '+%Y%m%d%H%M')
PROJECT?=github.com/Juev/test-tcp

clean:
	rm -f ${APP}

build: clean
	go build \
		-ldflags "-s -w -X ${PROJECT}/version.Release=${RELEASE} \
		-X ${PROJECT}/version.BuildTime=${BUILD_TIME}" \
		-o ${APP}

run: build
	./${APP} -c data/config.toml

test:
	go test -v -race ./..