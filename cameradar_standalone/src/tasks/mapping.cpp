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

#include <tasks/mapping.h>

namespace etix {
namespace cameradar {

//! The first command checks if dpkg finds nmap in the system by cutting the
//! result and grepping
//! nmap from it.
//!
//! The second command checks the version of nmap, right now it needs to be the
//! 6.47 but this could
//! be changed to 6 or greater depending on the needs. In a docker container
//! this should not be a
//! problem.
bool
nmap_is_ok() {
    return (
        launch_command("test `dpkg -l | cut -c 5-9 | grep nmap` = nmap")
        // && launch_command("test `nmap --version | cut -c 14-18  | head -n2 | tail -n1` = 6.47")
        &&
        launch_command(
            "mkdir -p scans")); // Creates the directory in which the scans will be stored
}

//! Launches and checks the return of the nmap command
//! Uses the subnets specified in the conf file to launch nmap
bool
mapping::run() const {
    if (nmap_is_ok()) {
        std::string subnets = this->conf.subnets;
        std::replace(subnets.begin(), subnets.end(), ',', ' ');
        LOG_INFO_("Nmap 6.0 or greater found", "mapping");
        LOG_INFO_("Beginning mapping task. This may take a while.", "mapping");
        std::string cmd =
            "nmap -T4 -A " + subnets + " -p " + this->conf.ports + " -oX " + nmap_output;
        LOG_DEBUG_("Launching nmap : " + cmd, "mapping");
        bool ret = launch_command(cmd);
        if (ret)
            LOG_INFO_("Nmap XML output successfully generated in file: " + nmap_output, "mapping");
        else
            LOG_ERR_("Nmap command failed", "mapping");
        return ret;
    } else {
        LOG_ERR_("Nmap 6.0 or greater is required to launch Cameradar", "mapping");
        return false;
    }
}
}
}
