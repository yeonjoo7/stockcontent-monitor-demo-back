ifndef BINARY
	BINARY=debug
endif

GENERATE_PATH := $(shell go env GOPATH)/bin

init:
	go mod download all
	go install github.com/swaggo/swag/cmd/swag
	go install github.com/google/wire/cmd/wire@v0.5.0

generate: swagger wire

swagger:
	${GENERATE_PATH}/swag init

wire:
	${GENERATE_PATH}/wire .