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

project(cameradar CXX)

# find gstreamer 1.x libraries
include(FindPkgConfig)

pkg_search_module(GSTREAMER REQUIRED gstreamer-1.0)
find_library(LIB_GSTREAMER NAMES ${GSTREAMER_LIBRARIES} HINTS ${GSTREAMER_LIBRARY_DIRS})

include_directories (
  ${GSTREAMER_INCLUDE_DIRS}
	${PROJECT_SOURCE_DIR}/include
	${VERSION_INCLUDE_DIR}
)

link_directories (
        ${GSTREAMER_LIBRARY_DIRS}
        "../deps/jsoncpp/src/deps.jsoncpp/src/lib_json"
)

if ("${CMAKE_C_COMPILER_ID}" STREQUAL "Clang")
    # search special osx gstreamer libs
    pkg_search_module(GSTREAMER_APP REQUIRED gstreamer-app-1.0)
    find_library(LIB_GSTREAMER NAMES ${GSTREAMER_APP_LIBRARIES} HINTS ${GSTREAMER_APP_LIBRARY_DIRS})

    include_directories (${GSTREAMER_APP_INCLUDE_DIRS})

    link_directories (${GSTREAMER_APP_LIBRARY_DIRS})
endif()

include (find_sources)
find_sources ("src" "include" "src/models" "src/repositories" "src/tasks")

add_executable (cameradar ${SOURCES})
target_link_libraries (cameradar pthread jsoncpp dl curl ${GSTREAMER_LIBRARIES})

# Add the conf files to the build dir
add_custom_command(TARGET cameradar PRE_BUILD
                   COMMAND ${CMAKE_COMMAND} -E copy_directory
                 ${CMAKE_SOURCE_DIR}/cameradar_standalone/conf $<TARGET_FILE_DIR:cameradar>/conf/)

set (BINARIES_NAME ${PROJECT_BINARY_DIR}/cameradar)
list (APPEND CAMERADAR_BINARIES ${BINARIES_NAME})
set (CAMERADAR_BINARIES ${CAMERADAR_BINARIES} PARENT_SCOPE)
