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

#include <fs.h>         // fs::home
#include <opt_parse.h>  // parsing opt
#include <dispatcher.h> // program loop

namespace cmrdr = etix::cameradar;

// Command line parsing
std::pair<bool, etix::tool::opt_parse>
parse_cmdline(int argc, char* argv[]) {
    auto opt_parse = etix::tool::opt_parse{ argc, argv };

    opt_parse.optional("-c", "Path to the configuration file (-c /path/to/conf)", true);
    opt_parse.optional("-l", "Set log level (-l 4 will only show warnings and errors)", true);
    opt_parse.optional("-d", "Launch the discovery tool on the given subnet", false);
    opt_parse.optional("-b", "Launch the bruteforce tool on all discovered devices", false);
    opt_parse.optional("-t", "Generate thumbnails from detected cameras", false);
    opt_parse.optional("-g", "Check if the stream can be opened with GStreamer", false);
    opt_parse.optional("-v", "Display Cameradar's version", false);
    opt_parse.optional("-h", "Display this help", false);
    opt_parse.execute();

    if (opt_parse.exist("-h")) {
        opt_parse.print_help();
        return std::make_pair(false, opt_parse);
    } else if (opt_parse.exist("-v")) {
        std::cout << "Cameradar 0.1" << std::endl;
        return std::make_pair(false, opt_parse);
    } else if (opt_parse.has_error()) {
        std::cout << "Usage: ./cameradar [option]\n\toptions:\n" << std::endl;
        opt_parse.print_help();
        return std::make_pair(false, opt_parse);
    }

    return std::make_pair(true, opt_parse);
}

// Check if a folder exists, is readable and writable
bool
check_folder(const std::string& path) {
    struct stat sb;

    if ((stat(path.c_str(), &sb) == 0) && (S_ISDIR(sb.st_mode)) && (sb.st_mode & S_IRUSR) &&
        (sb.st_mode & S_IWUSR)) {
        LOG_INFO_("Folder " + path + " is available and has sufficient rights", "main");
        return true;
    }
    LOG_ERR_("Folder " + path + " has insufficient rights, please check your configuration",
             "main");
    return false;
}

// Check if the storage path is available
bool
check_storage_path(const std::string& thumbnail_storage_path) {
    LOG_INFO_("Checking if storage path exists and are usable", "main");
    return (check_folder(thumbnail_storage_path));
}

int
main(int argc, char* argv[]) {
    etix::tool::logger::get_instance("cameradar");
    auto args = parse_cmdline(argc, argv);
    if (not args.first) return EXIT_FAILURE;

    // configure file configuration path
    auto conf_path = std::string{};
    if (not args.second.exist("-c")) {
        conf_path = etix::cameradar::default_configuration_path;
        LOG_WARN_("No custom path set, trying to use default path: " + conf_path, "main");
    } else {
        conf_path = args.second["-c"];
    }

    if (not args.second.exist("-l")) {
        etix::tool::logger::get_instance("cameradar").set_level(etix::tool::loglevel::INFO);
        LOG_INFO_("No log level set, using log level 2 (ignoring DEBUG)", "main");
    } else {
        try {
            int level = std::stoi(args.second["-l"]);
            etix::tool::logger::get_instance("cameradar")
                .set_level(static_cast<etix::tool::loglevel>(level));
        } catch (...) {
            LOG_ERR_("Invalid log level format, log level should be 1, 2, 4, 5 or 6", "main");
            return EXIT_FAILURE;
        }
    }

    // Try to load the configuration
    auto conf = cmrdr::load(conf_path);
    if (not conf.first) { return EXIT_FAILURE; }

    LOG_INFO_("Configuration successfully loaded", "main");

    // If one of the path is invalid, exit
    auto paths_ok = check_storage_path(conf.second.thumbnail_storage_path);
    if (not paths_ok) { return EXIT_FAILURE; }

    // Here we should get the cache manager but for now we will juste
    // make a dumb cache manager
    auto plug = std::make_shared<etix::cameradar::cache_manager>(conf.second.cache_manager_path,
                                                                 conf.second.cache_manager_name);

    if (not plug->make_instance()) {
        LOG_ERR_(std::string("Invalid cache manager "), "cameradar");
        return false;
    }

    LOG_INFO_("Launching Cameradar, press CTRL+C to gracefully stop", "main");

    etix::cameradar::dispatcher disp(conf.second, plug, args);

    disp.run();

    LOG_WARN_("See ya !", "cameradar");
    return EXIT_SUCCESS;
}
