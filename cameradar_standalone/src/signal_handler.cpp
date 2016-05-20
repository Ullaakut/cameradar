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

#include <signal_handler.h>

namespace etix {
namespace cameradar {

event_handler signal_handler::handler;

signal_handler::signal_handler() {}

void
signal_handler::call_handler(int signum) {
    handler.handle_signal(signum);
}

signal_handler&
signal_handler::instance(void) {
    static signal_handler singleton;

    struct sigaction sa;
    sa.sa_handler = call_handler;
    sigemptyset(&sa.sa_mask);
    sa.sa_flags = 0;
    sigaction(SIGINT, &sa, 0);

    return singleton;
}

etix::cameradar::stop_priority
signal_handler::should_stop(void) const {
    return handler.should_stop();
}
}
}
