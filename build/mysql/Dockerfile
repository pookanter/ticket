FROM mysql

COPY build/migration/schema.sql /docker-entrypoint-initdb.d/

ADD build/migration/schema.sql /docker-entrypoint-initdb.d

EXPOSE 3306