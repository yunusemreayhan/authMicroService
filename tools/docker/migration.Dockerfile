FROM ubuntu

RUN apt update
RUN apt install -y curl
RUN apt install -y unzip
RUN curl -L -v https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz -o /bin/migrate.tar.gz
RUN tar -xvf /bin/migrate.tar.gz -C /bin/
RUN rm /bin/migrate.tar.gz
RUN rm /bin/README.md
RUN rm /bin/LICENSE

COPY ./db/migration /migration/
RUN chmod +x /bin/migrate
RUN chmod -R 777 /migration

CMD /bin/migrate -database ${SQL_DSN} -path /migration up