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

#include <string>        // for string
#include <unordered_map> // for unordered_map
#include <utility>       // for pair
#include <vector>        // for vector

namespace etix {

namespace tool {

class opt_parse {
private:
    struct opt_param {
        bool required;
        bool need_arg;
        std::string name;
        std::string desc;
        std::string argument;
        bool is_passed = false;

        opt_param(bool required, bool need_arg, std::string name, std::string desc)
        : required(required), need_arg(need_arg), name(name), desc(desc) {}
    };

    std::unordered_map<std::string, opt_param> params;
    int argc;
    char** argv;
    int params_cnt = 0;

public:
    class iterator {
    private:
        std::vector<std::pair<std::string, std::string>> args;
        unsigned int opt_pos = 0;

    public:
        iterator(std::vector<std::pair<std::string, std::string>> args, unsigned int opt_pos)
        : args(args), opt_pos(opt_pos) {}
        iterator operator++() {
            this->opt_pos += 1;
            return *this;
        }
        std::pair<std::string, std::string>& operator*() { return this->args.at(this->opt_pos); }
        bool
        operator==(const iterator& rhs) const {
            return this->opt_pos == rhs.opt_pos;
        }
        bool
        operator!=(const iterator& rhs) const {
            return this->opt_pos != rhs.opt_pos;
        }
    };

    opt_parse() = delete;

    opt_parse(int argc, char* argv[]);

    ~opt_parse();

    void required(const std::string& name, const std::string& desc = "", bool need_arg = true);

    void optional(const std::string& name, const std::string& desc = "", bool need_arg = true);

    bool execute();

    iterator begin() const;

    iterator end() const;

    void print_usage() const;

    void print_help() const;

    bool has_error() const;

    bool exist(const std::string& opt) const;

    std::string operator[](const std::string& opt) const;
};

} // tool

} // etix
