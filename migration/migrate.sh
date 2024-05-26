#!/bin/bash

# Name of the Docker service running MySQL
SERVICE_NAME="mysql"

# Execute the mysql command in the Docker container to import the schema.sql file
docker-compose exec -T $SERVICE_NAME mysql -u $MYSQL_ROOT_USERNAME -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE < /usr/local/bin/schema.sql