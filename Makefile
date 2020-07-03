BINARY_NAME=go_weather

.PHONY: test show_coverage run format

all: test build

test:
	go test -v ./... -covermode=count -coverprofile=coverage.out

show_coverage:
	go tool cover -html=coverage.out

build:
	go build -o ${BINARY_NAME} -v main.go

clean:
	${RM} ${BINARY_NAME}
	${RM} bin/${BINARY_NAME}*

run: build
	go run main.go

compile:
	GOOS=linux GOARCH=arm64 go build -o bin/${BINARY_NAME}-linux-arm64 main.go

keycloak:
	docker-compose -f docker/keycloak.yaml up

lint:
	golangci-lint run