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

namespace etix {

namespace cameradar {

namespace mysql {

enum class execute_result { success, not_found, no_row_updated, sql_error, error };

//! Wrapper of a DB query result
//! Templated on the data type we want to return (list<model>, bool, whatever)
template <typename DataType>
struct query_result {
    DataType data;
    execute_result state;
    std::string error_msg;

    inline bool
    success(void) const {
        return state == execute_result::success;
    }
    inline bool
    error(void) const {
        return not success();
    }
};

//! Empty query result for when we just want to return the status
//! of the request with no associated data
template <>
struct query_result<void> {
    execute_result state;
    std::string error_msg;

    inline bool
    success(void) const {
        return state == execute_result::success;
    }
    inline bool
    error(void) const {
        return not success();
    }
};
typedef query_result<void> empty_result;

} //! mysql

} //! cameradar

} //! etix
