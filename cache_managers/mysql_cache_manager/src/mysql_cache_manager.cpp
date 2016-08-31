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

#include <mysql_cache_manager.h>

/* DATA FORMAT
**
**    Example :
**
**   "address" : "173.16.100.45",
**   "ids_found" : true,
**   "password" : "123456",
**   "path_found" : true,
**   "port" : 554,
**   "product" : "Vivotek FD9381-HTV",
**   "protocol" : "tcp",
**   "route" : "/live.sdp",
**   "service_name" : "rtsp",
**   "state" : "open",
**   "thumbnail_path" : "/tmp/127.0.0.1/1463735257.jpg",
**   "username" : "admin"
**
*/

namespace etix {

namespace cameradar {

const std::string mysql_cache_manager::create_table_query =
    "CREATE TABLE IF NOT EXISTS `results` ("
    "`id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT, "
    "`address` tinytext NOT NULL, "
    "`password` tinytext NOT NULL, "
    "`product` tinytext NOT NULL, "
    "`protocol` tinytext NOT NULL, "
    "`route` tinytext NOT NULL, "
    "`service_name` tinytext NOT NULL, "
    "`state` tinytext NOT NULL, "
    "`thumbnail_path` tinytext NOT NULL, "
    "`username` tinytext NOT NULL, "
    "`port` int(11) UNSIGNED NOT NULL, "
    "`ids_found` tinytext NOT NULL, "
    "`path_found` tinytext NOT NULL, "
    "PRIMARY KEY (`id`));";

const std::string mysql_cache_manager::insert_with_id_query =
    "INSERT INTO `%s`.`results`"
    " (`address`, `password`, `product`, `protocol`, `route`, `service_name`, `state`, "
    "`thumbnail_path`, `username`, `port`, `ids_found`, `path_found`)"
    " VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')";

const std::string mysql_cache_manager::update_result_query =
    "UPDATE `%s`.`results` SET"
    " `results`.`address` = '%s',"
    " `results`.`password` = '%s',"
    " `results`.`product` = '%s',"
    " `results`.`protocol` = '%s',"
    " `results`.`route` = '%s',"
    " `results`.`service_name` = '%s',"
    " `results`.`state` = '%s',"
    " `results`.`thumbnail_path` = '%s',"
    " `results`.`username` = '%s',"
    " `results`.`port` = '%s',"
    " `results`.`ids_found` = '%s',"
    " `results`.`path_found` = '%s'"
    " WHERE `results`.`address` LIKE '%s'";

const std::string mysql_cache_manager::exist_query =
    "SELECT * FROM `%s`.`results` WHERE `results`.`address` = '%s'";

const std::string mysql_cache_manager::get_results_query = "SELECT * FROM `%s`.`results`";

const std::string mysql_cache_manager::name = "mysql-cache-manager";

mysql_cache_manager::~mysql_cache_manager() {}

const std::string&
mysql_cache_manager::get_name() const {
    return mysql_cache_manager::static_get_name();
}

const std::string&
mysql_cache_manager::static_get_name() {
    return mysql_cache_manager::name;
}

bool
mysql_cache_manager::configure(std::shared_ptr<etix::cameradar::configuration> configuration) {
    return this->load_mysql_conf(configuration);
}

bool
mysql_cache_manager::execute_query(const std::string& query) {
    auto check_err = [](const auto& res) {
        if (res.state == mysql::execute_result::sql_error) {
            LOG_WARN_(res.error_msg, "mysql_cache_manager");
            return false;
        }
        return true;
    };
    return check_err(this->connection.execute(query));
}

bool
mysql_cache_manager::load_mysql_conf(
    std::shared_ptr<etix::cameradar::configuration> configuration) {
    this->configuration = configuration;

    try {
        this->db_conf.host = configuration->raw_conf["mysql_db"]["host"].asString();
        this->db_conf.port = configuration->raw_conf["mysql_db"]["port"].asUInt();
        this->db_conf.user = configuration->raw_conf["mysql_db"]["user"].asString();
        this->db_conf.password = configuration->raw_conf["mysql_db"]["password"].asString();
        this->db_conf.db_name = configuration->raw_conf["mysql_db"]["db_name"].asString();
    } catch (std::exception& e) {
        LOG_ERR_("Configuration of the MySQL db failed : " + std::string(e.what()),
                 "mysql_cache_manager");
        return false;
    }

    if (not this->connection
                .connect(db_conf.host + ":" + std::to_string(db_conf.port),
                         db_conf.user,
                         db_conf.password,
                         db_conf.db_name)
                .first) {
        LOG_ERR_("Configuration of the MySQL DB failed", "mysql_cache_manager");
        return false;
    }

    // Tries to create the Result table in the DB and returns the success state
    return (execute_query(create_table_query));
}

//! Replaces all cached streams by the content of the vector given as
//! parameter
void
mysql_cache_manager::set_streams(std::vector<etix::cameradar::stream_model> models) {
    LOG_DEBUG_("Beginning stream list DB insertion", "mysql_cache_manager");
    for (const auto& model : models) {
        if (!model.service_name.compare("rtsp") && !model.state.compare("open")) {
          auto query = tool::fmt(
              this->exist_query, this->connection.get_db_name().c_str(), model.address.c_str());
          auto result = this->connection.query(query);
          // If an entry already exists for this address in the database,
          // no need to insert it.

          // TODO : Update an entry if it already exists.

          if (result.data->next()) continue;

          query = tool::fmt(this->insert_with_id_query,
                            this->connection.get_db_name().c_str(),
                            model.address.c_str(),
                            model.password.c_str(),
                            model.product.c_str(),
                            model.protocol.c_str(),
                            model.route.c_str(),
                            model.service_name.c_str(),
                            model.state.c_str(),
                            model.thumbnail_path.c_str(),
                            model.username.c_str(),
                            std::to_string(model.port).c_str(),
                            std::to_string(model.ids_found).c_str(),
                            std::to_string(model.path_found).c_str());
          execute_query(query);
      }
    }
}

//! Inserts a single stream to the cache
void
mysql_cache_manager::update_stream(const etix::cameradar::stream_model& model) {
    auto query = tool::fmt(this->update_result_query,
                           this->connection.get_db_name().c_str(),
                           model.address.c_str(),
                           model.password.c_str(),
                           model.product.c_str(),
                           model.protocol.c_str(),
                           model.route.c_str(),
                           model.service_name.c_str(),
                           model.state.c_str(),
                           model.thumbnail_path.c_str(),
                           model.username.c_str(),
                           std::to_string(model.port).c_str(),
                           std::to_string(model.ids_found).c_str(),
                           std::to_string(model.path_found).c_str(),
                           model.address.c_str());
    execute_query(query);
}

//! Gets all cached streams
std::vector<etix::cameradar::stream_model>
mysql_cache_manager::get_streams() {
    auto query = tool::fmt(this->get_results_query, this->connection.get_db_name().c_str());
    auto result = this->connection.query(query);

    if (not result.data) {
        delete result.data;
        return {};
    }

    std::vector<stream_model> lst;
    while (result.data->next()) {
        // If it's an open RTSP stream
        if (not result.data->getString("state").compare("open") &&
            not result.data->getString("service_name").compare("rtsp")) {
            stream_model s{
                result.data->getString("address"),     result.data->getUInt("port"),
                result.data->getString("username"),    result.data->getString("password"),
                result.data->getString("route"),       result.data->getString("service_name"),
                result.data->getString("product"),     result.data->getString("protocol"),
                result.data->getString("state"),       result.data->getBoolean("ids_found"),
                result.data->getBoolean("path_found"), result.data->getString("thumbnail_path")
            };
            lst.push_back(s);
        }
    }

    delete result.data;
    return lst;
}

//! Gets all valid streams
std::vector<etix::cameradar::stream_model>
mysql_cache_manager::get_valid_streams() {
    auto query = tool::fmt(this->get_results_query, this->connection.get_db_name().c_str());
    auto result = this->connection.query(query);

    if (not result.data) {
        delete result.data;
        return {};
    }

    std::vector<stream_model> lst;
    while (result.data->next()) {
        // If the ID and the Path were found add this stream
        if (not result.data->getString("ids_found").compare("1") &&
            not result.data->getString("path_found").compare("1")) {
            stream_model s{
                result.data->getString("address"),     result.data->getUInt("port"),
                result.data->getString("username"),    result.data->getString("password"),
                result.data->getString("route"),       result.data->getString("service_name"),
                result.data->getString("product"),     result.data->getString("protocol"),
                result.data->getString("state"),       result.data->getBoolean("ids_found"),
                result.data->getBoolean("path_found"), result.data->getString("thumbnail_path")
            };
            lst.push_back(s);
        }
    }

    delete result.data;
    return lst;
}

extern "C" {
cache_manager_iface*
cache_manager_instance_new() {
    return new mysql_cache_manager();
}
}
}
}
