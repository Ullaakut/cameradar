#!/bin/bash

# check if a debug package exist in the current folder
if ! ls ./cameradar_*_Debug_Linux.tar.gz 1> /dev/null 2>&1; then
    (echo "no debug package in the current folder"; exit 137)
    exit 137
fi

cams_name_pattern="fake_camera_"
cmd=""

function make_docker_command {
    cmd="docker run --rm"

    # start cameras
    for (( i=1; i<=$1; i++ )); do
        name="$cams_name_pattern$i"
        cmd="$cmd --link=\"$name\""
    done

    # add mysql libk
    cmd="$cmd --link=\"cameradar-database\""
    # add cameradar srcs
    cmd="$cmd -v \"`pwd`/src:/go/src/cameradartest\""
    # add cmaeradar conf
    cmd="$cmd -v \"`pwd`/:/tmp/tests\""
    # add container name
    cmd="$cmd -v \"`pwd`/:/tmp/shared\""
    # add container name
    cmd="$cmd cameradartest"
}

function start_test {
    ./docker/gen_cameras.sh start $1 ./docker/cameratest.conf.tmpl.json
    make_docker_command $1
    eval $cmd
    ret=$?
    ./docker/gen_cameras.sh stop
    return $ret
}

# build images
echo "building docker images"
# building fake-camera container
docker build --no-cache -f Dockerfile-camera -t fake-camera .

# building cameradartest image
docker build --no-cache -t cameradartest .

# getting mysql
echo "starting mysql"
docker pull mysql:5.7
docker run --name cameradar-database -e MYSQL_DATABASE=cmrdr -e MYSQL_ROOT_PASSWORD=root -d mysql:5.7

start_test 5
ret=$?
echo "Tests returned ${ret}"

# stop mysql
echo "stopping mysql"
docker rm -f cameradar-database
exit $ret
