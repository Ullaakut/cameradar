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

#include "opt_parse.h"
#include <iostream>

namespace etix {

namespace tool {

opt_parse::opt_parse(int argc, char* argv[]) : argc(argc), argv(argv) {}

opt_parse::~opt_parse() {}

void
opt_parse::required(const std::string& name, const std::string& desc, bool need_arg) {
    this->params.emplace(name, opt_param(true, need_arg, name, desc));
}

void
opt_parse::optional(const std::string& name, const std::string& desc, bool need_arg) {
    this->params.emplace(name, opt_param(false, need_arg, name, desc));
}

bool
opt_parse::execute() {
    int i = 1;

    // if params are invalid
    if (this->argc < 1 || not this->argv) { return false; }

    while (i != this->argc) {
        // there is less argument than argc
        if (not this->argv[i]) { return false; }
        auto params = this->params.find(std::string(this->argv[i]));
        if (params != this->params.end()) {
            this->params_cnt += 1;
            (*params).second.is_passed = true;
            if ((*params).second.need_arg == true && (i + 1) != this->argc) {
                (*params).second.argument = this->argv[i + 1];
                i += 1;
            }
        }
        i += 1;
    }
    return true;
}

opt_parse::iterator
opt_parse::begin() const {
    std::vector<std::pair<std::string, std::string>> p;

    for (auto entry : this->params) {
        p.push_back(std::make_pair(entry.second.name, entry.second.argument));
    }
    return iterator(p, 0);
}

opt_parse::iterator
opt_parse::end() const {
    return iterator(std::vector<std::pair<std::string, std::string>>(), this->params_cnt);
}

void
opt_parse::print_usage() const {
    std::cout << "Usage: " << this->argv[0];

    for (auto entry : this->params) {
        if (entry.second.required == true) {
            if (entry.second.need_arg == true) { std::cout << " <arg>"; }
        }
    }
    std::cout << std::endl;
}

void
opt_parse::print_help() const {
    std::cout << "help: " << this->argv[0] << std::endl;

    for (auto entry : this->params) {
        std::cout << entry.second.name << "    " << entry.second.desc << std::endl;
    }
}

bool
opt_parse::has_error() const {
    for (auto entry : this->params) {
        // is the parameter required ?
        // the parameter need arguement ?
        if ((entry.second.required == true && entry.second.is_passed == false) ||
            (entry.second.is_passed == true && entry.second.need_arg == true &&
             entry.second.argument == "")) {
            return true;
        }
    }

    return false;
}

bool
opt_parse::exist(const std::string& opt) const {
    auto params = this->params.find(opt);

    if (params == this->params.end()) { return false; }

    return (*params).second.is_passed;
}

std::string opt_parse::operator[](const std::string& opt) const {
    std::string param("");

    auto opt_param = this->params.find(opt);
    if (opt_param != this->params.end()) { param = (*opt_param).second.argument; }

    return param;
}

} // tool

} // etix
