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

#include <stream_model.h>

namespace etix {
namespace cameradar {

Json::Value
deserialize(const stream_model& model) {
    Json::Value ret;

    ret["address"] = model.address;
    ret["port"] = model.port;
    ret["username"] = model.username;
    ret["password"] = model.password;
    ret["route"] = model.route;
    ret["service_name"] = model.service_name;
    ret["product"] = model.product;
    ret["protocol"] = model.protocol;
    ret["state"] = model.state;
    ret["path_found"] = model.path_found;
    ret["ids_found"] = model.ids_found;
    ret["thumbnail_path"] = model.thumbnail_path;
    return ret;
}
}
}
