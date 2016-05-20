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

#include <tasks/thumbnail.h>

namespace etix {
namespace cameradar {

std::string
remove_trailing_backslash(std::string s) {
    while (s.back() == '/') { s.pop_back(); }
    return s;
}

// Tranforms the path into a path for the thumbnail
// Example :
// rtsp://username:password@172.16.100.13/live.sdp
// will become /storage/path/172.16.100.13/1345425533.jpg
std::string
thumbnail::build_output_file_path(const std::string& path) const {
    auto ss = std::stringstream{};

    ss << remove_trailing_backslash(this->conf.thumbnail_storage_path);
    ss << "/";
    ss << path;
    ss << "/";
    ss << std::to_string(std::chrono::system_clock::to_time_t(std::chrono::system_clock::now()));
    ss << ".jpg";

    return ss.str();
}

//! Gets all the discovered streams with good routes and logs
//! And launches an ffmpeg command to generate a thumbnail
//! In order to check for the stream validity
bool
thumbnail::run() const {
    std::vector<stream_model> streams = (*cache)->get_valid_streams();
    LOG_INFO_("Started thumbnail generation, it may take a while", "thumbnail");
    for (const auto& stream : streams) {
        if (signal_handler::instance().should_stop() != etix::cameradar::stop_priority::running)
            break;
        std::string ffmpeg_cmd =
            "mkdir -p %s ; "
            "ffmpeg "
            "-y "
            "-nostdin "
            "-loglevel quiet "
            "-i '%s' "
            "-vcodec mjpeg "
            "-vframes 1 "
            "-an "
            "-f image2 "
            "-s 320x240 "
            "'%s'";
        std::string fullpath = make_path(stream);
        std::string output = build_output_file_path(stream.address);
        ffmpeg_cmd = tool::fmt(ffmpeg_cmd.c_str(),
                               output.substr(0, output.find_last_of("/")).c_str(),
                               fullpath.c_str(),
                               output.c_str());
        if (!launch_command(ffmpeg_cmd)) {
            LOG_WARN_("The following command [" + ffmpeg_cmd +
                          "] didn't work. That can either mean that the stream is "
                          "not valid or "
                          "that there is a problem with the camera.",
                      "thumbnail_generation");
        } else {
            LOG_DEBUG_("Generated thumbnail : " + ffmpeg_cmd, "thumbnail_generation");
            try {
                stream_model result{ stream.address,    stream.port,      stream.username,
                                     stream.password,   stream.route,     stream.service_name,
                                     stream.product,    stream.protocol,  stream.state,
                                     stream.path_found, stream.ids_found, output };
                (*cache)->update_stream(result);

            } catch (std::exception& e) { LOG_DEBUG_(e.what(), "thumbnail_generation"); }
        }
    }
    LOG_INFO_("All thumbnails have been successfully generated in " +
                  this->conf.thumbnail_storage_path,
              "thumbnail_generation");
    return true;
}
}
}
