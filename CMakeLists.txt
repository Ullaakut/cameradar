## Copyright 2016 Etix Labs
##
## Licensed under the Apache License, Version 2.0 (the "License");
## you may not use this file except in compliance with the License.
## You may obtain a copy of the License at
##
##     http://www.apache.org/licenses/LICENSE-2.0
##
## Unless required by applicable law or agreed to in writing, software
## distributed under the License is distributed on an "AS IS" BASIS,
## WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
## See the License for the specific language governing permissions and
## limitations under the License.

cmake_minimum_required (VERSION 2.8.1)
cmake_policy(SET CMP0042 NEW)

# compiler flags
set (CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++14") #enable C++14
set (CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -W -Wall -Wextra -Wno-unused-function") # extra warnings
set (CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fdiagnostics-color") #enable error coloration on gcc

# release specific flags
set (CMAKE_CXX_FLAGS_RELEASE "${CMAKE_CXX_FLAGS_RELEASE} -O2") #enable error coloration on gcc

#debug specific flags
set (CMAKE_CXX_FLAGS_DEBUG "${CMAKE_CXX_FLAGS_DEBUG} -g -fprofile-arcs -ftest-coverage")

# rpath enable osx
set(CMAKE_MACOSX_RPATH 1)

# list of all cache managers
set (CAMERADAR_CACHE_MANAGERS "")

# dependencies directory
set(DEPS_DIR ${CMAKE_SOURCE_DIR}/deps)

# output path for cache managers
set (CCTV_CACHE_MANAGER_OUTPUT_FOLDER cache_managers)
set (CCTV_CACHE_MANAGER_OUTPUT_PATH ${CMAKE_BINARY_DIR}/${CCTV_CACHE_MANAGER_OUTPUT_FOLDER})

set(CMAKE_MODULE_PATH ${PROJECT_SOURCE_DIR}/cmake)

include_directories (
        "cameradar_standalone/include"
        "deps/jsoncpp/src/deps.jsoncpp/include"
)

#build cache managers
add_subdirectory (deps)
add_subdirectory (cameradar_standalone)
add_subdirectory (cache_managers)
