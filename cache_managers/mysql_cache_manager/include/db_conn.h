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

#include <cppconn/resultset.h> // for ResultSet
#include <mutex>               // for mutex
#include <stdbool.h>           // for bool, false
#include <string>              // for string
#include <utility>             // for pair, make_pair

#include "query_result.h"

namespace sql {
class Connection;
class Driver;
class ResultSet;
}

namespace etix {

namespace cameradar {

namespace mysql {

//! MySQL Database connection handling
//! Abstracts all connection to the database
class db_connection {
private:
    static const std::string create_database_query;

    //! SQL driver
    sql::Driver* driver = nullptr;
    //! SQL connection
    sql::Connection* connection = nullptr;
    std::mutex access_mtx;
    bool connected = false;

    std::string db_name;

    //! Create the database if it doesn't exist at connector launch
    empty_result create_database(void);

public:
    db_connection(void);
    ~db_connection(void);

    //! Try to connect to the database
    std::pair<bool, std::string> connect(const std::string& host,
                                         const std::string& user,
                                         const std::string& pass,
                                         const std::string& db_name,
                                         bool create_db_if_not_exist = true);

    //! Execute a MySQL command
    empty_result execute(const std::string& request);

    //! Execute a query
    query_result<sql::ResultSet*> query(const std::string& query);

    bool is_connected();

    //! Return db_name
    const std::string&
    get_db_name(void) const {
        return this->db_name;
    }
};

} // mysql

} // cameradar

} // etix
