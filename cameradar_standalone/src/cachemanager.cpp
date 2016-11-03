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

#include "cachemanager.h" // for cache_manager
#include <algorithm>      // for move
#include <dlfcn.h>        // for dlerror, dlclose, dlopen, dlsym, etc
#include <logger.h>       // for LOG_ERR_
#include <stdbool.h>      // for bool, false, true

#include <errno.h>

namespace etix {

namespace cameradar {

#ifdef __APPLE__
const std::string cache_manager::PLUGIN_EXT = ".dylib";
#elif __linux__
const std::string cache_manager::PLUGIN_EXT = ".so";
#endif

const std::string cache_manager::default_symbol = "cache_manager_instance_new";

cache_manager::cache_manager(const std::string& path,
                             const std::string& name,
                             const std::string& symbol)
: name(name), path(path), symbol(symbol), handle(nullptr), ptr(nullptr) {}

cache_manager::cache_manager(cache_manager&& old)
: path(std::move(old.path)), symbol(std::move(old.symbol)) {
    this->handle = old.handle;
    old.handle = nullptr;
    this->ptr = old.ptr;
    old.ptr = nullptr;
}

cache_manager::~cache_manager() {
    delete this->ptr;
    if (this->handle) { dlclose(handle); }
}

bool
cache_manager::make_instance() {
    cache_manager_iface* (*new_fn)() = nullptr;

    // Gets the path to the dynamic library
    auto real_path = this->make_full_path();

    // Opens it to get the handle
    this->handle = dlopen(real_path.c_str(), RTLD_LAZY);
    if (this->handle == nullptr) {
        std::cout << "error: " << dlerror() << std::endl;
        LOG_ERR_("Failed to load cache manager: " + this->name + ", invalid path",
                 "cache manager loader");
        return false;
    } else {
        // Gets the symbol and checks if the library is valid
        *(void**)(&new_fn) = dlsym(this->handle, symbol.c_str());
        if (dlerror() != nullptr) {
            LOG_ERR_("Invalid cache manager package: " + this->name, "cache manager loader");
            return false;
        }
    }

    // Returns a string containing the most recent dl* error
    dlerror();

    // Instantiates the cache manager
    this->ptr = (*new_fn)();
    if (this->ptr == nullptr) {
        LOG_ERR_("Invalid cache manager format: " + this->name, "cache manager loader");
        return false;
    }

    return true;
}

// Generates a path as such : /libdumb_cache_manager.so
std::string
cache_manager::make_full_path() {
    std::string full_path = this->path;
    full_path += "/lib";
    full_path += this->name;
    full_path += "_cache_manager";
    full_path += PLUGIN_EXT;

    return full_path;
}

cache_manager_iface* cache_manager::operator->() { return this->ptr; }

const cache_manager_iface* cache_manager::operator->() const { return this->ptr; }

bool
operator==(std::nullptr_t nullp, const cache_manager& p) {
    return p.ptr == nullp;
}

bool
operator==(const cache_manager& p, std::nullptr_t nullp) {
    return p.ptr == nullp;
}

bool
operator!=(std::nullptr_t nullp, const cache_manager& p) {
    return p.ptr != nullp;
}

bool
operator!=(const cache_manager& p, std::nullptr_t nullp) {
    return p.ptr != nullp;
}

cache_manager_base&
cache_manager_base::get_instance() {
    return *this;
}

} // cameradar

} // etix
