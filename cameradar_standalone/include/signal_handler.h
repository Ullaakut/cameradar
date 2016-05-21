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

#include <assert.h> // assert
#include <csignal>  // sigint
#include <iostream> // stc::cout

// To avoid an unused warning for the asserted in handle_signal
#define _unused(x) ((void)(x))

namespace etix {

namespace cameradar {

enum class stop_priority { running, stop, force_stop };

class event_handler {
public:
    event_handler(void) : ss(stop_priority::running) {}

    virtual int
    handle_signal(int signum) {
        assert(signum == SIGINT);
        _unused(signum);
        std::cout << "\b\b\b\033[K";
        if (this->ss == stop_priority::running)
            this->ss = stop_priority::stop;
        else
            this->ss = stop_priority::force_stop;
        return 0;
    }

    etix::cameradar::stop_priority
    should_stop(void) const {
        return this->ss;
    }

private:
    stop_priority ss;
};

class signal_handler {
private:
    signal_handler(void);
    signal_handler(const signal_handler&);
    signal_handler& operator=(const signal_handler&);

    static void call_handler(int signum);

    static event_handler handler;

public:
    static signal_handler& instance(void);
    etix::cameradar::stop_priority should_stop(void) const;
};
}
}
