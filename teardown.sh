#!/bin/sh

CONTAINER_ID=$(docker ps | grep "postgres:13" | awk '{print $1}')
if [ -n "$CONTAINER_ID" ]; then
    docker exec $CONTAINER_ID /bin/sh -c "psql -d mydb -c \"TRUNCATE table1;TRUNCATE table2;\""
fi