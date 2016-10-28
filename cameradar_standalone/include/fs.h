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

#include <pwd.h>
#include <sys/types.h>
#include <unistd.h>
#include <fstream>
#include <string>

namespace etix {

namespace tool {

namespace fs {

enum class fs_error { is_dir, is_not_dir, dont_exist };

fs_error is_folder(const std::string& folder);
bool get_or_create_folder(const std::string& folder);
bool create_folder(const std::string& folder);
bool create_recursive_folder(const std::string& folder);
std::string home();

// this functions take a copy because we need to make some operations on the string
// for example, we need to apply std::string::pop_back
std::string get_file_folder(std::string full_file_path);

bool copy(const std::string& src, const std::string& dst);

} // fs

} // tool

} // etix
