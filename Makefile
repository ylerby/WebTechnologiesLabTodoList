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
	gofumpt -l -w .
	golangci-lint run
	docker compose up
rebuild:
	docker-compose up --force-recreate --build
full_rebuild:
	go mod tidy
	gofumpt -l -w .
	golangci-lint run
	docker compose up --force-recreate --build