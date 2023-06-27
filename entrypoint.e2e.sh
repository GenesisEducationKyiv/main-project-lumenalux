#!/bin/sh
if [ ! -f /app/storage/storage.csv ]; then
    touch /app/storage/storage.csv
else
    truncate -s 0 /app/storage/storage.csv
fi

./gses2-app