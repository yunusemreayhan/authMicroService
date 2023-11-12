
create_pg_for_db:
	docker run --name postgress_for_auth_micro_service -e POSTGRES_USER='root' -e POSTGRES_PASSWORD='root' -e PGPORT=${PGPORT} -p ${PGPORT}:${PGPORT} -d postgres:12-alpine
	sleep 3
	docker exec -it postgress_for_auth_micro_service createdb --username root --owner root auth_micro_service
	cd db
	make migrate_up
	cd ..

remove_pg_for_db:
	docker stop postgress_for_auth_micro_service || echo
	docker rm -f postgress_for_auth_micro_service || echo

logs_pg_for_db:
	docker logs postgress_for_auth_micro_service

bash_for_db:
	docker exec -it postgress_for_auth_micro_service bash

migrate_up:
	migrate -path ./db/migration/ -database ${SQL_DSN} -verbose up

migrate_down:
	migrate -path ./db/migration/ -database ${SQL_DSN} -verbose down

init_migration:
	migrate create --ext sql --dir db/migration/ --seq init_schema

reset_test_db: remove_pg_for_db create_pg_for_db

sqlc_generate:
	sqlc generate

unit_test:
	go clean -cache
	PGPORT=5435 SQL_DSN=postgresql://root:root@localhost:5435/auth_micro_service?sslmode=disable make reset_test_db
	PGPORT=5435 SQL_DSN=postgresql://root:root@localhost:5435/auth_micro_service?sslmode=disable go test -v ./...
	make remove_pg_for_db

integration_test:
	py.test-3 --capture=no ./test/integration/simple_auth_calls.py

generate: sqlc_generate

clean:
	rm -rf ./db/sqlc
	rm -rf ./build

key_generator:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags="all=-N -l"  -v  -o ./build/key_generator ./cmd/keygen/main.go

keys_for_auth:
	make key_generator
	./build/key_generator -key="./build/private.key"

auth_micro_service_binary:
	make generate
	mkdir -p ./build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags="all=-N -l"  -v  -o ./build/auth_micro_service ./cmd/auth/main.go

auth_micro_service_docker_image:
	docker image rm -f auth_micro_service
	docker build -f ./tools/docker/auth.Dockerfile -t auth_micro_service .

auth_micro_service: auth_micro_service_binary keys_for_auth

start_docker:
	docker rm -f auth_micro_service_instance
	docker run -t --name auth_micro_service_instance -p 3000:3000 -d auth_micro_service
	docker logs auth_micro_service_instance

all:
	./tools/install_local_dependencies.py
	make auth_micro_service 
	./tools/download_migrate.py

compose:
	make all
	make remove_pg_for_db
	docker compose up --detach	

compose_test:
	make all
	docker compose -f ./compose_test.yaml up --detach	

clean_compose_all:
	@echo -n "Are you sure about removing production database? [y/N] " && read ans && [ $${ans:-N} = y ]
	docker compose down -t 15 --remove-orphans --volumes
	sleep 5
	docker rmi authmicroservice-auth_micro_service_backend || echo "auth_micro_service_backend is not running"
	docker rmi authmicroservice-migration || echo "migration is not running"
	docker volume rm authmicroservice_db-data || echo "auth_micro_service_db_data is not running"

clean_compose:
	docker compose down -t 15 --remove-orphans --rmi "local"
	sleep 5
	docker ps -a
	docker image ls

clean_compose_test:
	docker compose down -t 15 --remove-orphans --rmi "local"
	sleep 5
	docker ps -a
	docker image ls
