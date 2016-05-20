if("master" STREQUAL "")
  message(FATAL_ERROR "Tag for git checkout should not be empty.")
endif()

set(run 0)

if("/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp-stamp/deps.jsoncpp-gitinfo.txt" IS_NEWER_THAN "/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp-stamp/deps.jsoncpp-gitclone-lastrun.txt")
  set(run 1)
endif()

if(NOT run)
  message(STATUS "Avoiding repeated git clone, stamp file is up to date: '/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp-stamp/deps.jsoncpp-gitclone-lastrun.txt'")
  return()
endif()

execute_process(
  COMMAND ${CMAKE_COMMAND} -E remove_directory "/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp"
  RESULT_VARIABLE error_code
  )
if(error_code)
  message(FATAL_ERROR "Failed to remove directory: '/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp'")
endif()

# try the clone 3 times incase there is an odd git clone issue
set(error_code 1)
set(number_of_tries 0)
while(error_code AND number_of_tries LESS 3)
  execute_process(
    COMMAND "/usr/bin/git" clone "https://github.com/open-source-parsers/jsoncpp.git" "deps.jsoncpp"
    WORKING_DIRECTORY "/home/ullaakut/Work/cameradar/deps/jsoncpp/src"
    RESULT_VARIABLE error_code
    )
  math(EXPR number_of_tries "${number_of_tries} + 1")
endwhile()
if(number_of_tries GREATER 1)
  message(STATUS "Had to git clone more than once:
          ${number_of_tries} times.")
endif()
if(error_code)
  message(FATAL_ERROR "Failed to clone repository: 'https://github.com/open-source-parsers/jsoncpp.git'")
endif()

execute_process(
  COMMAND "/usr/bin/git" checkout master
  WORKING_DIRECTORY "/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp"
  RESULT_VARIABLE error_code
  )
if(error_code)
  message(FATAL_ERROR "Failed to checkout tag: 'master'")
endif()

execute_process(
  COMMAND "/usr/bin/git" submodule init
  WORKING_DIRECTORY "/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp"
  RESULT_VARIABLE error_code
  )
if(error_code)
  message(FATAL_ERROR "Failed to init submodules in: '/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp'")
endif()

execute_process(
  COMMAND "/usr/bin/git" submodule update --recursive 
  WORKING_DIRECTORY "/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp"
  RESULT_VARIABLE error_code
  )
if(error_code)
  message(FATAL_ERROR "Failed to update submodules in: '/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp'")
endif()

# Complete success, update the script-last-run stamp file:
#
execute_process(
  COMMAND ${CMAKE_COMMAND} -E copy
    "/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp-stamp/deps.jsoncpp-gitinfo.txt"
    "/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp-stamp/deps.jsoncpp-gitclone-lastrun.txt"
  WORKING_DIRECTORY "/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp"
  RESULT_VARIABLE error_code
  )
if(error_code)
  message(FATAL_ERROR "Failed to copy script-last-run stamp file: '/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp-stamp/deps.jsoncpp-gitclone-lastrun.txt'")
endif()

