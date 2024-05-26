FROM mysql

COPY migration/schema.sql /docker-entrypoint-initdb.d/

ADD migration/schema.sql /docker-entrypoint-initdb.d

EXPOSE 3306