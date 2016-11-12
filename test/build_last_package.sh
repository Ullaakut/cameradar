#!/usr/bin/env bash

ESC_SEQ="\x1b["
COL_RESET=$ESC_SEQ"39;49;00m"
COL_RED=$ESC_SEQ"31;01m"
COL_GREEN=$ESC_SEQ"32;01m"
COL_YELLOW=$ESC_SEQ"33;01m"
COL_BLUE=$ESC_SEQ"34;01m"
COL_MAGENTA=$ESC_SEQ"35;01m"
COL_CYAN=$ESC_SEQ"36;01m"

echo -e $COL_YELLOW"Deleting old package ... "$COL_RESET
rm -f cameradar_*_Debug_Linux.tar.gz
echo -e $COL_GREEN"OK!"$COL_RESET

echo -e $COL_YELLOW"Creating package ... "$COL_RESET
{
  cd ..
  mkdir build
  cd build
  rm -f cameradar_*_Debug_Linux.tar.gz
  cmake .. -DCMAKE_BUILD_TYPE=Debug
  make package
  cp cameradar_*_Debug_Linux.tar.gz ../test
  cd ../test
  } &> /dev/null
echo -e $COL_GREEN"OK!"$COL_RESET
