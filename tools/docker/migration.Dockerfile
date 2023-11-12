FROM ubuntu

COPY ./tools/migrate /bin/
COPY ./db/migration /migration/
RUN chmod +x /bin/migrate
RUN chmod -R 777 /migration

CMD /bin/migrate -database ${SQL_DSN} -path /migration up