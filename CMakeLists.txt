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
cmake_policy(SET CMP0048 OLD)

set (PROJECT_NAME cameradar)

project (${PROJECT_NAME})

message ("Here")

set (${PROJECT_NAME}_VERSION_MAJOR 1)
set (${PROJECT_NAME}_VERSION_MINOR 1)
set (${PROJECT_NAME}_VERSION_PATCH 0)
set (${PROJECT_NAME}_VERSION "${${PROJECT_NAME}_VERSION_MAJOR}.${${PROJECT_NAME}_VERSION_MINOR}.${${PROJECT_NAME}_VERSION_PATCH}${${PROJECT_NAME}_SUFFIX}")

find_package(Git REQUIRED)

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
set (CAMERADAR_CACHE_MANAGER_OUTPUT_FOLDER cache_managers)
set (CAMERADAR_CACHE_MANAGER_OUTPUT_PATH ${CMAKE_BINARY_DIR}/${CAMERADAR_CACHE_MANAGER_OUTPUT_FOLDER})

set(CMAKE_MODULE_PATH ${PROJECT_SOURCE_DIR}/cmake)

# the place where the version.h file is generated, used from the main.cpp of cameradar
set (VERSION_INCLUDE_DIR ${PROJECT_BINARY_DIR})

# get the git revision
message (STATUS "retrieve current git revision SHA1 of cameradar")
execute_process(
    COMMAND "git" "rev-parse" "HEAD"
    WORKING_DIRECTORY ${PROJECT_SOURCE_DIR}
    OUTPUT_VARIABLE CAMERADAR_GIT_SHA1
)

# remove last character of the git output (\n)
string(LENGTH ${CAMERADAR_GIT_SHA1} CAMERADAR_GIT_SHA1_LEN)
math(EXPR CAMERADAR_GIT_SHA1_LEN "${CAMERADAR_GIT_SHA1_LEN} - 1")
string(SUBSTRING ${CAMERADAR_GIT_SHA1} 0 ${CAMERADAR_GIT_SHA1_LEN} CAMERADAR_GIT_SHA1)

# print the SHA1
message (STATUS "current cameradar git revision SHA1 is ${CAMERADAR_GIT_SHA1}")

# generate build number from the current timestamp
string(TIMESTAMP CAMERADAR_VERSION_BUILD "%Y%m%d%H%M%S" "UTC")

# print version
message (STATUS "current cameradar build version will be ${CAMERADAR_VERSION_BUILD}")

configure_file (
    "${PROJECT_SOURCE_DIR}/version.h.in"
    "${PROJECT_BINARY_DIR}/version.h"
)

# add all deps libraries to the link directories path
link_directories (
    # third party libraries
    "deps/jsoncpp/src/deps.jsoncpp/src/lib_json"
    "deps/boost/src/deps.boost/libs"
    "deps/mysql-connector/lib"
)

include_directories (
    "cameradar_standalone/include"
    "deps/jsoncpp/src/deps.jsoncpp/include"
    "deps/boost/src/deps.boost/include"
    "deps/mysql-connector/include"
)

set (${CAMERADAR_BINARIES} "")
set (${CAMERADAR_LIBRARIES} "")

#build cache managers
add_subdirectory (deps)
message ("Debug")
add_subdirectory (cameradar_standalone)
add_subdirectory (cache_managers)

list (APPEND CAMERADAR_LIBRARIES ${CAMERADAR_INSTALL_DEPENDENCIES} ${CAMERADAR_LIBRARIES})

install (PROGRAMS ${CAMERADAR_BINARIES} DESTINATION bin)
install (FILES ${CAMERADAR_CACHE_MANAGERS} DESTINATION cache_managers)
install (FILES ${CAMERADAR_LIBRARIES} DESTINATION libraries)
install (DIRECTORY ${CMAKE_SOURCE_DIR}/deps/licenses DESTINATION libraries)

# cpack configuration
include (InstallRequiredSystemLibraries)
set (CPACK_PACKAGE_DESCRIPTION_SUMMARY "cameradar")
set (CPACK_PACKAGE_VENDOR "Etix Labs")
set (CPACK_PACKAGE_DESCRIPTION_SUMMARY "cameradar tool")
set (CPACK_PACKAGE_FILE_NAME "${PROJECT_NAME}_${${PROJECT_NAME}_VERSION}_${CMAKE_BUILD_TYPE}_${CMAKE_SYSTEM_NAME}")
set (CPACK_PACKAGE_DESCRIPTION_FILE "${CMAKE_CURRENT_SOURCE_DIR}/README.md")
set (CPACK_RESOURCE_FILE_LICENSE "${CMAKE_CURRENT_SOURCE_DIR}/LICENSE")
set (CPACK_PACKAGE_VERSION_MAJOR "0")
set (CPACK_PACKAGE_VERSION_MINOR "2")
set (CPACK_PACKAGE_VERSION_PATCH "2")
set (CPACK_PACKAGE_INSTALL_DIRECTORY "${PROJECT_NAME}_${${PROJECT_NAME}_VERSION}")
set (CPACK_GENERATOR "TGZ")
set (CPACK_SOURCE_GENERATOR "TGZ")

include(CPack)
