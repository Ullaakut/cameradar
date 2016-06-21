#!/bin/bash

# check if a debug package exist in the current folder
if ! ls ./cctv_*_Debug_Linux.tar.gz 1> /dev/null 2>&1; then
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
    cmd="$cmd --link=\"mysql_cameradar\""
    # add cameradar srcs
    cmd="$cmd -v \"`pwd`/src:/go/src/cameradartest\""
    # add cmaeradar conf
    cmd="$cmd -v \"`pwd`/:/tmp/tests\""
    # add container name
    cmd="$cmd cameradartest"
}

function start_test {
    make_docker_command $1
    ./docker/gen_cameras.sh start $1 ./docker/cameratest.conf.tmpl.json
    eval $cmd
    ./docker/gen_cameras.sh stop
}

# build images
echo "building docker images"
# building fake-camera container
docker build -f Dockerfile-camera -t fake-camera .

# building cameradartest image
docker build -t cameradartest .

# getting mysql
echo "starting mysql"
docker pull mysql:5.7
docker run --name mysql_cameradar -e MYSQL_DATABASE=cctv -e MYSQL_ROOT_PASSWORD=root -d mysql:5.7

start_test 1
start_test 5
# start_test 10
# start_test 20

# stop mysql
echo "stopping mysql"
docker rm -f mysql_cameradar