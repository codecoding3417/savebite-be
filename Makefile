run:
	@go run cmd/api/main.go
build:
	@go build -o ./build/bitesave ./cmd/api/
migrate:
	@go run cmd/migrate/main.go
