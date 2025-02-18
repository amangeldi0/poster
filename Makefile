.PHONY: default
default: run;

generate:
	sqlc -f ./config/sqlc.yaml generate
m-up:
	goose -env=./config/.env up
m-down:
	goose -env=./config/.env down
run:
	go run cmd/poster/main.go
dbc:
	rm -rf ./internal/database
