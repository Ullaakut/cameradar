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

#include <tasks/print.h>

namespace etix {
namespace cameradar {

// Launches and checks the return of the nmap command
// Uses the subnets specified in the conf file to launch nmap
bool
print::run() const {
    std::ofstream file;
    std::vector<stream_model> results = (*cache)->get_valid_streams();
    bool first = true;

    file.open(default_result_path);
    if (file.fail()) {
        LOG_ERR_("Could not open file " + default_result_path, "print");
        return false;
    }

    file << "[\n";
    for (const auto& stream : results) {
        file << deserialize(stream).toStyledString();
        if (first)
            first = false;
        else
            file << ",";
        LOG_INFO_("Generated JSON Result : " + deserialize(stream).toStyledString(), "print");
    }
    file << "]";
    file.close();
    return true;
}
}
}
