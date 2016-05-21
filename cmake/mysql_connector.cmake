# Copyright (C) 2015 Etix Labs - All Rights Reserved.
# All information contained herein is, and remains the property of Etix Labs and its suppliers,
# if any. The intellectual and technical concepts contained herein are proprietary to Etix Labs
# Dissemination of this information or reproduction of this material is strictly forbidden unless
# prior written permission is obtained from Etix Labs.

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
