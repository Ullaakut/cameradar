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

#include <describe.h>

namespace etix {
namespace cameradar {

//! Sends a request to the camera using the OPTION method,
//! then a DESCRIBE to check for valid IDs
//! then another DESCIBE with IDs if an authentication is needed
bool
curl_describe(const std::string& path, bool logs) {
    CURL* csession;
    CURLcode res;
    struct curl_slist* custom_msg = NULL;
    char URL[256];
    long rc;
    FILE* protofile = NULL;
    protofile = fopen("/dev/null", "wb");
    csession = curl_easy_init();
    if (csession == NULL) return -1;
    sprintf(URL, "%s", path.c_str());
    // These are the options for all following cURL requests
    // Activate verbose if debug is needed
    curl_easy_setopt(csession, CURLOPT_TIMEOUT, 1);
    curl_easy_setopt(csession, CURLOPT_NOBODY, 1);
    curl_easy_setopt(csession, CURLOPT_URL, URL);
    curl_easy_setopt(csession, CURLOPT_RTSP_STREAM_URI, URL);
    curl_easy_setopt(csession, CURLOPT_FOLLOWLOCATION, 0);
    curl_easy_setopt(csession, CURLOPT_HEADER, 0);
    curl_easy_setopt(csession, CURLOPT_INTERLEAVEDATA, protofile);
    curl_easy_setopt(csession, CURLOPT_VERBOSE, 0);
    curl_easy_setopt(csession, CURLOPT_RTSP_REQUEST, CURL_RTSPREQ_OPTIONS);
    curl_easy_setopt(csession, CURLOPT_WRITEDATA, protofile);
    // This request will handshake the stream's server, it should always return 200 OK
    curl_easy_perform(csession);
    curl_easy_getinfo(csession, CURLINFO_RESPONSE_CODE, &rc);
    custom_msg = curl_slist_append(
        custom_msg, "Accept: application/x-rtsp-mh, application/rtsl, application/sdp");
    curl_easy_setopt(csession, CURLOPT_RTSPHEADER, custom_msg);
    curl_easy_setopt(csession, CURLOPT_RTSP_REQUEST, CURL_RTSPREQ_DESCRIBE);
    curl_easy_setopt(csession, CURLOPT_WRITEDATA, protofile);
    // This request will check if the given path is right without the need of encrypted ids
    curl_easy_perform(
        csession); // will return 404 if no ids and bad route, 401 if ids, 200 is all ok
    res = curl_easy_getinfo(csession, CURLINFO_RESPONSE_CODE, &rc);
    unsigned long pos = path.find("@");
    if (pos != std::string::npos) {
        std::string encoded = etix::tool::encode::encode64(path.substr(7, pos - 7));
        custom_msg =
            curl_slist_append(custom_msg, std::string("Authorization: Basic " + encoded).c_str());
        curl_easy_setopt(csession, CURLOPT_RTSPHEADER, custom_msg);
        curl_easy_setopt(csession, CURLOPT_RTSP_REQUEST, CURL_RTSPREQ_DESCRIBE);
        curl_easy_setopt(csession, CURLOPT_WRITEDATA, protofile);
        // This request will check if the given ids are good
        curl_easy_perform(csession); // will return 404 if good ids, 401 if bad ids
        res = curl_easy_getinfo(csession, CURLINFO_RESPONSE_CODE, &rc);
    }
    curl_easy_cleanup(csession);
    fclose(protofile);
    curl_global_cleanup();
    if (logs) {
        if (rc != 401 && pos == std::string::npos)
            LOG_INFO_("Unprotected camera discovered.", "brutelogs");
        return ((res == CURLE_OK) && rc != 401);
    }
    return ((res == CURLE_OK) && rc != 404);
}
}
}
