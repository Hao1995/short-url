
install:
	# for enum
	go install github.com/abice/go-enum@v0.6.0
	# for mocks
	go install github.com/vektra/mockery/v2@v2.52.1
	# for database migration
	go install github.com/pressly/goose/v3/cmd/goose@v3.24.1

generate:
	# gen enum
	go generate ./...
	# gen mocks
	mockery

test:
	go test --race ./...

up:
	docker compose up --build -d

down:
	docker compose down