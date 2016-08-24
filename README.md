# Cameradar

## An RTSP surveillance camera access multitool

[![cameradar License](https://img.shields.io/badge/license-Apache-blue.svg)](#license)
[![Latest release](https://img.shields.io/badge/release-1.0.2-green.svg)](https://github.com/EtixLabs/cameradar/releases/latest)


#### Cameradar allows you to:

* **Detect open RTSP hosts** on any accessible subnetwork
* Get their public info (hostname, port, camera model, etc.)
* Bruteforce your way into them to get their **stream route** (for example /live.sdp)
* Bruteforce your way into them to get the **username and password** of the cameras
* **Generate thumbnails** from them to check if the streams are valid and to have a quick preview of their content
* Try to create a Gstreamer pipeline to check if they are **properly encoded**
* Print a summary of all the informations Cameradar could get

#### And all of this in a _single command-line_.

Of course, you can also call for individual tasks if you plug in a Database to Cameradar, but for now this repo only contains a basic cache manager. You can however create your own by following the simple example of the **dumb cache manager**.

## Table of content

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
- [Under the hood](#under-the-hood)
- [Contribution](#contribution)
- [Next improvements](#next-improvements)
- [License](#license)

## Quick install

The quick install uses docker to build Cameradar without polluting your machine with dependencies and makes it easy to deploy Cameradar in a few commands. **However, it may require networking knowledge, as your docker containers will need access to the cameras subnetwork.**

### Dependencies

The only dependencies are `docker`, `docker-tools`, `git` and `make`.

### Five steps guide

1. `git clone https://github.com/EtixLabs/cameradar.git`
2. Go into the Cameradar repository, then to the `deployment` directory
3. Tweak the `conf/cameradar.conf.json` as you need (see [the onfiguration guide here](#configuration) for more information)
4. Run `docker-compose build cameradar` to build the cameradar container
5. Run `docker-compose up cameradar` to launch Cameradar

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
2. Go into the Cameradar repository, create a directory named `build` and go in it
3. In the build directory, run `cmake ..` This will generate the Makefiles you need to build Cameradar
4. Run the command `make`
5. This should compile Cameradar. Go into the `cameradar_standalone` directory
6. You can now customize the `conf/cameradar.conf.json` file to set the subnetworks and specific ports you want to scan, as well as the thumbnail generation path. More information will be given about the configuration file in another part of this document.
7. You are now ready to launch Cameradar by launching `./cameradar` in the cameradar_standalone directory.

## Advanced Docker deployment

### Dependencies

The only dependencies are `docker` and `docker-compose`.

### Deploy a custom version of Cameradar

2. Go into the Cameradar repository, create a directory named `build` and go in it
3. In the build directory, run `cmake .. -DCMAKE_BUILD_TYPE=Release` This will generate the Makefiles you need to build Cameradar
4. Run the command `make package` to compile it into a package
5. Copy your package into the `deployment` directory
6. Run `docker-compose build cameradar` to build the cameradar container using your custom package
5. Run `docker-compose up cameradar` to launch Cameradar


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
  "subnets" : "SUBNET1,SUBNET2,SUBNET3,[...]",
  "ports" : "PORT1,PORT2,[...]",
  "rtsp_url_file" : "conf/url.json",
  "rtsp_ids_file" : "conf/ids.json",
  "thumbnail_storage_path" : "/valid/path/to/a/storage/directory",
  "cache_manager_path" : "../cache_managers/dumb_cache_manager",
  "cache_manager_name" : "dumb"
}
```

The subnetworks should be passed separated by commas only, and their subnet format should be the same as used in nmap.
```json
"subnets" : "172.100.16.0/24,172.100.17.0/24,localhost,192.168.1.13"
```

The **RTSP ports for most cameras are 554**, so you should probably specify 554 as one of the ports you scan. Not giving any ports in the configuration will scan every port of every host found on the subnetworks..How is formatted Cameradar's result
You **can use your own files for the ids and routes dictionaries** used to bruteforce the cameras, but the Cameradar repo already gives you a good base that works with most cameras.

The thumbnail storage path should be a **valid and accessible directory** in which the thumbnails will be stored.

The cache manager path and name variables are used to change the cache manager you want to load into Cameradar. If you want to, you can code your own cache manager using a database, a file, a remote server, [...]. Feel free to share it by creating a merge request on this repo if you developed a generic manager (It must not be specific to your company's infrastructure).

## Output

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

If you have vlc, you should be able to use the GUI to connect to the RTSP stream using this format : `username:password@address:port/route`

With the above result, the RTSP URL would be `admin:123456@173.16.100.45:554/live.sdp`

If you're still in your console however, you can go even faster by using **vlc in commmand-line** and just run `vlc username:password@address:port/route` with the camera's info instead of the placeholders.

## Command line options

* **"-c"** : Set a custom path to the configuration file (-c /path/to/conf)
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
* **"-b"** : Launch the bruteforce tool on all discovered devices
  * Needs either to be launched with the -d option or to use an advanced cache manager (DB, file, ...) with data already present
* **"-t"** : Generate thumbnails from detected cameras
  * Needs either to be launched with the -d option or to use an advanced cache manager (DB, file, ...) with data already present
* **"-g"** : Check if the stream can be opened with GStreamer
  * Needs either to be launched with the -d option or to use an advanced cache manager (DB, file, ...) with data already present
* **"-v"** : Display Cameradar's version
* **"-h"** : Display this help

## Under the hood

Cameradar uses **nmap** to map all of the subnetworks you specified in the configuration file (_cameradar.conf.json_), then parses its result to get all of the open RTSP streams that were detected.

After that, it uses **cURL** to send requests to the cameras and to try routes and ids for each camera until it is accessed or until all of the most used routes/ids (that you can modify in _conf/ids.json_ and _conf/url.json_) were tried

Then, it uses **FFMPEG** to generate a lightweight thumbnail from the stream, which you could use to get a quick preview of the camera's view.

Finally, it tries to access the stream using a simple **Gstreamer pipeline** to check for the stream's encoding.

The output of Cameradar will be printed on the standard output and will also be accessible in the result.json file.

Cameradar uses **nmap** to map all of the subnetworks you specified in the configuration file (_cameradar.conf.json_), then parses its result to get all of the open RTSP streams that were detected.

After that, it uses **cURL** to send requests to the cameras and to try routes and ids for each camera until it is accessed or until all of the most used routes/ids (that you can modify in _conf/ids.json_ and _conf/url.json_) were tried

Then, it uses **FFMPEG** to generate a lightweight thumbnail from the stream, which you could use to get a quick preview of the camera's view.

Finally, it tries to access the stream using a simple **Gstreamer pipeline** to check for the stream's encoding.

The output of Cameradar will be printed on the standard output and will also be accessible in the result.json file.

## Contribution

Well there are many things we could code in order to add features to Cameradar. Adding other protocols than RTSP would be really cool, as well as making generic cache managers. Creating an HTTP server with an API that would launch cameradar upon recieving requests ans answer with Cameradar's result would also be potentially really useful.

If you're not into software development or not into C++, even updating the dictionaries would be a really cool contribution! Just make sure the ids and routes you add are **default constructor credentials** and not custom credentials.

If you have other cool ideas, feel free to share them with me at brendan.leglaunec@etixgroup.com !

## Next improvements

- [x] Add a docker deployment to avoid the current deps hell
- [x] Development of a MySQL cache manager
- [ ] Development of a JSON file cache manager
- [ ] Development of an XML file cache manager

## License

Copyright 2016 Etix Labs

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.

See the License for the specific language governing permissions and limitations under the License.
