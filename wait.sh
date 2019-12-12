#!/bin/bash

cmd="$@"

ls /src/*.zip | xargs -n1 sh -c 'unzip -o $1 -d /src/dump' argv0
ls /src/dump

until mysql -u root --host=127.0.0.1 --password=$MYSQL_ROOT_PASSWORD --wait --connect-timeout=300 -e '\q'; do
  >&2 echo "MySQL is unavailable - sleeping"
  sleep 1
done

result=$(mysql -u root --password=$MYSQL_ROOT_PASSWORD --host=127.0.0.1 -s -N -e "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME='cabs_db'")
while [ "$result" != "cabs_db" ]; do
    echo "waiting for DB"
    sleep 1
    result=$(mysql -u root --password=$MYSQL_ROOT_PASSWORD --host=127.0.0.1 -s -N -e "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME='cabs_db'")
done
echo "db exists";


result=$(mysql -u root --password=$MYSQL_ROOT_PASSWORD --host=127.0.0.1 -s -N -e "SELECT count(*) FROM cabs_db.cab_trip_data")
while [ $result -lt 1 ]; do
    echo "waiting for data"
    sleep 1
    result=$(mysql -u root --password=$MYSQL_ROOT_PASSWORD --host=127.0.0.1 -s -N -e "SELECT count(*) FROM cabs_db.cab_trip_data")
done
echo "data present";

exec $cmd
