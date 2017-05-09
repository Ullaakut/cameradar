#!/bin/bash

ESC_SEQ="\x1b["
COL_RESET=$ESC_SEQ"39;49;00m"
COL_GREEN=$ESC_SEQ"32;01m"

# if command starts with an option, prepend /cameradar/bin/cameradar
if [ "${1:0:1}" = '-' ]; then
	set -- /cameradar/bin/cameradar "$@"
fi

# skip setup if they want an option that stops cameradar
wantHelp=
for arg; do
	case "$arg" in
		-v|-h)
			wantHelp=1
			break
			;;
	esac
done

if [ "$CACHE_MANAGER" == "" ]; then
	export CACHE_MANAGER="dumb"
fi

if [ "$CAMERAS_TARGET" == "" ]; then
	export CAMERAS_TARGET="0.0.0.0"
fi

if [ "$CAMERAS_PORTS" == "" ]; then
	export CAMERAS_PORTS="554,8554"
fi

envsubst < /cameradar/conf/cameradar.tmpl.conf.json > /cameradar/conf/cameradar.conf.json


if [ "$CACHE_MANAGER" == "mysql" ] && [ "$1" = '/cameradar/bin/cameradar' -a -z "$wantHelp" ]; then
  echo -n "Waiting for cameradar-database to be ready..."
  while ! mysqladmin ping -h "cameradar-database" -P3306 --silent; do
      sleep 1; echo -n "."
  done
  echo -e $COL_GREEN"ok"$COL_RESET

  echo "Cameradar init finished. Starting it."
fi

exec "$@"
