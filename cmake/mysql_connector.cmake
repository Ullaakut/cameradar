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

# MySQL Connector dependency
message(STATUS "Configuring deps.mysqlconnector")

set (MYSQL_CONNECTOR_VERSION 1.1.6)
set (MD5 9e49dcfc1408b18b3d3ca02781ff7efb)
set (MYSQL_CONNECTOR_DIR mysql-connector)
set (MYSQL_CONNECTOR_PATH ${DEPS_DIR}/${MYSQL_CONNECTOR_DIR})

set (BOOST_ROOT_DIR ${DEPS_DIR}/boost/src/deps.boost)

# include(ExternalProject)
ExternalProject_Add(
    deps.mysql_connector
    PREFIX ${MYSQL_CONNECTOR_PATH}
    URL http://dev.mysql.com/get/Downloads/Connector-C++/mysql-connector-c++-${MYSQL_CONNECTOR_VERSION}.tar.gz
    URL_HASH MD5=${MD5}
    CONFIGURE_COMMAND ${CMAKE_COMMAND} -DBOOST_ROOT=${BOOST_ROOT_DIR} "-DCMAKE_INSTALL_PREFIX=${MYSQL_CONNECTOR_PATH}" -DBUILD_TYPE=Release -DMYSQL_CXXFLAGS=-fexceptions <SOURCE_DIR>
    BUILD_IN_SOURCE ON
    UPDATE_COMMAND ""
    BUILD_COMMAND ${CMAKE_MAKE_PROGRAM}
    INSTALL_COMMAND ${CMAKE_MAKE_PROGRAM} install
    LOG_DOWNLOAD ON
    LOG_UPDATE ON
    LOG_CONFIGURE ON
    LOG_BUILD ON
)

set (MYSQL_CONNECTOR_INCLUDE_DIR "${MYSQL_CONNECTOR_PATH}/include" PARENT_SCOPE)
set (MYSQL_CONNECTOR_LIBRARY_DIR "${MYSQL_CONNECTOR_PATH}/lib")

set (MYSQL_CONNECTOR_LIBRARY_DIR ${MYSQL_CONNECTOR_LIBRARY_DIR} PARENT_SCOPE)

# list all the hiredis libraries
file(GLOB MYSQL_CONNECTOR_INSTALL_DEPENDENCIES "${MYSQL_CONNECTOR_LIBRARY_DIR}/${CMAKE_SHARED_LIBRARY_PREFIX}mysqlcppconn*${CMAKE_SHARED_LIBRARY_SUFFIX}*")
list (APPEND CAMERADAR_INSTALL_DEPENDENCIES ${MYSQL_CONNECTOR_INSTALL_DEPENDENCIES})

set(CAMERADAR_INSTALL_DEPENDENCIES ${CAMERADAR_INSTALL_DEPENDENCIES} PARENT_SCOPE)
