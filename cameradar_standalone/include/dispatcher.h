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

#include <list>             // sig
#include <memory>           // std::shared_ptr
#include <opt_parse.h>      // parsing opt
#include <logger.h>         // LOG
#include <configuration.h>  // conf
#include <thread>           // std::thread
#include <chrono>           // operator""ms
#include <signal_handler.h> // sig

// All the tasks managed by the dispatcher
#include <tasks/mapping.h>
#include <tasks/parsing.h>
#include <tasks/creds_attack.h>
#include <tasks/path_attack.h>
#include <tasks/thumbnail.h>
#include <tasks/stream_check.h>
#include <tasks/print.h>

namespace etix {
namespace cameradar {

enum class task {
    init,
    preparation,
    mapping,
    parsing,
    path_attack,
    creds_attack,
    thumb_generation,
    print,
    finished
};

class dispatcher {
private:
    bool busy;
    task current;
    std::string nmap_output;
    const configuration& conf;
    std::shared_ptr<etix::cameradar::cache_manager> cache;
    const std::pair<bool, etix::tool::opt_parse>& opts;
    std::list<cameradar_task*> queue;

public:
    dispatcher() = delete;
    dispatcher(const configuration& conf,
               std::shared_ptr<etix::cameradar::cache_manager> cache,
               const std::pair<bool, etix::tool::opt_parse>& opts)
    : busy(false)
    , current(task::init)
    , nmap_output("/tmp/scans/scan" + std::to_string(std::chrono::system_clock::to_time_t(
                               std::chrono::system_clock::now())) +
                  ".xml")
    , conf(conf)
    , cache(cache)
    , opts(opts){};
    ~dispatcher() = default;
    bool
    doing_stuff() const {
        return this->busy;
    }

    void do_stuff();

    void run();
};
}
}
