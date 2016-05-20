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
#include <cachemanager.h>   // cacheManager
#include <describe.h>       // send DESCRIBE through cURL
#include <signal_handler.h> // signals
#include <stream_model.h>   // data model

namespace etix {
namespace cameradar {

class brutelogs : public etix::cameradar::cameradar_task {
    const configuration& conf;
    std::shared_ptr<cache_manager> cache;
    std::string nmap_output;

public:
    brutelogs() = delete;
    brutelogs(std::shared_ptr<cache_manager> cache,
              const configuration& conf,
              std::string nmap_output)
    : conf(conf), cache(cache), nmap_output(nmap_output) {}
    brutelogs(const brutelogs& ref) = delete;

    virtual bool run() const;

    bool test_ids(const etix::cameradar::stream_model& cit,
                  const std::string& pit,
                  const std::string& uit) const;
};
}
}
