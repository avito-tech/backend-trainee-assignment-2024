migrate-up:
	migrate -path ./schema -database 'postgres://postgres:qwerty123@localhost:5432/postgres?sslmode=disable' up

migrate-down:
	migrate -path ./schema -database 'postgres://postgres:qwerty123@localhost:5432/postgres?sslmode=disable' down
