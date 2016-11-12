#!/bin/bash

ports=('8554' '8554' '8554' '8554' '8554' '8554')
users=('admin' 'root' 'ubnt' 'Admin' 'supervisor' '')
passwords=('admin' 'root' '12345' 'ubnt' 'password' '')
routes=('cam0_0' 'live.sdp' 'ch001.sdp' 'cam' 'invalid' 'live_mpeg4.sdp')
cams_name_pattern="fake_camera_"

# json generation variable only
json="[\n"
first=true
# $1 = adress, $2 = port, $3 = path, $4 = usernam $5 = password, $6 = valid
function make_json {
    # Get all data about the container, this will return three lines
    # One empty that we ignore
    # the two other ones with the IP of our container
    # We take the second one using sed and cut to get only the IPAddress
    address="$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $CID)"
    if [ "$first" = true ] ; then first=false
    else json="$json,\n"; fi
    json="$json{"
    json="$json\"address\":\"$address\","
    json="$json\"port\":$2,"
    json="$json\"route\":\"$3\","
    json="$json\"username\":\"$4\","
    json="$json\"password\":\"$5\","
    json="$json\"valid\":$6"
    json="$json}"
}

# $1 = configuration template path
function generate_conf {
    echo "generate configuration"
    sed s#__CAMERAS__#$json#g $1 > cameradartest.conf.json
}

# $1 = numbers of cameras to generate
function start {
    # Seed random generator
    RANDOM=$(date +%s)

    # start cameras
    for (( i=1; i<=$1; i++ )); do
        name="$cams_name_pattern$i"
        # random conf
        conf_idx=$(($RANDOM % ${#ports[@]}))

        # get conf variables
        port=${ports[$conf_idx]}
        user=${users[$conf_idx]}
        passw=${passwords[$conf_idx]}
        route=${routes[$conf_idx]}
        is_valid=true

        # if conf_idx = 4 -> invalid conf
        if [ "$conf_idx" == "4" ] ; then is_valid=false; fi

        CID=$(docker run -d --name "$name" fake-camera /start.sh "$port" "$user" "$passw" "$route");
        make_json "$name" "$port" "$route" "$user" "$passw" $is_valid $CID
    done

    # finalize json
    json="$json]"
}

function stop {
    # if no cameras containers are started just exit
    camera_count="`docker ps -a -q --filter="name=$cams_name_pattern" | wc -l`"
    if [ "$camera_count" == "0" ]; then
        echo "error: no cameras started"; exit 1
    fi

    echo "stopping and removing $camera_count containers"
    # docker stop $(docker ps -a -q --filter="name=$cams_name_pattern")
    docker rm -f $(docker ps -a -q --filter="name=$cams_name_pattern") > /dev/null
}

# need first argument at least
if [ "$1" == "" ]; then
    echo "error: invalid number of argument"
    exit 1
fi
case $1 in
"start")
    # check if the argument is a number.
    re='^[0-9]+$'
    if ! [[ $2 =~ $re ]] ; then
        echo "error: argument is not a number"; exit 1
    fi
    if [[ "$3" == "" ]] ; then
        echo "error: missing path to the configuration file template"; exit 1
    fi
    echo "starting $2 cameras"
    start $2
    generate_conf $3
  ;;
"stop")
    echo "stopping all cameras tests"
    stop
  ;;
"help")
    echo "./gen_cameras.sh start CAMS_NB - start CAMS_NB cameras"
    echo "                 stop - stop all started cameras"
    echo "                 help - display this help"
    exit 0
  ;;
*)
    echo "invalid test name"
    exit 1
  ;;
esac
