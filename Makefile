start:
	go run cmd/api/main.go

swagger-gen:
	swag init -g cmd/api/main.go --parseDependency --parseInternal

test:
	go test ./... -v

mock-gen-all:
	go generate ./...
