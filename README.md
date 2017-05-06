# Cameradar

## An RTSP surveillance camera access multitool

[![cameradar License](https://img.shields.io/badge/license-Apache-blue.svg?style=flat)](#license)
[![Docker Pulls](https://img.shields.io/docker/pulls/ullaakut/cameradar.svg?style=flat)](https://hub.docker.com/r/ullaakut/cameradar/)
[![Build](https://img.shields.io/travis/EtixLabs/cameradar/master.svg?style=flat)](https://travis-ci.org/EtixLabs/cameradar)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/6ab80cfa7069413e8e7d7e18320309e3)](https://www.codacy.com/app/brendan-le-glaunec/cameradar?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=EtixLabs/cameradar&amp;utm_campaign=Badge_Grade)
[![Latest release](https://img.shields.io/github/release/EtixLabs/cameradar.svg?style=flat)](https://github.com/EtixLabs/cameradar/releases/latest)


#### Cameradar allows you to:

* **Detect open RTSP hosts** on any accessible target
* Get their public info (hostname, port, camera model, etc.)
* Launch automated dictionary attacks to get their **stream route** (for example /live.sdp)
* Launch automated dictionary attacks to get the **username and password** of the cameras
* **Generate thumbnails** from them to check if the streams are valid and to have a quick preview of their content
* Try to create a Gstreamer pipeline to check if they are **properly encoded**
* Print a summary of all the informations Cameradar could get

#### And all of this in a _single command-line_.

Of course, you can also call for individual tasks if you plug in a Database to Cameradar using the MySQL cache manager for example. You can create your own cache manager by following the simple example of the **dumb cache manager**.

<p align="center"><img src="https://raw.githubusercontent.com/EtixLabs/cameradar/master/Cameradar.png" width="350"/></p>

## Table of content

- [Docker Image](#docker-image)
- [Quick install](#quick-install)
  - [Dependencies](#quick-install###dependencies)
  - [Five steps guide](#quick-install###five-steps-guide)
- [Manual installation](#manual-installation)
  - [Dependencies](#manual-installation###dependencies)
  - [Steps](#manual-installation###Steps)
- [Advanced docker deployment](#advanced-docker-deployment)
  - [Dependencies](#advanced-docker-deployment###dependencies)
  - [Deploy a custom version of Cameradar](#advanced-docker-deployment###deploy-a-custom-version-of-cameradar)
- [Configuration](#configuration)
- [Output](#output)
- [Check camera access](#check-camera-access)
- [Command line options](#command-line-options)
- [Next improvements](#next-improvements)
- [Contribution](#contribution)
- [Frequently Asked Questions](#frequently-asked-questions)
- [License](#license)

## Docker Image

This is the fastest and simplest way to use Cameradar. To do this you will just need `docker` on your machine.

Run

```
docker run  -v /tmp/thumbs/:/tmp/thumbs \
            -e CAMERAS_TARGET=your_target \
            ullaakut/cameradar:tag
```

* `your_target` can be a subnet (e.g.: `172.16.100.0/24`) or even an IP (e.g.: `172.16.100.10`), a range of IPs (e.g.: `172.16.100.10-172.16.100.20`) or a mix of all those separated by commas (e.g.: `172.17.100.0/24,172.16.100.10-172.16.100.20,0.0.0.0`).
* `tag` allows you to specify a specific version for camerada. If you don't specify any tag, you will use the latest version by default (recommended)

Check [Cameradar's readme on the Docker Hub](https://hub.docker.com/r/ullaakut/cameradar/) for more information and more command-line options.

The generated thumbnails will be in `/tmp/thumbs` on both your machine and the `cameradar` container.

For more complex use of the Docker image, see the `Environment variables` part of [Cameradar's readme on the Docker Hub](https://hub.docker.com/r/ullaakut/cameradar/).

## Quick install

The quick install uses docker to build Cameradar without polluting your machine with dependencies and makes it easy to deploy Cameradar in a few commands. **However, it may require networking knowledge, as your docker containers will need access to the cameras subnetwork.**

### Dependencies

The only dependencies are `docker`, `docker-tools`, `git` and `make`.

### Five steps guide

1. `git clone https://github.com/EtixLabs/cameradar.git`
2. `cd cameradar/deployment`
3. Tweak the `conf/cameradar.conf.json` as you need (see [the configuration guide here](#configuration) for more information)
4. `docker-compose build ; docker-compose up`

By default, the version of the package in the deployment should be the last stable release.

If you want to scan a different target or different ports, change the values `CAMERAS_TARGET` and `CAMERAS_PORTS` in the `docker-compose.yml` file.

The generated thumbnails will be in the `cameradar_thumbnails` folder after Cameradar has finished executing.

If you want to deploy your custom version of Cameradar using the same method, you should check the [advanced docker deployment](#advanced-docker-deployment) tutorial here.

## Manual installation

The manual installation is recommended if you want to tweak Cameradar and quickly test them using CMake and running Cameradar in command-line. If you just want to use Cameradar, it is recommended to use the [quick install](#quick-install) instead.

### Dependencies

To install Cameradar you will need these packages

* cmake (`cmake`)
* git (`git`)
* gstreamer1.x (`libgstreamer1.0-dev`)
* ffmpeg (`ffmpeg`)
* boost (`libboost-all-dev`)
* libcurl (`libcurl4-openssl-dev`)

### Steps

The simplest way would be to follow these steps :

1. `git clone https://github.com/EtixLabs/cameradar.git`
2. `cd cameradar`
3. `mkdir build`
4. `cd build`
5. `cmake ..`
6. `make`
7. `cd cameradar_standalone`
8. `./cameradar -s the_target_you_want_to_scan`

## Advanced Docker deployment

In case you want to use Docker to deploy your custom version of Cameradar.

### Dependencies

The only dependencies are `docker` and `docker-compose`.

### Using the package generation script
1. `git clone https://github.com/EtixLabs/cameradar.git`
2. `cd cameradar/deployment`
3. `rm *.tar.gz`
4. `./build_last_package.sh`
5. `docker-compose build cameradar`
6. `docker-compose up cameradar`

### Deploy a custom version of Cameradar by hand

1. `git clone https://github.com/EtixLabs/cameradar.git`
2. `cd cameradar`
3. `mkdir build`
4. `cd build`
5. `cmake .. -DCMAKE_BUILD_TYPE=Release`
6. `make package`
7. `cp cameradar_*_Release_Linux.tar.gz ../deployment`
8. `cd ../deployment`
9. `docker-compose build cameradar`
10. `docker-compose up cameradar`

### Configuration

Here is the basic content of the configuration file with simple placeholders :
```json
{
  "mysql_db" : {
     "host" : "MYSQL_SERVER_IP_ADDRESS",
     "port" : MYSQL_SERVER_PORT,
     "user": "root",
     "password": "root",
     "db_name": "cmrdr"
  },
  "target" : "target1,target2,target3,[...]",
  "ports" : "PORT1,PORT2,[...]",
  "rtsp_url_file" : "/path/to/url/dictionary",
  "rtsp_ids_file" : "/path/to/url/dictionary",
  "thumbnail_storage_path" : "/valid/path/to/a/storage/directory",
  "cache_manager_path" : "/path/to/cache/manager",
  "cache_manager_name" : "CACHE_MANAGER_NAME"
}
```

This **configuration is needed only if you want to overwrite the default values**, which are :

```json
{
  "target" : "localhost",
  "ports" : "554,8554",
  "rtsp_url_file" : "conf/url.json",
  "rtsp_ids_file" : "conf/ids.json",
  "thumbnail_storage_path" : "/tmp",
  "cache_manager_path" : "../cache_managers/dumb_cache_manager",
  "cache_manager_name" : "dumb"
}
```

This means that **by default Cameradar will not use a database**, will scan localhost and the ports 554 (default RTSP port) and 8554 (default emulated RTSP port), use the default constructor dictionaries and store the thumbnails in `/tmp`. If you need to override simply the target or ports, you can use the [command line options](#command-line-options).

The targets should be passed separated by commas only, and their target format should be the same as used in nmap.
```json
"target" : "172.100.16.0/24,172.100.17.0/24,localhost,192.168.1.13"
```

The **RTSP ports for most cameras are 554**, so you should probably specify 554 as one of the ports you scan. Not giving any ports in the configuration will scan every port of every host found on the target.

You **can use your own files for the ids and routes dictionaries** used to attack the cameras, but the Cameradar repository already gives you a good base that works with most cameras.

The thumbnail storage path should be a **valid and accessible directory** in which the thumbnails will be stored.

The cache manager path and name variables are used to change the cache manager you want to load into Cameradar. If you want to, you can code your own cache manager using a database, a file, a remote server, [...]. Feel free to share it by creating a merge request on this repository if you developed a generic manager (It must not be specific to your company's infrastructure).

## Output

For each camera, Cameradar will output these JSON objects :

```json
{
   "address" : "173.16.100.45",
   "ids_found" : true,
   "password" : "123456",
   "path_found" : true,
   "port" : 554,
   "product" : "Vivotek FD9381-HTV",
   "protocol" : "tcp",
   "route" : "/live.sdp",
   "service_name" : "rtsp",
   "state" : "open",
   "thumbnail_path" : "/tmp/127.0.0.1/1463735257.jpg",
   "username" : "admin"
}
```

## Check camera access

If you have [VLC Media Player](http://www.videolan.org/vlc/), you should be able to use the GUI to connect to the RTSP stream using this format : `rtsp://username:password@address:port/route`

With the above result, the RTSP URL would be `rtsp://admin:123456@173.16.100.45:554/live.sdp`

If you're still in your console however, you can go even faster by using **vlc in commmand-line** and just run `vlc rtsp://username:password@address:port/route` with the camera's info instead of the placeholders.

## Command line options

* **"-c"** : Set a custom path to the configuration file (-c /path/to/conf)
<<<<<<< HEAD
* **"-s"** : Set custom subnets (overrides configuration) : You can use this argument in many ways, using a subnet (e.g.: `172.16.100.0/24`) or even an IP (e.g.: `172.16.100.10`), a range of IPs (e.g.: `172.16.100.10-172.16.100.20`) or a mix of all those (e.g.: `172.17.100.0/24,172.16.100.10-172.16.100.20,0.0.0.0`).
=======
* **"-s"** : Set custom target (overrides configuration)
>>>>>>> 5489969... v2.0.0: Rename subnet to target to avoid confusion
* **"-p"** : Set custom ports (overrides configuration)
* **"-m"** : Set number of threads (*Default value : 1*)
* **"-l"** : Set log level
  * **"-l 1"** : Log level DEBUG
    * _Will print everything including debugging logs_
  * **"-l 2"** : Log level INFO
    * _Prints every normal information_
  * **"-l 4"** : Log level WARNING
    * _Only prints warning and errors_
  * **"-l 5"** : Log level ERROR
    * _Only prints errors_
  * **"-l 6"** : Log level CRITICAL
    * _Doesn't print anything since Cameradar can't have critical failures right now, however you can use this level to debug your own code easily or if you add new critical layers_
* **"-d"** : Launch the discovery tool
* **"-b"** : Launch the dictionary attack tool on all discovered devices
  * Needs either to be launched with the -d option or to use an advanced cache manager (DB, file, ...) with data already present
* **"-t"** : Generate thumbnails from detected cameras
  * Needs either to be launched with the -d option or to use an advanced cache manager (DB, file, ...) with data already present
* **"-g"** : Check if the stream can be opened with GStreamer
  * Needs either to be launched with the -d option or to use an advanced cache manager (DB, file, ...) with data already present
* **"-v"** : Display Cameradar's version
* **"-h"** : Display this help
* **"--gst-rtsp-server"** : Use this option if the attack does not seem to work (only detects the username but not the path, or the opposite). This option will switch the order of the attacks to prioritize path over credentials, which is the way priority is handled for cameras that use GStreamer's RTSP server.

## Contribution

See [the contribution document](/CONTRIBUTION.md) to get started.

## Frequently Asked Questions

> My camera's credentials are guessed by Cameradar but the RTSP URL is not!

Your camera probably uses GST RTSP Server internally. Try the `--gst-rtsp-server` command-line option, and if it does not work, send me the Cameradar output in DEBUG mode (`-l 1`) and I will help you.

> Cameradar does not detect any camera!

That means that either your cameras are not streaming in RTSP or that they are not on the target you are scanning. In most cases, CCTV cameras will be on a private subnetwork. Use the `-s` option to specify your target.

> Cameradar detects my cameras, but does not manage to access them at all!

Maybe your cameras have been configured and the credentials / URL have been changed. Cameradar only guesses using default constructor values. However, you can use your own dictionary in which you just have to add your passwords. To do that, see how the [configuration](#configuration) works. Also, maybe your camera's credentials are not yet known, in which case if you find them it would be very nice to add them to the Cameradar dictionaries to help other people in the future.

> It does not compile :(

You probably missed the part with the dependencies! Use the quick docker deployment, it will be easier and will not pollute your machine with useless dependencies! `;)`

## License

Copyright 2016 Etix Labs

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.

See the License for the specific language governing permissions and limitations under the License.
