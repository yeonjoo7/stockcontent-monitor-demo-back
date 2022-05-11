ifndef BINARY
	BINARY=debug
endif

GO_BIN_PATH := $(shell go env GOPATH)/bin

init:
	go mod download all
	go install github.com/swaggo/swag/cmd/swag
	go install github.com/google/wire/cmd/wire@v0.5.0
	go install entgo.io/ent/cmd/ent@v0.10.1

gen: swagger wire ent-gen

swagger:
	${GO_BIN_PATH}/swag init

entity:
	${GO_BIN_PATH}/ent init $(name)

entity-gen:
	go generate ./ent

wire:
	${GO_BIN_PATH}/wire .

go-run:
	go run .

go-build:
	go build -ldflags="-X 'stockcontent-monitor-demo-back/build.runtimeProfile=PRODUCTION'" -o ${BINARY} .