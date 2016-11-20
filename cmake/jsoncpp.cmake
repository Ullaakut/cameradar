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

cmake_minimum_required(VERSION 2.8.8)
include(ExternalProject)

message(STATUS "Configuring deps.jsoncpp")

set(JSONCPP_DIR jsoncpp)
set(JSONCPP_PATH ${DEPS_DIR}/${JSONCPP_DIR})

ExternalProject_Add(
    deps.jsoncpp
    PREFIX ${JSONCPP_PATH}
    GIT_REPOSITORY https://github.com/open-source-parsers/jsoncpp.git
    TIMEOUT 10
    CONFIGURE_COMMAND ${CMAKE_COMMAND} "-DCMAKE_INSTALL_PREFIX=${JSONCPP_PATH}" -DBUILD_TYPE=Release -DBUILD_STATIC_LIBS=OFF -DBUILD_SHARED_LIBS=ON -DJSONCPP_WITH_TESTS=OFF -DJSONCPP_WITH_POST_BUILD_UNITTEST=OFF <SOURCE_DIR>
    BUILD_IN_SOURCE ON
    UPDATE_COMMAND ""
    BUILD_COMMAND ${CMAKE_MAKE_PROGRAM}
    INSTALL_COMMAND ""
    LOG_DOWNLOAD ON
    LOG_UPDATE ON
    LOG_CONFIGURE ON
    LOG_BUILD ON
)

ExternalProject_Get_Property(deps.jsoncpp SOURCE_DIR)

set (JSONCPP_INCLUDE_DIR "${SOURCE_DIR}/include" PARENT_SCOPE)
set (JSONCPP_LIBRARY_DIR "${SOURCE_DIR}/src/lib_json")
set (JSONCPP_LIBRARY_DIR ${JSONCPP_LIBRARY_DIR} PARENT_SCOPE)

file(GLOB JSONCPP_INSTALL_DEPENDENCIES "${JSONCPP_LIBRARY_DIR}/${CMAKE_SHARED_LIBRARY_PREFIX}jsoncpp${CMAKE_SHARED_LIBRARY_SUFFIX}*")
list (APPEND CAMERADAR_INSTALL_DEPENDENCIES ${JSONCPP_INSTALL_DEPENDENCIES})
set(CAMERADAR_INSTALL_DEPENDENCIES ${CAMERADAR_INSTALL_DEPENDENCIES} PARENT_SCOPE)
