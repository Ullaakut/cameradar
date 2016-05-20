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

#include "dispatcher.h"

namespace etix {
namespace cameradar {

// The main loop of the binary
void
dispatcher::run() {
    std::thread worker(&dispatcher::do_stuff, this);
    using namespace std::chrono_literals;
    // catch CTRL+C signal
    signal_handler::instance();

    // wait for event or end
    while (signal_handler::instance().should_stop() not_eq stop_priority::stop &&
           current != task::finished) {
        std::this_thread::sleep_for(30ms);
    }

    if (doing_stuff()) {
        LOG_INFO_("Waiting for a task to terminate", "dispatcher");
        LOG_INFO_("Press CTRL+C again to force stop", "dispatcher");
    }

    // waiting for task to cleanup / force stop command
    while ((signal_handler::instance().should_stop() not_eq stop_priority::force_stop) and
           doing_stuff()) {
        std::this_thread::sleep_for(std::chrono::milliseconds(30));
    }
    worker.join();
}

//! This loop is used to add all the tasks specified in the command line
//! And then run them successively
void
dispatcher::do_stuff() {
    if (opts.second.exist("-d")) {
        queue.push_back(new etix::cameradar::mapping(cache, conf, nmap_output));
        queue.push_back(new etix::cameradar::parsing(cache, conf, nmap_output));
    }
    if (opts.second.exist("-b")) {
        queue.push_back(new etix::cameradar::brutelogs(cache, conf, nmap_output));
        queue.push_back(new etix::cameradar::brutepath(cache, conf, nmap_output));
    }
    if (opts.second.exist("-t")) {
        queue.push_back(new etix::cameradar::thumbnail(cache, conf, nmap_output));
    }
    if (opts.second.exist("-g")) {
        queue.push_back(new etix::cameradar::stream_check(cache, conf, nmap_output));
    }
    if (!opts.second.exist("-d") && !opts.second.exist("-b") && !opts.second.exist("-t") &&
        !opts.second.exist("-g")) {
        queue.push_back(new etix::cameradar::mapping(cache, conf, nmap_output));
        queue.push_back(new etix::cameradar::parsing(cache, conf, nmap_output));
        queue.push_back(new etix::cameradar::brutelogs(cache, conf, nmap_output));
        queue.push_back(new etix::cameradar::brutepath(cache, conf, nmap_output));
        queue.push_back(new etix::cameradar::thumbnail(cache, conf, nmap_output));
        queue.push_back(new etix::cameradar::stream_check(cache, conf, nmap_output));
    }
    queue.push_back(new etix::cameradar::print(cache, conf, nmap_output));
    while (queue.size() > 0 && signal_handler::instance().should_stop() == stop_priority::running) {
        if (queue.front()->run())
            queue.pop_front();
        else {
            LOG_ERR_("An error occured in one of the tasks, Cameradar will now stop.", "dispatcher");
            break;
        }
    }
    this->current = task::finished;
}
}
}
