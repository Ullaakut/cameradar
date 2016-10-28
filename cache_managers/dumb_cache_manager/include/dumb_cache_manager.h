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

#include <cachemanager.h>
#include <configuration.h>
#include <logger.h>
#include <stream_model.h>
#include <vector>

namespace etix {
namespace cameradar {

class dumb_cache_manager : public cache_manager_base {
private:
    static const std::string name;
    std::vector<etix::cameradar::stream_model> streams;
    std::shared_ptr<etix::cameradar::configuration> configuration;

    std::mutex m;

public:
    using cache_manager_base::cache_manager_base;
    ~dumb_cache_manager();

    const std::string& get_name() const override;
    static const std::string& static_get_name();
    bool load_dumb_conf(std::shared_ptr<etix::cameradar::configuration> configuration);
    bool configure(std::shared_ptr<etix::cameradar::configuration> configuration) override;

    bool has_changed(const etix::cameradar::stream_model&);

    void set_streams(std::vector<etix::cameradar::stream_model> model);

    void update_stream(const etix::cameradar::stream_model& newmodel);

    std::vector<etix::cameradar::stream_model> get_streams();

    std::vector<etix::cameradar::stream_model> get_valid_streams();
};
}
}
