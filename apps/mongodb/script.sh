#!/bin/sh
set -e
mongoimport --host localhost --db test --collection users --drop --file /app/data/users_data.json
