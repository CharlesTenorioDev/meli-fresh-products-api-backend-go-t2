#!/bin/bash -i

# Start MySQL in the background.
mysqld --initialize-insecure && mysqld --user=root --init-file=<(echo "ALTER USER 'root'@'localhost' IDENTIFIED BY 'test';") --console &

# Wait for MySQL to be ready.
while ! mysqladmin ping --silent &> /dev/null; do
    sleep 1
done

# Wait for MySQL to accept connections.
while ! mysql -u root -ptest -e "SELECT 1" &> /dev/null; do
    sleep 1
done
echo "SUCCESS TO ACCESS MYSQL"

mysql -u root -ptest -e "CREATE DATABASE fresh_products;"
mysql -u root -ptest fresh_products < ./docs/db/mysql/create.sql

