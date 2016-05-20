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

message(STATUS "Configuring deps.boost")

set(BOOST_VERSION 1.60.0)
# set(BoostSHA1 2fc96c1651ac6fe9859b678b165bd78dc211e881)

# Set up general b2 (bjam) command line arguments
set(b2Args <SOURCE_DIR>/b2
           # link=static
           threading=multi
           runtime-link=shared
           --layout=tagged
           --build-dir=build
           --without-wave
           --without-python
           stage
           -d+2
)

if(TARGET_ARCH STREQUAL "x86_64")
    list(APPEND b2Args address-model=64)
endif()

string(REPLACE "." "_" BOOST_VERSION_UNDERSCORE ${BOOST_VERSION})

set(BOOST_DIR boost)
set(BOOST_PATH ${DEPS_DIR}/${BOOST_DIR})

# Set up build steps
include(ExternalProject)
ExternalProject_Add(
    deps.boost
    PREFIX ${BOOST_PATH}
    URL http://sourceforge.net/projects/boost/files/boost/${BOOST_VERSION}/boost_${BOOST_VERSION_UNDERSCORE}.tar.bz2/download
    TIMEOUT 600
    CONFIGURE_COMMAND ${CMAKE_COMMAND} -E make_directory <SOURCE_DIR>/build
    BUILD_COMMAND "${b2Args}"
    # BUILD_COMMAND "<SOURCE_DIR>/b2 address-model=64 threading=multi runtime-link=shared --layout=tagged --build-dir=<SOURCE_DIR>/build"
    BUILD_IN_SOURCE ON
    INSTALL_COMMAND ""
    # INSTALL_COMMAND <SOURCE_DIR>/b2 install --prefix=${BOOST_PATH}
    LOG_DOWNLOAD ON
    LOG_UPDATE ON
    LOG_CONFIGURE ON
    LOG_BUILD ON
    LOG_TEST ON
    LOG_INSTALL ON
)

# Set extra step to build b2 (bjam)
set(b2Bootstrap "./bootstrap.sh")
ExternalProject_Add_Step(
    deps.boost
    make_b2
    COMMAND ${b2Bootstrap}
    COMMENT "Building b2..."
    DEPENDEES download
    DEPENDERS configure
    WORKING_DIRECTORY <SOURCE_DIR>
    LOG ON
)


ExternalProject_Get_Property(deps.boost SOURCE_DIR)
set(BOOST_INCLUDE_DIR ${SOURCE_DIR} PARENT_SCOPE)
set(BOOST_LIBRARY_DIR "${SOURCE_DIR}/stage/lib")
set(BOOST_LIBRARY_DIR ${BOOST_LIBRARY_DIR} PARENT_SCOPE)

# list all the boost libraries .dylib/.so
file(GLOB BOOST_INSTALL_DEPENDENCIES "${BOOST_LIBRARY_DIR}/${CMAKE_SHARED_LIBRARY_PREFIX}boost_*${CMAKE_SHARED_LIBRARY_SUFFIX}")
list (APPEND CCTV_INSTALL_DEPENDENCIES ${BOOST_INSTALL_DEPENDENCIES})
# on linux
if (CMAKE_SYSTEM_NAME STREQUAL "Linux")
  file(GLOB BOOST_INSTALL_DEPENDENCIES "${BOOST_LIBRARY_DIR}/${CMAKE_SHARED_LIBRARY_PREFIX}boost_*${CMAKE_SHARED_LIBRARY_SUFFIX}.${BOOST_VERSION}")
  list (APPEND CCTV_INSTALL_DEPENDENCIES ${BOOST_INSTALL_DEPENDENCIES})
endif()

set(CCTV_INSTALL_DEPENDENCIES ${CCTV_INSTALL_DEPENDENCIES} PARENT_SCOPE)
