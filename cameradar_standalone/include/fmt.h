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

#include <string>
#include <iostream>
#include <mutex>

namespace etix {

namespace tool {

static std::mutex mutex;

//! Format a string with the given arguments
//! same behavior as sprintf.
template <class... Args>
std::string
fmt(const std::string& base, Args... args) {
    std::lock_guard<std::mutex> guard(mutex);
    static char buf[512];

    std::sprintf(buf, base.c_str(), args...);

    return std::string(buf);
}

} // tool

} // etix
