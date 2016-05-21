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

project(mysql_cache_manager CXX)

find_package(PkgConfig)

include_directories (${PROJECT_SOURCE_DIR}/include ${CAMERADAR_INCLUDES})

include (find_sources)
find_sources ("src" "include")

add_library (mysql_cache_manager SHARED ${SOURCES})
set(CMAKE_SHARED_LINKER_FLAGS "-Wl,--no-undefined")
target_link_libraries (mysql_cache_manager jsoncpp mysqlcppconn pthread)

set (CACHE_MANAGER_NAME ${CAMERADAR_CACHE_MANAGER_OUTPUT_PATH}/${CMAKE_SHARED_LIBRARY_PREFIX}mysql_cache_manager${CMAKE_SHARED_LIBRARY_SUFFIX})
list (APPEND CAMERADAR_CACHE_MANAGERS ${CACHE_MANAGER_NAME})
set (CAMERADAR_CACHE_MANAGERS ${CAMERADAR_CACHE_MANAGERS} PARENT_SCOPE)
