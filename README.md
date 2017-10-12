# Cameradar

<p align="center"><img src="https://raw.githubusercontent.com/EtixLabs/cameradar/master/images/Cameradar.gif" width="100%"/></p>

## An RTSP stream access tool that comes with its library

[![cameradar License](https://img.shields.io/badge/license-Apache-blue.svg?style=flat)](#license)
[![Docker Pulls](https://img.shields.io/docker/pulls/ullaakut/cameradar.svg?style=flat)](https://hub.docker.com/r/ullaakut/cameradar/)
[![Build](https://img.shields.io/travis/EtixLabs/cameradar/master.svg?style=flat)](https://travis-ci.org/EtixLabs/cameradar)
[![Go Report Card](https://goreportcard.com/badge/github.com/EtixLabs/cameradar)](https://goreportcard.com/report/github.com/EtixLabs/cameradar)
[![GoDoc](https://godoc.org/github.com/EtixLabs/cameradar?status.svg)](https://godoc.org/github.com/EtixLabs/cameradar)
[![Latest release](https://img.shields.io/github/release/EtixLabs/cameradar.svg?style=flat)](https://github.com/EtixLabs/cameradar/releases/latest)

### Cameradar allows you to

* **Detect open RTSP hosts** on any accessible target host
* Detect which device model is streaming
* Launch automated dictionary attacks to get their **stream route** (e.g.: `/live.sdp`)
* Launch automated dictionary attacks to get the **username and password** of the cameras
* Retrieve a complete and user-friendly report of the results

<p align="center"><img src="https://raw.githubusercontent.com/EtixLabs/cameradar/master/images/Cameradar.png" width="250"/></p>

## Table of content

* [Docker Image](#docker-image)
* [Configuration](#configuration)
* [Output](#output)
* [Check camera access](#check-camera-access)
* [Command line options](#command-line-options)
* [Contribution](#contribution)
* [Frequently Asked Questions](#frequently-asked-questions)
* [License](#license)

## Docker Image for Cameradar

Install [docker](https://docs.docker.com/engine/installation/) on your machine, and run the following command:

```bash
docker run -t ullaakut/cameradar -t <target> <other command-line options>
```

[See command-line options](#command-line-options).

e.g.: `docker run -t ullaakut/cameradar -t 192.168.100.0/24 -l` will scan the ports 554 and 8554 of hosts on the 192.168.100.0/24 subnetwork and attack the discovered RTSP streams and will output debug logs.

* `YOUR_TARGET` can be a subnet (e.g.: `172.16.100.0/24`), an IP (e.g.: `172.16.100.10`), or a range of IPs (e.g.: `172.16.100.10-20`).
* If you want to get the precise results of the nmap scan in the form of an XML file, you can add `-v /your/path:/tmp/cameradar_scan.xml` to the docker run command, before `ullaakut/cameradar`.
* If you use the `-r` and `-c` options to specify your custom dictionaries, make sure to also use a volume to add them to the docker container. Example: `docker run -t -v /path/to/dictionaries/:/tmp/ ullaakut/cameradar -r /tmp/myroutes -c /tmp/mycredentials.json -t mytarget`

## Installing the binary

### Dependencies

* `go`
* `glide`

#### Installing [glide](https://github.com/Masterminds/glide)

* OSX: `brew install glide`
* Linux: `curl https://glide.sh/get | sh`
* Windows: Download the binary package [here](https://github.com/Masterminds/glide/releases)

### Steps to install

Make sure you installed the dependencies mentionned above.

1. `go get github.com/EtixLabs/cameradar`
2. `cd $GOPATH/src/github.com/EtixLabs/cameradar`
3. `glide install`
4. `cd cameradar`
5. `go install`

The `cameradar` binary is now in your `$GOPATH/bin` ready to be used. See command line options [here](#command-line-options).

## Library

### Dependencies of the library

* `curl-dev` / `libcurl` (depending on your OS)
* `nmap`
* `github.com/pkg/errors`
* `gopkg.in/go-playground/validator.v9`
* `github.com/andelf/go-curl`

#### Installing the library

`go get github.com/EtixLabs/cameradar`

After this command, the *cameradar* library is ready to use. Its source will be in:

    $GOPATH/src/pkg/github.com/EtixLabs/cameradar

You can use `go get -u` to update the package.

Here is an overview of the exposed functions of this library:

#### Discovery

You can use the cameradar library for simple discovery purposes if you don't need to access the cameras but just to be aware of their existence.

<p align="center"><img  width="90%" src="https://raw.githubusercontent.com/EtixLabs/cameradar/master/images/NmapPresets.png"/></p>
This describes the nmap time presets. You can pass a value between 1 and 5 as described in this table, to the NmapRun function.

#### Attack

If you already know which hosts and ports you want to attack, you can also skip the discovery part and use directly the attack functions. The attack functions also take a timeout value as a parameter.

#### Data models

Here are the different data models useful to use the exposed functions of the cameradar library.

<p align="center"><img width="60%" src="https://raw.githubusercontent.com/EtixLabs/cameradar/master/images/Models.png"/></p>

#### Dictionary loaders

The cameradar library also provides two functions that take file paths as inputs and return the appropriate data models filled.

## Configuration

The **RTSP port used for most cameras is 554**, so you should probably specify 554 as one of the ports you scan. Not specifying any ports to the cameradar application will scan the 554 and 8554 ports.

`docker run -t ullaakut/cameradar -p "18554,19000-19010" -t localhost` will scan the ports 18554, and the range of ports between 19000 and 19010 on localhost.

You **can use your own files for the ids and routes dictionaries** used to attack the cameras, but the Cameradar repository already gives you a good base that works with most cameras, in the `/dictionaries` folder.

```bash
docker run -t -v /my/folder/with/dictionaries:/tmp/dictionaries \
           ullaakut/cameradar \
           -r "/tmp/dictionaries/my_routes" \
           -c "/tmp/dictionaries/my_credentials.json" \
           -t 172.19.124.0/24
```

This will put the contents of your folder containing dictionaries in the docker image and will use it for the dictionary attack instead of the default dictionaries provided in the cameradar repo.

## Check camera access

If you have [VLC Media Player](http://www.videolan.org/vlc/), you should be able to use the GUI or the command-line to connect to the RTSP stream using this format : `rtsp://username:password@address:port/route`

With the above result, the RTSP URL would be `rtsp://admin:12345@173.16.100.45:554/live.sdp`

## Command line options

* **"-t, --target"**: Set custom target. Required.
* **"-p, --ports"**: (Default: `554,8554`) Set custom ports.
* **"-s, --speed"**: (Default: `4`) Set custom nmap discovery presets to improve speed or accuracy. It's recommended to lower it if you are attempting to scan an unstable and slow network, or to increase it if on a very performant and reliable network. See [this for more info on the nmap timing templates](https://nmap.org/book/man-performance.html).
* **"-T, --timeout"**: (Default: `1000`) Set custom timeout value in miliseconds after which an attack attempt without an answer should give up.
* **"-r, --custom-routes"**: (Default: `dictionaries/routes`) Set custom dictionary path for routes
* **"-c, --custom-credentials"**: (Default: `dictionaries/credentials.json`) Set custom dictionary path for credentials
* **"-o, --nmap-output"**: (Default: `/tmp/cameradar_scan.xml`) Set custom nmap output path
* **"-l, --log"**: Enable debug logs (nmap requests, curl describe requests, etc.)
* **"-h"** : Display the usage information

## Contribution

### Build

#### Docker build

To build the docker image, simply run `docker build -t . cameradar` in the root of the project.

Your image will be called `cameradar` and NOT `ullaakut/cameradar`.

#### Go build

To build the project without docker:

1. install [glide](https://github.com/Masterminds/glide)
    * OSX: `brew install glide`
    * Linux: `curl https://glide.sh/get | sh`
    * Windows: Download the binary package [here](https://github.com/Masterminds/glide/releases)
2. `glide install`
3. `go build` to build the library
4. `cd cameradar && go build` to build the binary

The cameradar binary is now in the root of the directory.

See [the contribution document](/CONTRIBUTING.md) to get started.

## Frequently Asked Questions

> Cameradar does not detect any camera!

That means that either your cameras are not streaming in RTSP or that they are not on the target you are scanning. In most cases, CCTV cameras will be on a private subnetwork, isolated from the internet. Use the `-t` option to specify your target.

> Cameradar detects my cameras, but does not manage to access them at all!

Maybe your cameras have been configured and the credentials / URL have been changed. Cameradar only guesses using default constructor values if a custom dictionary is not provided. You can use your own dictionaries in which you just have to add your credentials and RTSP routes. To do that, see how the [configuration](#configuration) works. Also, maybe your camera's credentials are not yet known, in which case if you find them it would be very nice to add them to the Cameradar dictionaries to help other people in the future.

> What happened to the C++ version?

You can still find it under the 1.1.4 tag on this repo, however it was less performant and stable than the current version written in Golang.

> How to use the Cameradar library for my own project?

See the example in `/cameradar`. You just need to run `go get github.com/EtixLabs/cameradar` and to use the `cmrdr` package in your code.

> I want to scan my own localhost for some reason and it does not work! What's going on?

Use the `--net=host` flag when launching the cameradar image, or use the binary by running `go run cameradar/cameradar.go` or [installing it](#installing-the-binary)

> I don't see a colored output :(

You forgot the `-t` flag before `ullaakut/cameradar` in your command-line. This tells docker to allocate a pseudo-tty for cameradar, which makes it able to use colors.

## Known issues

* When running Cameradar in a docker container, specifying multiple targets does not work. Using subnetworks (such as `182.49.20.0/24`) or ranges (`182.49.20.0-44`) works.
* There is currently no way to use environment variables instead of command-line arguments in Cameradar. This will be done at some point, but isn't a priority right now.

## License

Copyright 2017 Etix Labs

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.

See the License for the specific language governing permissions and limitations under the License.
