

set(command "/usr/bin/cmake;-P;/home/ullaakut/Work/cameradar/deps/jsoncpp/tmp/deps.jsoncpp-gitclone.cmake")
execute_process(
  COMMAND ${command}
  RESULT_VARIABLE result
  OUTPUT_FILE "/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp-stamp/deps.jsoncpp-download-out.log"
  ERROR_FILE "/home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp-stamp/deps.jsoncpp-download-err.log"
  )
if(result)
  set(msg "Command failed: ${result}\n")
  foreach(arg IN LISTS command)
    set(msg "${msg} '${arg}'")
  endforeach()
  set(msg "${msg}\nSee also\n  /home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp-stamp/deps.jsoncpp-download-*.log")
  message(FATAL_ERROR "${msg}")
else()
  set(msg "deps.jsoncpp download command succeeded.  See also /home/ullaakut/Work/cameradar/deps/jsoncpp/src/deps.jsoncpp-stamp/deps.jsoncpp-download-*.log")
  message(STATUS "${msg}")
endif()
