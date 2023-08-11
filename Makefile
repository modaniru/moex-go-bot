.PHONY: run
run: fmt
	go run cmd/main.go

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: mup
mup:
	migrate -path db/migrations -database 'postgres://postgres:postgres@localhost:1111/postgres?sslmode=disable' up

.PHONY: mdown
mdown:
	migrate -path db/migrations -database 'postgres://postgres:postgres@localhost:1111/postgres?sslmode=disable' down