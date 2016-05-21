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
#include <db_conn.h>
#include <fmt.h>
#include <logger.h>
#include <stream_model.h>
#include <vector>

namespace etix {

namespace cameradar {

struct mysql_configuration {
    unsigned int port;
    std::string host;
    std::string db_name;
    std::string user;
    std::string password;

    mysql_configuration() = default;

    mysql_configuration(unsigned int port,
                        const std::string& host,
                        const std::string& db_name,
                        const std::string& user = "",
                        const std::string& password = "")
    : port(port), host(host), db_name(db_name), user(user), password(password) {}
};

class mysql_cache_manager : public cache_manager_base {
private:
    static const std::string name;
    std::vector<etix::cameradar::stream_model> streams;
    std::shared_ptr<etix::cameradar::configuration> configuration;
    etix::cameradar::mysql_configuration db_conf;
    etix::cameradar::mysql::db_connection connection;

    static const std::string create_table_query;
    static const std::string insert_with_id_query;
    static const std::string exist_query;
    static const std::string get_results_query;
    static const std::string update_result_query;

public:
    using cache_manager_base::cache_manager_base;
    ~mysql_cache_manager();

    // Specific to MySQL
    bool execute_query(const std::string& query);

    const std::string& get_name() const override;
    static const std::string& static_get_name();
    bool load_mysql_conf(std::shared_ptr<etix::cameradar::configuration> configuration);
    bool configure(std::shared_ptr<etix::cameradar::configuration> configuration) override;

    void set_streams(std::vector<etix::cameradar::stream_model> model);

    void update_stream(const etix::cameradar::stream_model& newmodel);

    std::vector<etix::cameradar::stream_model> get_streams();

    std::vector<etix::cameradar::stream_model> get_valid_streams();
};
}
}
