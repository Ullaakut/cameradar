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

#include <fstream>         // std::ifstream
#include <unistd.h>        // access, F_OK
#include <configuration.h> // configuration

namespace etix {

namespace cameradar {

const std::string configuration::name_ = "configuration";

// read a file at the path "path"
// if the file is available we return the whole content as an std::string inside
// a pair
// otherwise return false and an empty string inside a pair
std::pair<bool, std::string>
read_file(const std::string& path) {
    auto line = std::string{};
    auto content = std::string{};
    auto file = std::ifstream{ path };

    if (file.is_open()) {
        while (getline(file, line)) { content += line + "\n"; }
        file.close();
    } else {
        return std::make_pair(false, std::string{});
    }

    return std::make_pair(true, content);
}

// Loads the IDS dictionary
bool
configuration::load_ids() {
    std::string content;

    LOG_DEBUG_("Trying to open ids file from " + this->rtsp_ids_file, "configuration");
    if (this->rtsp_ids_file.size()) {
        content = read_file(this->rtsp_ids_file.c_str()).second;
    } else {
        LOG_WARN_(
            "No ids file detected in your configuration, Cameradar will use "
            "the default one "
            "instead.",
            "configuration");
        content = read_file(default_ids_file_path_).second;
    }
    if (content.size()) {
        auto root = Json::Value();
        auto reader = Json::Reader();
        reader.parse(content, root);

        for (unsigned int i = 0; i < root["username"].size(); i++) {
            if (not root["username"][i].isString()) {
                LOG_ERR_("\"username\" should be of type string", "configuration");
                return false;
            }
            this->usernames.push_back(root["username"][i].asString());
        }
        for (unsigned int i = 0; i < root["password"].size(); i++) {
            if (not root["password"][i].isString()) {
                LOG_ERR_("\"password\" should be of type string", "configuration");
                return false;
            }
            this->passwords.push_back(root["password"][i].asString());
        }
        return true;
    } else {
        LOG_ERR_(
            "Could not load ids file. Make sure you provided a valid path in your "
            "configuration file.",
            "configuration");
        return false;
    }
}

// Loads the URL dictionary
bool
configuration::load_url() {
    std::string content;

    LOG_DEBUG_("Trying to open ids file from " + this->rtsp_ids_file, "configuration");
    if (this->rtsp_url_file.size()) {
        content = read_file(this->rtsp_url_file.c_str()).second;
    } else {
        LOG_WARN_(
            "No ids file detected in your configuration, Cameradar will use "
            "the default one "
            "instead.",
            "configuration");
        content = read_file(default_urls_file_path_).second;
    }
    if (content.size()) {
        auto root = Json::Value();
        auto reader = Json::Reader();
        reader.parse(content, root);
        //    auto result = tool::json::check_fields(
        //        {{"urls", Json::arrayValue, root["urls"]}}, "general
        //        configuration");

        //    if (not result.first) {
        //      LOG_ERR_(result.second, "general configuration");
        //      return false;
        //    }

        for (unsigned int i = 0; i < root["urls"].size(); i++) {
            if (not root["urls"][i].isString()) {
                LOG_ERR_("\"urls\" should be of type string", "configuration");
                return false;
            }
            this->paths.push_back(root["urls"][i].asString());
        }

        return true;
    } else {
        LOG_ERR_(
            "Could not load ids file. Make sure you provided a valid path in your "
            "configuration file.",
            "configuration");
        return false;
    }
}

std::pair<bool, configuration>
serialize(const Json::Value& root) {
    std::pair<bool, configuration> ret;

    try {
        ret.second.ports = root["ports"].asString();
        ret.second.subnets = root["subnets"].asString();
        ret.second.rtsp_ids_file = root["rtsp_ids_file"].asString();
        ret.second.rtsp_url_file = root["rtsp_url_file"].asString();
        ret.second.thumbnail_storage_path = root["thumbnail_storage_path"].asString();
        ret.second.cache_manager_path = root["cache_manager_path"].asString();
        ret.second.cache_manager_name = root["cache_manager_name"].asString();
        ret.first = true;
    } catch (std::exception& e) {
        LOG_ERR_("Configuration failed : " + std::string(e.what()), "configuration");
        ret.first = false;
    }
    return ret;
}

Json::Value
configuration::get_raw() const {
    return this->raw_conf;
}

// Loads the configuration from a path
// Returns a pair containing a boolean value & the configuration.
// Will return true & valid configuration if success
// Otherwise false & empty configuration
std::pair<bool, configuration>
load(const std::string& path) {
    // Check if the file exists at the given path
    if (access(path.c_str(), F_OK) == -1) {
        LOG_ERR_("Can't access: " + path, "configuration");
        return std::make_pair(false, configuration{});
    }

    // Get the content of the file
    auto content = read_file(path);
    if (not content.first) {
        LOG_ERR_(
            "Can't open configuration file, you should check your rights to "
            "access the file",
            "configuration");
        return std::make_pair(false, configuration{});
    }

    // Parse & validate the json
    auto root = Json::Value();

    auto reader = Json::Reader();
    auto parse_succes = reader.parse(content.second, root);
    if (not parse_succes) {
        LOG_ERR_("Can't load configuration, invalid json format:\n" +
                     reader.getFormattedErrorMessages(),
                 "configuration");
        return std::make_pair(false, configuration{});
    }
    // Deserialize the json to a configuration struct
    // and return
    // REPLACE THIS WITH JSONCPP
    std::pair<bool, configuration> conf = serialize(root);
    conf.second.raw_conf = root;
    conf.first &= conf.second.load_url();
    conf.first &= conf.second.load_ids();

    return conf;
}
}
}
