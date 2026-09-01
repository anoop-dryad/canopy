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

# ----------------- golang-migrate ---------------------

migrate-up:
	migrate -database "$(DB_DSN)" -path app/migrations up

migrate-down:
	migrate -database "$(DB_DSN)" -path app/migrations down 1

migrate-create:
	migrate create -ext sql -dir app/migrations -seq $(name)
