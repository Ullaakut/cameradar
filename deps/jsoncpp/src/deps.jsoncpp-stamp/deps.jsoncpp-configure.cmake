

set(command "/usr/bin/cmake;-DCMAKE_INSTALL_PREFIX=/home/ullaakut/Work/cameradar/deps/jsoncpp;-DBUILD_TYPE=Release;-DBUILD_STATIC_LIBS=OFF;-DBUILD_SHARED_LIBS=ON;-DJSONCPP_WITH_TESTS=OFF;-DJSONCPP_WITH_POST_BUILD_UNITTEST=OFF;/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp")
execute_process(
  COMMAND ${command}
  RESULT_VARIABLE result
  OUTPUT_FILE "/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp-stamp/deps.jsoncpp-configure-out.log"
  ERROR_FILE "/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp-stamp/deps.jsoncpp-configure-err.log"
  )
if(result)
  set(msg "Command failed: ${result}\n")
  foreach(arg IN LISTS command)
    set(msg "${msg} '${arg}'")
  endforeach()
  set(msg "${msg}\nSee also\n  /home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp-stamp/deps.jsoncpp-configure-*.log")
  message(FATAL_ERROR "${msg}")
else()
  set(msg "deps.jsoncpp configure command succeeded.  See also /home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp-stamp/deps.jsoncpp-configure-*.log")
  message(STATUS "${msg}")
endif()
