#!/usr/bin/env bash
# start-server.sh
cd diy
if [[ "$MIGRATE" == "true" ]]
then
	(python manage.py makemigrations;python manage.py migrate)
fi
(gunicorn diy.wsgi --user www-data --bind 0.0.0.0:8010 --workers 3) &
nginx -g "daemon off;"
