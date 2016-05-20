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

//! Parse command line arguments
class opt_parse {
private:
    //! An argumetn representation to be passed to the program
    struct opt_param {
        //! is it required
        bool required;
        //! Does he needs arguments
        bool need_arg;
        //! What is its name
        std::string name;
        //! Description
        std::string desc;
        //! the argument of the arguments !
        std::string argument;
        bool is_passed = false;

        opt_param(bool required, bool need_arg, std::string name, std::string desc)
        : required(required), need_arg(need_arg), name(name), desc(desc) {}
    };

    //! Map of the different possibles argument as string and their
    //! rertpresntation
    std::unordered_map<std::string, opt_param> params;
    //! The total count of arguments for this program
    int argc;
    //! The list of arguments as a String
    char** argv;
    //! The total count of params
    int params_cnt = 0;

public:
    //! An iterator to iterate over all the arguments of
    //! the program
    class iterator {
    private:
        //! The arguments vector and their argumetns
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
        bool operator==(const iterator& rhs) const { return this->opt_pos == rhs.opt_pos; }
        bool operator!=(const iterator& rhs) const { return this->opt_pos != rhs.opt_pos; }
    };

    opt_parse() = delete;

    //! \param argc Total count of arguements
    //! \param argv Cmdline arguments from program startup
    opt_parse(int argc, char* argv[]);

    ~opt_parse();

    //! Add a argument required for your program
    //!
    //! If the specified argument is not given in cmdline, a error will be
    //! generated
    //! \param name The name of the parameter as a string (e.g "-l")
    //! \param desc A description that will be used by the function `print_help`
    //! \param need_arg Does the argument require a parameter
    void required(const std::string& name, const std::string& desc = "", bool need_arg = true);

    //! Add an optional argument for your program
    //!
    //! If the specified argument is not given in cmdline, a error will be
    //! generated
    //! \param name The name of the parameter as a string (e.g "-l")
    //! \param desc A description that will be used by the function `print_help`
    //! \param need_arg Does the argument require a parameter
    void optional(const std::string& name, const std::string& desc = "", bool need_arg = true);

    //! Process the parsing of the arguments
    bool execute();

    //! \return an iterator on the begin of the arguments
    iterator begin() const;

    //! \return the iterator on the end of the arguments
    iterator end() const;

    //! Print the usage using the parameter setted when referencing the arguments
    //! for the program
    void print_usage() const;

    //! Print an help message generated using all the specified arguments
    void print_help() const;

    //! Is there on the parameters (missing parameter ? unknows ? missing
    //! arguments ?)
    //! \return true if there is error, false otherwise
    bool has_error() const;

    //! Does the option exist or not ?
    //! \param opt The name of the option to check
    //! \return true if the param exist, false otherwise
    bool exist(const std::string& opt) const;

    //! Acces to an argument from its name
    //! \param opt The name of the option to check
    //! \return the the argument of the param as a string
    std::string operator[](const std::string& opt) const;
};

} // tool

} // etix
