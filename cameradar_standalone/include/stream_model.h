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

#pragma once

#include <string>
#include <json/value.h>

namespace etix {
namespace cameradar {

struct stream_model {
    // Ex : "172.16.100.113"
    std::string address;
    // Ex : 8554
    unsigned short port;
    // Ex : "admin"
    std::string username = "";
    // Ex : "123456"
    std::string password = "";
    // Ex : "live.sdp"
    std::string route = "";

    // Ex : "rtsp"
    std::string service_name;
    // Ex : "Vivotek HDCam"
    std::string product;
    // Ex : "RTSP"
    std::string protocol;
    // Ex : "Open"
    std::string state;

    // Ex : "true"
    bool path_found = false;
    // Ex : "true"
    bool ids_found = false;

    // Ex : "/thumbnails/cameradar"
    std::string thumbnail_path = "";
};
Json::Value deserialize(const stream_model& model);
}
}
