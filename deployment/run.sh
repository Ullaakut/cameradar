#!/usr/bin/env bash

ESC_SEQ="\x1b["
COL_RESET=$ESC_SEQ"39;49;00m"
COL_RED=$ESC_SEQ"31;01m"
COL_GREEN=$ESC_SEQ"32;01m"
COL_YELLOW=$ESC_SEQ"33;01m"
COL_BLUE=$ESC_SEQ"34;01m"
COL_MAGENTA=$ESC_SEQ"35;01m"
COL_CYAN=$ESC_SEQ"36;01m"

# declare usefuls vars
CONF=/conf/cameradar.conf.json

# copy configuration
cp /tmp/conf/* /conf/

echo -n "replacing cameras subnetworks in configuration "
sed -i s#__CAMERAS_SUBNETWORKS__#$CAMERAS_SUBNETWORKS#g $CONF
echo -e $COL_GREEN"ok"$COL_RESET

echo -n "replacing cameras ports in configuration "
sed -i s#__PORTS_TO_CHECK__#$CAMERAS_PORTS#g $CONF
echo -e $COL_GREEN"ok"$COL_RESET

# Replace ext_cctv_mysql with the IP address of your DB or the name of its Docker
# container. The container has to be linked in docker-compose.yml for cameradar
# to be able to interact with it.
echo -n "replacing mysql host and port in configuration "
sed -i s#__MYSQL_ADDR__#ext_cctv_mysql#g $CONF

# Reaplce 3306 with the port of your DB
sed -i s#__MYSQL_PORT__#3306#g $CONF
echo -e $COL_GREEN"ok"$COL_RESET

/cameradar/bin/cameradar -l 1 -c /conf/cameradar.conf.json &
cameradar_pid=$!

trap 'kill -2 $cameradar_pid; wait $cameradar_pid; exit $?' SIGTERM SIGINT
wait $cameradar_pid
