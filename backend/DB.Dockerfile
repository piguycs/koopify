FROM docker.io/library/postgres:18.1-alpine3.23

ENV POSTGRES_DB=kopify_dev_db \
    POSTGRES_USER=dev_user \
    POSTGRES_PASSWORD=postgres_dev_123

COPY ./sql/schema.sql /docker-entrypoint-initdb.d/

EXPOSE 5432
