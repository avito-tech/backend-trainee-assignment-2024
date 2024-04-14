wait-postgres:
	@echo "Waiting for PostgreSQL to be ready..."
	@while ! docker-compose exec db pg_isready -U postgers; do \
  		sleep 2; \
 	done

env:
	go mod tidy
	@echo "Generating .env file"
	@if [ -f .env ]; then \
		echo ".env file already exists. Removing old .env file..."; \
		rm .env; \
	fi
	@echo DB_PASSWORD=$$(openssl rand -hex 16) >> .env
	@echo SIGNING_KEY=$$(openssl rand -base64 32 | tr -d /=+ | cut -c1-32) >> .env
	@echo ".env file has been created with random credentials."

run:
	go run cmd/main.go

migrate-up:
	migrate -path ./schema -database 'postgres://postgres:qwerty123@localhost:5432/banner?sslmode=disable' up

migrate-down:
	migrate -path ./schema -database 'postgres://postgres:qwerty123@localhost:5432/banner?sslmode=disable' down

test:
	go clean -testcache
	go test -v ./tests/ | grep -e Test -e FAIL

test-with-docker-db:
	docker compose up -d db
	$(MAKE) wait-postgres
	$(MAKE) test
	docker compose down

docker-db-init:
	docker compose up -d
	$(MAKE) wait-postgres
	$(MAKE) migrate-up

docs-fmt:
	swag fmt

docs-gen:
	swag init -g cmd/main.go
