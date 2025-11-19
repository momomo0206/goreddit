.PHONY: postgres adminer migrate

postgres:
	docker run --rm -d -p 5432:5432 --name postgres_db -e POSTGRES_PASSWORD=secret postgres

adminer:
	docker run --rm -d -p 8080:8080 --name adminer_app -e ADMINER_DEFAULT_SERVER=host.docker.internal adminer

migrate:
	migrate -source file://migrations \
					-database postgres://postgres:secret@localhost:5432/postgres?sslmode=disable up

migrate-down:
	migrate -source file://migrations \
					-database postgres://postgres:secret@localhost:5432/postgres?sslmode=disable down