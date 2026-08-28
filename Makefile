clean:
	cd app && go clean -modcache

tidy:
	cd app && go mod tidy

deps:
	cd app && go mod download

lint:
	cd app && golangci-lint run ./...

swagger:
	cd app && swag init -g cmd/server/main.go -o infra/http/swagger/docs

run: swagger
	cd app && go run cmd/server/main.go
