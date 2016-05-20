// Copyright (C) 2016 Etix Labs - All Rights Reserved.
// All information contained herein is, and remains the property of Etix Labs
// and its suppliers,
// if any. The intellectual and technical concepts contained herein are
// proprietary to Etix Labs
// Dissemination of this information or reproduction of this material is
// strictly forbidden unless
// prior written permission is obtained from Etix Labs.

#include <tasks/brutepath.h>

namespace etix {
namespace cameradar {

static const std::string no_route_found_ =
    "The url.json files' default paths didn't match with the discovered "
    "cameras. Either "
    "they have a custom path, or your url.json file does not contain enough "
    "default "
    "routes. Thumbnail generation is impossible without the path.";

//! Tries to match the detected combination of Username / Password
//! with a route for the camera stream. Creates a resource in the DB upon
//! valid discovery
bool
brutepath::test_path(const stream_model& stream, const std::string& route) const {
    bool found = false;
    std::string path = stream.service_name + "://" + stream.username + ":" + stream.password + "@" +
                       stream.address + ":" + std::to_string(stream.port);
    if (route.front() != '/') { path += "/"; }
    path += route;
    LOG_DEBUG_("Testing path : " + path, "brutepath");
    try {
        if (curl_describe(path, false)) {
            // insert in DB and go to the next port, print a cool message
            found = true;
            LOG_INFO_("Discovered a valid path : [" + path + "]", "brutepath");
            stream_model newstream{
                stream.address,      stream.port,          stream.username, stream.password, route,
                stream.service_name, stream.product,       stream.protocol, stream.state,    true,
                stream.ids_found,    stream.thumbnail_path
            };
            (*cache)->update_stream(newstream);
        } else {
            stream_model newstream{
                stream.address,      stream.port,          stream.username, stream.password, route,
                stream.service_name, stream.product,       stream.protocol, stream.state,    false,
                stream.ids_found,    stream.thumbnail_path
            };
            (*cache)->update_stream(newstream);
        }
    } catch (const std::runtime_error& e) { LOG_INFO_(e.what(), "brutepath"); }
    return found;
}

bool
path_already_found(std::vector<stream_model> streams, stream_model model) {
    for (const auto& stream : streams) {
        if ((model.address == stream.address) && (model.port == stream.port) && stream.path_found)
            return true;
    }
    return false;
}

//! Tries to discover a route on all RTSP streams in DB
//! Uses the url.json file to try different routes
bool
brutepath::run() const {
    LOG_INFO_("Beginning bruteforce of the camera paths task, it may take a while.", "bruteforce");
    std::vector<stream_model> streams = (*cache)->get_streams();
    int found = 0;
    for (const auto& stream : streams) {
        if (signal_handler::instance().should_stop() != etix::cameradar::stop_priority::running)
            break;
        if (path_already_found(streams, stream)) {
            LOG_INFO_(stream.address +
                          " : This camera's path was already discovered in the database."
                          "Skipping to the next camera.",
                      "brutepath");
            ++found;
        } else {
            for (const auto& route : conf.paths) {
                if (signal_handler::instance().should_stop() !=
                    etix::cameradar::stop_priority::running)
                    break;
                if (test_path(stream, route)) {
                    found++;
                    break;
                }
            }
        }
    }
    if (!found) {
        LOG_WARN_(no_route_found_, "brutepath");

    } else
        LOG_INFO_("Found " + std::to_string(found) + " routes for " +
                      std::to_string(streams.size()) + " cameras",
                  "brutepath");
    return true;
}
}
}
