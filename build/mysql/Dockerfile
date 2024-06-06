FROM mysql:8.4.0

COPY migration/schema.sql /docker-entrypoint-initdb.d/

ADD migration/schema.sql /docker-entrypoint-initdb.d