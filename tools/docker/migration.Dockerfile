FROM ubuntu

COPY ./tools/migrate /bin/
COPY ./db/migration /migration/
RUN chmod +x /bin/migrate
RUN chmod -R 777 /migration

CMD /bin/migrate -database "postgres://root:root@auth_micro_service_db:5431/auth_micro_service?sslmode=disable" -path /migration up