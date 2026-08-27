clean:
	cd app && go clean -modcache

tidy:
	cd app && go mod tidy

deps:
	cd app && go mod download

lint:
	cd app && golangci-lint run ./...
