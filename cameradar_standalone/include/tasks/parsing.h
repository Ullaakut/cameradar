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

#include <cameradar_task.h> // task interface
#include <tinyxml.h>        // parsing
#include <stream_model.h>   // data model
#include <cachemanager.h>   // cacheManager

namespace etix {
namespace cameradar {

class parsing : public etix::cameradar::cameradar_task {
    const configuration& conf;
    std::shared_ptr<cache_manager> cache;
    std::string nmap_output;

public:
    parsing() = delete;
    parsing(std::shared_ptr<cache_manager> cache,
            const configuration& conf,
            std::string nmap_output)
    : conf(conf), cache(cache), nmap_output(nmap_output) {}
    parsing(const parsing& ref) = delete;

    virtual bool run() const;

    void parse_camera(TiXmlElement*, std::vector<etix::cameradar::stream_model>& data) const;

    bool print_detected_cameras(const std::vector<etix::cameradar::stream_model>& data) const;
};
}
}
