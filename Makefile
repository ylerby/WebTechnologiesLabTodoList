run:
	docker compose up
lint:
	golangci-lint run
stop:
	docker compose down
format:
	gofumpt -l -w .
full_run:
	go mod tidy
	golangci-lint run
	gofumpt -l -w .
	docker compose up
