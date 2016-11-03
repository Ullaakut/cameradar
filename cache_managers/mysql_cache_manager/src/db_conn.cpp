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

#include "db_conn.h"            // for db_connection
#include "cppconn/connection.h" // for Connection
#include "query_result.h"       // for queries
#include <cppconn/driver.h>     // for get_driver_instance, etc
#include <cppconn/exception.h>  // for SQLException
#include <cppconn/statement.h>  // for Statement
#include <fmt.h>                // for fmt
#include <logger.h>             // for LOG_

namespace etix {

namespace cameradar {

namespace mysql {

const std::string db_connection::create_database_query = "CREATE DATABASE IF NOT EXISTS %s";

db_connection::db_connection() : connected(false) {}

db_connection::~db_connection() { delete this->connection; }

std::pair<bool, std::string>
db_connection::connect(const std::string& host,
                       const std::string& user,
                       const std::string& pass,
                       const std::string& db_name,
                       bool create_db_if_not_exist) {
    this->db_name = db_name;

    try {
        this->driver = get_driver_instance();
        if (this->driver == nullptr) {
            return std::make_pair(false, "Cannot instantiate sql_driver");
        }
        this->connection = driver->connect(host, user, pass);
        if (this->connection == nullptr) return std::make_pair(false, "Cannot connect to mysql");

        this->connected = true;
        if (create_db_if_not_exist) {
            auto cdb = this->create_database();
            if (cdb.state == mysql::execute_result::sql_error) { return { false, cdb.error_msg }; }
            this->connection->setSchema(db_name);
        }
    } catch (sql::SQLException& e) {
        this->connected = false;
        return { false, e.what() };
    }

    return std::make_pair(true, "");
}

empty_result
db_connection::execute(const std::string& request) {
    std::lock_guard<std::mutex> lock(this->access_mtx);

    sql::Statement* stmt = nullptr;
    empty_result return_value = { execute_result::success, "" };

    if (!this->is_connected()) {
        return { execute_result::sql_error, "Error, not connected to MySQL database" };
    }

    try {
        stmt = this->connection->createStatement();
        stmt->execute(request);
        if (stmt->getUpdateCount() == 0) {
            return_value = { execute_result::no_row_updated, "No row updated" };
        }
    } catch (sql::SQLException& e) { return_value = { execute_result::sql_error, e.what() }; }
    delete stmt;

    return return_value;
}

query_result<sql::ResultSet*>
db_connection::query(const std::string& query) {
    std::lock_guard<std::mutex> lock(this->access_mtx);

    sql::Statement* stmt = nullptr;
    query_result<sql::ResultSet*> return_value = { nullptr, execute_result::success, "" };

    if (!this->is_connected()) {
        return { nullptr, execute_result::sql_error, "Error, not connected to MySQL database" };
    }

    try {
        stmt = this->connection->createStatement();
        return_value = { stmt->executeQuery(query), execute_result::success, "" };
    } catch (sql::SQLException& e) {
        return_value = { nullptr, execute_result::sql_error, e.what() };
    }
    delete stmt;

    return return_value;
}

bool
db_connection::is_connected() {
    if (this->connection == nullptr) return false;

    // check if our connection is always valid
    if (this->connection->isClosed() || not this->connection->isValid()) {
        LOG_INFO_("MySQL database connection is either closed or invalid, try to reconnect.",
                  "db_connection");
        this->connection->reconnect();
        if (this->connection->isClosed() || not this->connection->isValid()) {
            this->connected = false;
            LOG_ERR_("Unable to reconnect to MySQL.", "db_connection");
        }
    }
    return this->connected;
}

empty_result
db_connection::create_database() {
    auto query = tool::fmt(this->create_database_query, this->db_name.c_str());
    return this->execute(query);
}

} // mysql

} // cameradar

} // etix
