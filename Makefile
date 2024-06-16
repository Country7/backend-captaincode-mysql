# DB_URL=postgresql://root:secret@localhost:5432/main_db?sslmode=disable
# DB_URL=mysql://root:secret@tcp(localhost:3306)/main_db
# jdbc:mysql://localhost:3306/db?allowPublicKeyRetrieval=true&useSSL=false
DB_URL=mysql://root:secret@tcp(localhost:3306)/main_db

help: ## Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST) | column -tl 2

pull-images: ## Pull the needed docker images. (docker pull postgres:16-alpine)
	docker pull mysql:8.0

create-network: ## Create the ww-network.
	docker network create ww-network

create-db: ## Create the database. (docker exec -it postgres16 createdb --username=root --owner=root main_db)
	docker exec -it mysql8 createdb --username=root --owner=root main_db

drop-db: ## Drop the database.  (docker exec -it postgres16 dropdb main_db)
	docker exec -it mysql8 dropdb main_db

migrate-up: ## Apply all up migrations.
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrate-up-1: ## Apply the last up migration.
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migrate-down: ## Apply all down migrations.
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migrate-down-1: ## Apply the last down migration.
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

# https://hub.docker.com/_/postgres
# run-postgres: ## Run postgresql database docker image.
#	docker run --name postgres16 --network ww-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine
run-mysql8: ## Run mysql database docker image.
	docker run --name mysql8 -p 3306:3306 -e MYSQL_DATABASE=main_db -e MYSQL_ROOT_PASSWORD=secret -d mysql:8.0

#start-postgres16: ## Start available postgresql database docker container.
#	docker start postgres16
start-mysql8: ## Start available mysql database docker container.
	docker start mysql8

#stop-postgres: ## Stop postgresql database docker image.
#	docker stop postgres16
stop-mysql8: ## Stop mysql database docker image.
	docker stop mysql8

#run-postgres-cli:    ## Run psql on the postgres16 docker container.
#	docker exec -it -u root postgres16 psql

db-docs: ## Generate the database documentation.
	dbdocs build doc/db.dbml

#db-schema: ## Generate the database schema.
#	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc: ## sqlc generate.
	sqlc generate

docker-system-clean: ## Docker system clean.
	docker system prune -f

test: ## Test go files and report coverage.
	go test -v -cover ./...

server: ## Run the application server.
	go run main.go

mock: ## Generate a store mock.
	mockgen -package mockdb -destination db/mock/store.go github.com/Country7/backend-captaincode-mysql/db/sqlc Store

build-docker-image: ## Build the Docker image.
	docker build -t backend-captaincode:latest .

.PHONY: run-mysql8 start-mysql8 stop-mysql8 create-db drop-db migrate-up migrate-down \
 run-postgres-cli docker-system-clean sqlc test mock migrate-up-1 migrate-down-1 db-docs db-schema
