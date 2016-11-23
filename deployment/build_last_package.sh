#!/usr/bin/env bash

ESC_SEQ="\x1b["
COL_RESET=$ESC_SEQ"39;49;00m"
COL_RED=$ESC_SEQ"31;01m"
COL_GREEN=$ESC_SEQ"32;01m"
COL_YELLOW=$ESC_SEQ"33;01m"

echo -e $COL_YELLOW"Deleting old package ... "$COL_RESET
rm -f cameradar_*_${1:-"Release"}_Linux.tar.gz
echo -e $COL_GREEN"OK!"$COL_RESET

echo -e $COL_YELLOW"Creating package ... "$COL_RESET

cd ..
ret=$?
if [ "$ret" -ne "0" ]; then
  echo -e $COL_RED"KO!"$COL_RESET;
  exit 1;
fi

mkdir build

cd build
ret=$?
if [ "$ret" -ne "0" ]; then
  echo -e $COL_RED"KO!"$COL_RESET;
  exit 1;
fi

rm -f cameradar_*_${1:-"Release"}_Linux.tar.gz

cmake .. -DCMAKE_BUILD_TYPE=${1:-"Release"}
ret=$?
if [ "$ret" -ne "0" ]; then
  echo -e $COL_RED"KO!"$COL_RESET;
  exit 1;
fi

make package
ret=$?
if [ "$ret" -ne "0" ]; then
  echo -e $COL_RED"KO!"$COL_RESET;
  exit 1;
fi

cp cameradar_*_${1:-"Release"}_Linux.tar.gz ../deployment

cd ../deployment
ret=$?
if [ "$ret" -ne "0" ]; then
  echo -e $COL_RED"KO!"$COL_RESET;
  exit 1;
fi
echo -e $COL_GREEN"OK!"$COL_RESET
