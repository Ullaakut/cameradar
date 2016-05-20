

set(command "/usr/bin/make")
execute_process(
  COMMAND ${command}
  RESULT_VARIABLE result
  OUTPUT_FILE "/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp-stamp/deps.jsoncpp-build-out.log"
  ERROR_FILE "/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp-stamp/deps.jsoncpp-build-err.log"
  )
if(result)
  set(msg "Command failed: ${result}\n")
  foreach(arg IN LISTS command)
    set(msg "${msg} '${arg}'")
  endforeach()
  set(msg "${msg}\nSee also\n  /home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp-stamp/deps.jsoncpp-build-*.log")
  message(FATAL_ERROR "${msg}")
else()
  set(msg "deps.jsoncpp build command succeeded.  See also /home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp-stamp/deps.jsoncpp-build-*.log")
  message(STATUS "${msg}")
endif()
