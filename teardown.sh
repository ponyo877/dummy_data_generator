#!/bin/sh

CONTAINER_ID=$(docker ps | grep "postgres:13" | awk '{print $1}')
if [ -n "$CONTAINER_ID" ]; then
    docker exec $CONTAINER_ID /bin/sh -c "psql -d mydb -c \"TRUNCATE table1;TRUNCATE table2;TRUNCATE table3;\""
fi

CONTAINER_ID=$(docker ps | grep "mysql:5.7" | awk '{print $1}')
if [ -n "$CONTAINER_ID" ]; then
    docker exec $CONTAINER_ID /bin/sh -c "MYSQL_PWD=password mysql mydb -e \"TRUNCATE table1;TRUNCATE table2;TRUNCATE table3;\""
fi