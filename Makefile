
swag-init: ## init swagger
	swag init -g cmd/main.go -o api/docs

run: ## run application
	go run cmd/main.go