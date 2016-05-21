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

#include <configuration.h>
#include <memory>
#include <stream_model.h>
#include <vector>

namespace etix {
namespace cameradar {

//! The interface a cache_manager should implement to be valid
class cache_manager_iface {
public:
    virtual ~cache_manager_iface() {}

    //! Launches the manager configuration
    //! \return false if failed
    virtual bool configure(std::shared_ptr<etix::cameradar::configuration> configuration) = 0;

    //! get the name of the cache manager
    virtual const std::string& get_name() const = 0;

    //! Replaces all cached streams by the content of the vector given as
    //! parameter
    virtual void set_streams(std::vector<etix::cameradar::stream_model> model) = 0;

    //! Inserts a single stream to the cache
    virtual void update_stream(const etix::cameradar::stream_model& newmodel) = 0;

    //! Gets all cached streams
    virtual std::vector<etix::cameradar::stream_model> get_streams() = 0;

    //! Gets all valid streams which have been accessed
    virtual std::vector<etix::cameradar::stream_model> get_valid_streams() = 0;
};

class cache_manager_base : public cache_manager_iface {
public:
    cache_manager_base() = default;
    virtual ~cache_manager_base() = default;

    //! Launches the cache manager configuration
    //! \return false if failed
    virtual bool configure(std::shared_ptr<etix::cameradar::configuration> configuration) = 0;

    //! get the name of the cache manager
    virtual const std::string& get_name() const = 0;

    //! Replaces all cached streams by the content of the vector given as
    //! parameter
    virtual void set_streams(std::vector<etix::cameradar::stream_model> model) = 0;

    //! Updates a single stream to the cache
    virtual void update_stream(const etix::cameradar::stream_model& newmodel) = 0;

    //! Gets all cached streams
    virtual std::vector<etix::cameradar::stream_model> get_streams() = 0;

    //! Gets all valid streams which have been accessed
    virtual std::vector<etix::cameradar::stream_model> get_valid_streams() = 0;

    //! Get the manager's instance
    cache_manager_base& get_instance();

    template <typename I, typename T>
    std::shared_ptr<T>
    get() {
        static_assert(std::is_base_of<cache_manager_base, I>::value,
                      "I must implement cache_manager_base");
        std::shared_ptr<I> cache_manager(dynamic_cast<I*>(this));
        if (not cache_manager) return nullptr;
        return cache_manager->template get<T>();
    }
};

//! The representation of a cache manager
//!
//! This class loads a shared library, and tries to call an extern "C"
//! function which should instanciate a new instance of the plugin.
class cache_manager {
private:
    static const std::string PLUGIN_EXT;
    static const std::string default_symbol;

    //! the name of the cache manager
    std::string name;

    //! The path where the manager is located
    //! should be specified in the configuration file
    std::string path;

    //! The symbol entry point of the manager to
    //! call to create an instance from the shared library
    std::string symbol;

    //! The handle to the shared library where is stored the manager
    void* handle = nullptr;

    //! The cache manager instance if it is successfully loaded
    cache_manager_iface* ptr = nullptr;

    //! Internal function that creates the full path of the cache manager
    //!
    //! full path is composed of: the path, the name, the string "_cache-manager"
    //! and the extension PLUGIN_EXT depending of the platform
    std::string make_full_path();

public:
    //! Delete constructor
    cache_manager() = delete;

    //! The manager needs a path and a symbol to be instantiated.
    //! The symbol can be changed if the plugin entry point
    //! is different than the standard one.
    cache_manager(const std::string& path,
                  const std::string& name,
                  const std::string& symbol = default_symbol);

    //  //! Copy constructor
    //  cache_manager(cache_manager &other);

    //! Move constructor
    cache_manager(cache_manager&& old);

    ~cache_manager();

    //! Creates the instance of the cache_manager
    //!
    // \return false if the cache_manager failed to be instantiated or if
    // the cache_manager is not a valid cache manager, true otherwise
    bool make_instance();

    template <typename I, typename T>
    std::shared_ptr<T>
    get() {
        static_assert(std::is_base_of<cache_manager_base, I>::value,
                      "I must implement plugin_base");
        return this->get<I, T>();
    }

    //! Helper to access internal loaded cache_manager
    //!
    //! Gives access to the methods of the cache_manager using the operator
    //! -> (e.g.: cache_manager->get_name());
    cache_manager_iface* operator->();
    const cache_manager_iface* operator->() const;

    //! helper function to check if a cache_manager is instantiated or not
    friend bool operator==(std::nullptr_t nullp, const cache_manager& p);

    //! helper function to check if a cache_manager is instantiated or not
    friend bool operator==(const cache_manager& p, std::nullptr_t nullp);

    //! helper function to check if a cache_manager is instantiated or not
    friend bool operator!=(std::nullptr_t nullp, const cache_manager& p);

    //! helper function to check if a cache_manager is instantiated or not
    friend bool operator!=(const cache_manager& p, std::nullptr_t nullp);
};
}
}
