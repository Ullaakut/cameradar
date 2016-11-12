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

#include <logger.h>
#include <rtsp_path.h>

namespace etix {

namespace cameradar {

const std::string
make_path(const stream_model& model) {
    if (model.password != "" || model.username != "") {
        std::string ret(model.service_name + "://" + model.username + ":" + model.password + "@" +
                        model.address + ":" + std::to_string(model.port) + model.route);
        LOG_DEBUG_(ret, "debug");
        return ret;
    } else {
        std::string ret(model.service_name + "://" + model.address + ":" +
                        std::to_string(model.port) + model.route);
        LOG_DEBUG_(ret, "debug");
        return ret;
    }
}
}
}
