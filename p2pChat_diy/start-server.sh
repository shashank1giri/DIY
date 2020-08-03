#!/usr/bin/env bash
#RETRIES=10
#export
#PGPASSWORD=$POSTGRES_PGPASSWORD psql -d $POSTGRES_DB -c "select 1;" > /dev/null 2>&1
#ready=$?
#echo $ready
#echo $ready
#until [ $ready -eq 0 ] || [ $RETRIES -eq 0 ]; do
#  echo "Waiting for postgres server, $((RETRIES--)) remaining attempts..."
#  sleep 1
#  PGPASSWORD=$POSTGRES_PGPASSWORD psql -d $POSTGRES_DB -c "select 1;" > /dev/null 2>&1
#  ready=$?
#  echo $ready
#done
sleep 5
./main
