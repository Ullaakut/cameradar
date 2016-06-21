#!/bin/bash

while ! mysqladmin ping -h"mysql_cameradar" -P3306 --silent; do
    sleep 1
done

ls -alhR /conf
cat /etc/hosts

# build
go build
# run test
./cameradartest /tmp/tests/cameradartest.conf.json

cp cameratest.log.xml /tmp/tests/