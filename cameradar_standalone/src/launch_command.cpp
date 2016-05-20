// Copyright 2016 Etix Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#include <launch_command.h>

namespace etix {
namespace cameradar {

//! Launches a command and checks for the return value
bool
launch_command(const std::string& cmd) {
    int status = system(cmd.c_str());
    if (status < 0) {
        LOG_ERR_("Error: " + std::string(strerror(errno)) + "", "dispatcher");
        return false;
    } else {
        if (WIFEXITED(status)) {
            LOG_DEBUG_("Program returned normally, exit code " +
                           std::to_string(WEXITSTATUS(status)),
                       "dispatcher");
            return true;
        } else
            LOG_WARN_("Program exited abnormaly.", "dispatcher");
    }
    return false;
}
}
}
