#!/bin/bash

ls /src/*.zip | xargs -n1 sh -c 'unzip -o $1 -d /src/dump' argv0
ls /src/dump

until mysql -u root --host=127.0.0.1 --password=$MYSQL_ROOT_PASSWORD --wait --connect-timeout=300 -e '\q'; do
  >&2 echo "MySQL is unavailable - sleeping"
  sleep 1
done

result=$(mysql -u root --password=$MYSQL_ROOT_PASSWORD --host=127.0.0.1 -s -N -e "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME='cabs_db'")
if [ -z "$result" ];
then
  mysql -u root --host=127.0.0.1 --password=$MYSQL_ROOT_PASSWORD < /src/createdb.sql
  ls /src/dump/*.sql | xargs -n1 sh -c 'mysql -u root --host="127.0.0.1" --password=$MYSQL_ROOT_PASSWORD cabs_db < $1' argv0
else
  echo "db exists";
fi
