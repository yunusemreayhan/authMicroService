
createdb:
	docker exec -it postgress_for_auth_micro_service createdb --username admin --owner admin auth_micro_service

dropdb:
	docker exec -it postgress_for_auth_micro_service dropdb --username admin  auth_micro_service

build-sqlc:
	mkdir -p ../build/sqlc
	go build ./sqlc/... -o ../build/sqlc

.PHONY: createdb
