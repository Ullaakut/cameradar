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

#include <tasks/parsing.h>

namespace etix {
namespace cameradar {

static const std::string no_hosts_found_ =
    "No hosts were discovered on your network. Please check your internet "
    "connexion "
    "and verify that the subnetworks you specified in your configuration file "
    "were "
    "accessible";

// Avoids segfaults on unknown xml structure
std::string
xml_safe_get(const TiXmlElement* elem, const std::string& attr) {
    if (elem == nullptr) return "closed";
    if (elem->Attribute(attr.c_str()) != nullptr) return std::string(elem->Attribute(attr.c_str()));
    return "closed";
}

// Parse a single host node (generally containing only one camera)
// Pushes it back to the data structure
void
parsing::parse_camera(TiXmlElement* xml_host, std::vector<stream_model>& data) const {
    TiXmlElement* xml_streams = xml_host->FirstChild("ports")->ToElement();
    stream_model stream;
    for (TiXmlElement* xml_stream = xml_streams->FirstChild("port")->ToElement(); xml_stream;
         xml_stream = xml_stream->NextSiblingElement("port")) {
        stream.address = xml_safe_get(xml_host->FirstChild("address")->ToElement(), "addr");
        stream.protocol = xml_safe_get(xml_stream, "protocol");
        stream.port = static_cast<unsigned short>(std::stoi(xml_safe_get(xml_stream, "portid")));
        TiXmlElement* state = xml_stream->FirstChild("state")->ToElement();
        stream.state = xml_safe_get(state, "state");
        TiXmlElement* service;
        if (state->NextSibling("service") &&
            (service = state->NextSibling("service")->ToElement())) {
            stream.service_name = xml_safe_get(service, "name");
            stream.product = xml_safe_get(service, "product");
        } else {
            stream.service_name = "closed";
            stream.product = "closed";
        }
        data.push_back(stream);
    }
}

// Prints all detected cameras into the data structure and stops the program if
// no open RTSP streams were found
bool
parsing::print_detected_cameras(const std::vector<stream_model>& data) const {
    int added = 0;
    for (const auto& stream : data) {
        if (!stream.service_name.compare("rtsp") && !stream.state.compare("open")) {
            try {
                LOG_INFO_("Generated JSON Result : " + deserialize(stream).toStyledString(),
                          "print");
                added++;
            } catch (const std::runtime_error& e) {
                LOG_WARN_("Port already scanned : " + std::string(e.what()), "parsing");
                added++;
            }
        }
    }
    if (!added) {
        LOG_WARN_(
            "Mapping unsuccessful, no rtsp streams were discovered. You "
            "should try other "
            "subnetworks",
            "parsing");
        return false;
    }
    LOG_INFO_("Mapping successfuly ended, " + std::to_string(added) +
                  " RTSP streams were discovered.",
              "parsing");
    (*cache)->set_streams(data);
    return true;
}

// Opens the nmap output file, parses the data of each discovered port
// Adds the RTSP ports only into the DB
bool
parsing::run() const {
    std::vector<stream_model> data;
    try {
        TiXmlDocument doc(nmap_output.c_str());
        doc.LoadFile();
        TiXmlHandle docHandle(&doc);

        TiXmlElement* nmaprun = docHandle.FirstChild("nmaprun").ToElement();
        TiXmlNode* xml_node = nmaprun->FirstChild("host");
        if (xml_node == NULL) return false;
        TiXmlElement* xml_host;
        if ((xml_host = xml_node->ToElement()) && xml_host->Attribute("endtime"))
            for (xml_host = xml_node->ToElement(); xml_host;
                 xml_host = xml_host->NextSiblingElement("host")) {
                parse_camera(xml_host, data);
            }
        else
            LOG_WARN_(no_hosts_found_, "parsing");
        if (data.size() == 0) { LOG_WARN_("No cameras were discovered", "parsing"); }
        return print_detected_cameras(data);
    } catch (const std::exception& e) {
        LOG_ERR_("Error during parsing. brutepath aborted : " + std::string(e.what()), "parsing");
        return false;
    }
}
}
}
