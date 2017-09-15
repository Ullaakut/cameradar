# Cameradar

## An RTSP surveillance camera access multitool

[![cameradar License](https://img.shields.io/badge/license-Apache-blue.svg?style=flat)](#license)
[![Docker Pulls](https://img.shields.io/docker/pulls/ullaakut/cameradar.svg?style=flat)](https://hub.docker.com/r/ullaakut/cameradar/)
[![Build](https://img.shields.io/travis/EtixLabs/cameradar/master.svg?style=flat)](https://travis-ci.org/EtixLabs/cameradar)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/6ab80cfa7069413e8e7d7e18320309e3)](https://www.codacy.com/app/brendan-le-glaunec/cameradar?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=EtixLabs/cameradar&amp;utm_campaign=Badge_Grade)
[![Latest release](https://img.shields.io/github/release/EtixLabs/cameradar.svg?style=flat)](https://github.com/EtixLabs/cameradar/releases/latest)


#### Cameradar allows you to:

* **Detect open RTSP hosts** on any accessible target
* Get their public info (hostname, port, etc.)
* Launch automated dictionary attacks to get their **stream route** (for example /live.sdp)
* Launch automated dictionary attacks to get the **username and password** of the cameras

<p align="center"><img src="https://raw.githubusercontent.com/EtixLabs/cameradar/master/Cameradar.png" width="350"/></p>

## Table of content

- [Docker Image](#docker-image)
- [Configuration](#configuration)
- [Output](#output)
- [Check camera access](#check-camera-access)
- [Command line options](#command-line-options)
- [Contribution](#contribution)
- [Frequently Asked Questions](#frequently-asked-questions)
- [License](#license)

## Docker Image

Install docker on your machine, and run the following command:

```
docker run  --net=host \
            -v /tmp/thumbs/:/tmp/thumbs \
            -e CAMERAS_TARGET=your_target \
            ullaakut/cameradar
```

* `your_target` can be a subnet (e.g.: `172.16.100.0/24`) or even an IP (e.g.: `172.16.100.10`), a range of IPs (e.g.: `172.16.100.10-172.16.100.20`) or a mix of all those separated by commas (e.g.: `172.17.100.0/24,172.16.100.10-172.16.100.20,0.0.0.0`).
* if you want to get the precise results of the nmap scan in the form of an XML file, you can add `-v /your/path:/tmp/cameradar_scan.xml` to the docker run command, before `ullaakut/cameradar`.

Check [Cameradar's readme on the Docker Hub](https://hub.docker.com/r/ullaakut/cameradar/) for more information and more command-line options.

For more complex use of the Docker image, see the `Environment variables` part of [Cameradar's readme on the Docker Hub](https://hub.docker.com/r/ullaakut/cameradar/).

### Configuration

The **RTSP port used for most cameras is 554**, so you should probably specify 554 as one of the ports you scan. Not specifying any ports to the cameraccess application will scan the 554 and 8554 ports.

You **can use your own files for the ids and routes dictionaries** used to attack the cameras, but the Cameradar repository already gives you a good base that works with most cameras, in the `/dictionaries` folder.

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
   "device" : "Vivotek FD9381-HTV",
   "route" : "/live.sdp",
   "username" : "admin"
}
```

## Check camera access

If you have [VLC Media Player](http://www.videolan.org/vlc/), you should be able to use the GUI to connect to the RTSP stream using this format : `rtsp://username:password@address:port/route`

With the above result, the RTSP URL would be `rtsp://admin:123456@173.16.100.45:554/live.sdp`

If you're still in your console however, you can go even faster by using **vlc in commmand-line** and just run `vlc rtsp://username:password@address:port/route` with the camera's info instead of the placeholders.

## Command line options

* **"-c"** : Set a custom path to the configuration file (-c /path/to/conf)
* **"-s"** : Set custom target (overrides configuration) : You can use this argument in many ways, using a subnet (e.g.: `172.16.100.0/24`) or even an IP (e.g.: `172.16.100.10`), a range of IPs (e.g.: `172.16.100.10-172.16.100.20`) or a mix of all those (e.g.: `172.17.100.0/24,172.16.100.10-172.16.100.20,0.0.0.0`)
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
