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

#include "spdlog/spdlog.h"
#include <sstream>
#include <string>

namespace etix {

namespace tool {

enum class loglevel { DEBUG = 1, INFO = 2, WARN = 4, ERR = 5, CRITICAL = 6 };

inline std::string
format_output(const std::string& from, const std::string& message) {
    auto ss = std::stringstream{};

    ss << "(" << from << "): ";
    ss << message;

    return ss.str();
}

class logger {
    std::string name;
    std::shared_ptr<spdlog::logger> console;

    logger(const std::string& plugin)
    : name(plugin), console(spdlog::stdout_logger_mt("cameradar")) {}

public:
    static logger&
    get_instance(const std::string& name = "") {
        static logger self(name);
        return self;
    }
    void
    set_level(loglevel level) {
        switch (level) {
        case loglevel::DEBUG: this->console->set_level(spdlog::level::level_enum::debug); break;

        case loglevel::INFO: this->console->set_level(spdlog::level::level_enum::info); break;

        case loglevel::WARN: this->console->set_level(spdlog::level::level_enum::warn); break;

        case loglevel::ERR: this->console->set_level(spdlog::level::level_enum::err); break;

        case loglevel::CRITICAL:
            this->console->set_level(spdlog::level::level_enum::critical);
            break;
        }
    }

    static void
    info(const std::string& message) {
        etix::tool::logger::get_instance().console->info(message);
    }

    static void
    warn(const std::string& message) {
        etix::tool::logger::get_instance().console->warn(message);
    }

    static void
    err(const std::string& message) {
        etix::tool::logger::get_instance().console->error(message);
    }

    static void
    debug(const std::string& message) {
        etix::tool::logger::get_instance().console->debug(message);
    }
};
}
}

// Should be replaced to calls to spdlog::logger::getlogger(const std::string&
// name)
#define LOG_WARN_(message, from)                                                                   \
    etix::tool::logger::get_instance().warn(etix::tool::format_output(                             \
        std::string(from) + "::" + __FUNCTION__ + ":" + std::to_string(__LINE__), message))
#define LOG_ERR_(message, from)                                                                    \
    etix::tool::logger::get_instance().err(etix::tool::format_output(                              \
        std::string(from) + "::" + __FUNCTION__ + ":" + std::to_string(__LINE__), message))
#define LOG_DEBUG_(message, from)                                                                  \
    etix::tool::logger::get_instance().debug(etix::tool::format_output(                            \
        std::string(from) + "::" + __FUNCTION__ + ":" + std::to_string(__LINE__), message))
#define LOG_INFO_(message, from)                                                                   \
    etix::tool::logger::get_instance().info(etix::tool::format_output(                             \
        std::string(from) + "::" + __FUNCTION__ + ":" + std::to_string(__LINE__), message))
