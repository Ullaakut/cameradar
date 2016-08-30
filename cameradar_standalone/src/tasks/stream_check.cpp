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

#include <tasks/stream_check.h>

namespace etix {
namespace cameradar {

//! Gets all the discovered streams with good routes and logs
//! And launches an ffmpeg command to generate a thumbnail
//! In order to check for the stream validity
bool
stream_check::run() const {
    GstElement* pipeline;
    GstElement* elem;

    gst_init(nullptr, nullptr);

    std::vector<stream_model> streams = (*cache)->get_valid_streams();

    if (not streams.size()) {
      LOG_WARN_("There were no valid streams to check. Cameradar will stop.", "stream_check");
      return false;
    }
    for (const auto& stream : streams) {
        GError* error = NULL;

        pipeline =
            gst_parse_launch("rtspsrc name=source ! rtph264depay ! h264parse ! fakesink", &error);

        std::string location = "rtsp://";
        location += stream.username + ":" + stream.password + "@" + stream.address + ":" + std::to_string(stream.port);
        if (pipeline == NULL) {
            LOG_ERR_("[" + stream.address + "] Can't configure pipeline", "stream_check");
            return false;
        } else {
            elem = gst_bin_get_by_name(GST_BIN(pipeline), "source");
            LOG_DEBUG_("Launching gstreamer check on rtsp://" + stream.username + ":" + stream.password + "@" + stream.address + ":" + std::to_string(stream.port), "gstreamer check");
            g_object_set(G_OBJECT(elem), "location", location.c_str(), "latency", 20, NULL);

            if (gst_element_set_state(pipeline, GST_STATE_PLAYING) == GST_STATE_CHANGE_FAILURE) {
                LOG_ERR_(
                    "This stream is unaccessible with GStreamer, there must be a "
                    "configuration issue",
                    "stream_check");
                gst_object_unref(pipeline);
                stream_model invalidstream{
                    stream.address,   stream.port,         stream.username,  stream.password,
                    stream.route,     stream.service_name, stream.product,   stream.protocol,
                    "invalid stream", stream.path_found,   stream.ids_found, stream.thumbnail_path
                };
                (*cache)->update_stream(invalidstream);
                return false;
            }
            LOG_INFO_("[" + stream.address + "] Set pipeline to playing", "stream_check");
        }
    }
    LOG_INFO_("All streams could be accessed with GStreamer", "stream_check");
    return true;
}
}
}
