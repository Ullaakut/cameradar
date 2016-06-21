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
./etix_rtsp_server -u $s -p $3 -r $4
# cvlc /vlc/screen.png -I dummy --sout-keep --no-drop-late-frames --no-skip-frames --image-duration 9999 --sout="#transcode{vcodec=h264,fps=15,venc=x264{preset=ultrafast,tune=zerolatency,keyint=30,bframes=0,ref=1,level=30,profile=baseline,hrd=cbr,crf=20,ratetol=1.0,vbv-maxrate=1200,vbv-bufsize=1200,lookahead=0}}:rtp{sdp=$url}" --sout-all
