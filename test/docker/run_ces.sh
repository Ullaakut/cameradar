#!/bin/bash

port=$1
user=$2
passw=$3
route=$4
url=""

# need first argument at least
if [ "$2" == "" ]; then
    url="rtsp://:$port/$route"
else
    url="rtsp://$user:$passw@:$port/$route"
fi

./camera_emulation_server -u $2 -p $3 -r $4
echo "Stream started on ${url}"
