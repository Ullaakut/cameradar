# Cameradar

<p align="center">
    <img src="https://raw.githubusercontent.com/Ullaakut/cameradar/master/images/Cameradar.gif" width="100%"/>
</p>

<p align="center">
    <a href="#license">
        <img src="https://img.shields.io/badge/license-Apache-blue.svg?style=flat" />
    </a>
    <a href="https://hub.docker.com/r/ullaakut/cameradar/">
        <img src="https://img.shields.io/docker/pulls/ullaakut/cameradar.svg?style=flat" />
    </a>
    <a href="https://travis-ci.org/Ullaakut/cameradar">
        <img src="https://travis-ci.org/Ullaakut/cameradar.svg?branch=master" />
    </a>
    <a href='https://coveralls.io/github/Ullaakut/cameradar?branch=master'>
        <img src='https://coveralls.io/repos/github/Ullaakut/cameradar/badge.svg?branch=master' alt='Coverage Status' />
    </a>
    <a href="https://golangci.com/r/github.com/ullaakut/cameradar">
        <img src="https://golangci.com/badges/github.com/ullaakut/cameradar.svg" />
    </a>
    <a href="https://goreportcard.com/report/github.com/ullaakut/cameradar">
        <img src="https://goreportcard.com/badge/github.com/ullaakut/cameradar" />
    </a>
    <a href="https://github.com/ullaakut/cameradar/releases/latest">
        <img src="https://img.shields.io/github/release/Ullaakut/cameradar.svg?style=flat" />
    </a>
    <a href="https://godoc.org/github.com/ullaakut/cameradar">
        <img src="https://godoc.org/github.com/ullaakut/cameradar?status.svg" />
    </a>
</p>

## An RTSP stream access tool that comes with its library

### Cameradar allows you to

* **Detect open RTSP hosts** on any accessible target host
* Detect which device model is streaming
* Launch automated dictionary attacks to get their **stream route** (e.g.: `/live.sdp`)
* Launch automated dictionary attacks to get the **username and password** of the cameras
* Retrieve a complete and user-friendly report of the results

<p align="center"><img src="https://raw.githubusercontent.com/Ullaakut/cameradar/master/images/Cameradar.png" width="250"/></p>

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

e.g.: `docker run -t ullaakut/cameradar -t 192.168.100.0/24 -l` will scan the ports 554, 5554 and 8554 of hosts on the 192.168.100.0/24 subnetwork and attack the discovered RTSP streams and will output debug logs.

* `YOUR_TARGET` can be a subnet (e.g.: `172.16.100.0/24`), an IP (e.g.: `172.16.100.10`), or a range of IPs (e.g.: `172.16.100.10-20`).
* If you want to get the precise results of the nmap scan in the form of an XML file, you can add `-v /your/path:/tmp/cameradar_scan.xml` to the docker run command, before `ullaakut/cameradar`.
* If you use the `-r` and `-c` options to specify your custom dictionaries, make sure to also use a volume to add them to the docker container. Example: `docker run -t -v /path/to/dictionaries/:/tmp/ ullaakut/cameradar -r /tmp/myroutes -c /tmp/mycredentials.json -t mytarget`

## Installing the binary on your machine

Only use this solution if for some reason using docker is not an option for you or if you want to locally build Cameradar on your machine.

### Dependencies

* `go`
* `dep`

#### Installing dep

* OSX: `brew install dep` and `brew upgrade dep`
* Others: Download the release package for your OS [here](https://github.com/golang/dep/releases)

### Steps to install

Make sure you installed the dependencies mentionned above.

1. `go get github.com/ullaakut/cameradar`
2. `cd $GOPATH/src/github.com/ullaakut/cameradar`
3. `dep ensure`
4. `cd cameradar`
5. `go install`

The `cameradar` binary is now in your `$GOPATH/bin` ready to be used. See command line options [here](#command-line-options).

## Library

### Dependencies of the library

* `curl-dev` / `libcurl` (depending on your OS)
* `nmap`
* `github.com/pkg/errors`
* `gopkg.in/go-playground/validator.v9`
* `github.com/ullaakut/go-curl`

#### Installing the library

`go get github.com/ullaakut/cameradar`

After this command, the _cameradar_ library is ready to use. Its source will be in:

    $GOPATH/src/pkg/github.com/ullaakut/cameradar

You can use `go get -u` to update the package.

Here is an overview of the exposed functions of this library:

#### Discovery

You can use the cameradar library for simple discovery purposes if you don't need to access the cameras but just to be aware of their existence.

<p align="center"><img  width="90%" src="https://raw.githubusercontent.com/Ullaakut/cameradar/master/images/NmapPresets.png"/></p>
This describes the nmap time presets. You can pass a value between 1 and 5 as described in this table, to the NmapRun function.

#### Attack

If you already know which hosts and ports you want to attack, you can also skip the discovery part and use directly the attack functions. The attack functions also take a timeout value as a parameter.

#### Data models

Here are the different data models useful to use the exposed functions of the cameradar library.

<p align="center"><img width="60%" src="https://raw.githubusercontent.com/Ullaakut/cameradar/master/images/Models.png"/></p>

#### Dictionary loaders

The cameradar library also provides two functions that take file paths as inputs and return the appropriate data models filled.

## Configuration

The **RTSP port used for most cameras is 554**, so you should probably specify 554 as one of the ports you scan. Not specifying any ports to the cameradar application will scan the 554, 5554 and 8554 ports.

`docker run -t --net=host ullaakut/cameradar -p "18554,19000-19010" -t localhost` will scan the ports 18554, and the range of ports between 19000 and 19010 on localhost.

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

* **"-t, --targets"**: Set target. Required. Target can be a file (see [instructions on how to format the file](#format-input-file)), an IP, an IP range, a subnetwork, or a combination of those. Example: `--targets="192.168.1.72,192.168.1.74"`
* **"-p, --ports"**: (Default: `554,5554,8554`) Set custom ports.
* **"-s, --speed"**: (Default: `4`) Set custom nmap discovery presets to improve speed or accuracy. It's recommended to lower it if you are attempting to scan an unstable and slow network, or to increase it if on a very performant and reliable network. You might also want to keep it low to keep your discovery stealthy. See [this for more info on the nmap timing templates](https://nmap.org/book/man-performance.html).
* **"-T, --timeout"**: (Default: `2000`) Set custom timeout value in miliseconds after which an attack attempt without an answer should give up. It's recommended to increase it when attempting to scan unstable and slow networks or to decrease it on very performant and reliable networks.
* **"-r, --custom-routes"**: (Default: `<CAMERADAR_GOPATH>/dictionaries/routes`) Set custom dictionary path for routes
* **"-c, --custom-credentials"**: (Default: `<CAMERADAR_GOPATH>/dictionaries/credentials.json`) Set custom dictionary path for credentials
* **"-o, --nmap-output"**: (Default: `/tmp/cameradar_scan.xml`) Set custom nmap output path
* **"-l, --log"**: Enable debug logs (nmap requests, curl describe requests, etc.)
* **"-h"** : Display the usage information

## Format input file

The file can contain IPs, hostnames, IP ranges and subnetwork, separated by newlines. Example:

```go
0.0.0.0
localhost
192.17.0.0/16
192.168.1.140-255
192.168.2-3.0-255
```

## Environment Variables

### `CAMERADAR_TARGET`

This variable is mandatory and specifies the target that cameradar should scan and attempt to access RTSP streams on.

Examples:

* `172.16.100.0/24`
* `192.168.1.1`
* `localhost`
* `192.168.1.140-255`
* `192.168.2-3.0-255`

### `CAMERADAR_PORTS`

This variable is optional and allows you to specify the ports on which to run the scans.

Default value: `554,5554,8554`

It is recommended not to change these except if you are certain that cameras have been configured to stream RTSP over a different port. 99.9% of cameras are streaming on these ports.

### `CAMERADAR_NMAP_OUTPUT_FILE`

This variable is optional and allows you to specify on which file nmap will write its output.

Default value: `/tmp/cameradar_scan.xml`

This can be useful only if you want to read the files yourself, if you don't want it to write in your `/tmp` folder, or if you want to use only the RunNmap function in cameradar, and do its parsing manually.

### `CAMERADAR_CUSTOM_ROUTES`, `CAMERADAR_CUSTOM_CREDENTIALS`

These variables are optional, allowing to replace the default dictionaries with custom ones, for the dictionary attack.

Default values: `<CAMERADAR_GOPATH>/dictionaries/routes` and `<CAMERADAR_GOPATH>/dictionaries/credentials.json`

### `CAMERADAR_SPEED`

This optional variable allows you to set custom nmap discovery presets to improve speed or accuracy. It's recommended to lower it if you are attempting to scan an unstable and slow network, or to increase it if on a very performant and reliable network. See [this for more info on the nmap timing templates](https://nmap.org/book/man-performance.html).

Default value: `4`

### `CAMERADAR_TIMEOUT`

This optional variable allows you to set custom timeout value in miliseconds after which an attack attempt without an answer should give up. It's recommended to increase it when attempting to scan unstable and slow networks or to decrease it on very performant and reliable networks.

Default value: `2000`

### `CAMERADAR_LOGGING`

This optional variable allows you to enable a more verbose output to have more information about what is going on.

It will output nmap results, cURL requests, etc.

Default: `false`

## Contribution

### Build

#### Docker build

To build the docker image, simply run `docker build -t . cameradar` in the root of the project.

Your image will be called `cameradar` and NOT `ullaakut/cameradar`.

#### Go build

To build the project without docker:

1. Install dep
   * OSX: `brew install dep` and `brew upgrade dep`
   * Others: Download the release package for your OS [here](https://github.com/golang/dep/releases)
2. `dep ensure`
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

See the example in `/cameradar`. You just need to run `go get github.com/ullaakut/cameradar` and to use the `cmrdr` package in your code. You can find the documentation on [godoc](https://godoc.org/github.com/ullaakut/cameradar).

> I want to scan my own localhost for some reason and it does not work! What's going on?

Use the `--net=host` flag when launching the cameradar image, or use the binary by running `go run cameradar/cameradar.go` or [installing it](#installing-the-binary)

> I don't see a colored output :(

You forgot the `-t` flag before `ullaakut/cameradar` in your command-line. This tells docker to allocate a pseudo-tty for cameradar, which makes it able to use colors.

> I don't have a camera but I'd like to try Cameradar!

Simply run `docker run -p 8554:8554 -e RTSP_USERNAME=admin -e RTSP_PASSWORD=12345 -e RTSP_PORT=8554 ullaakut/rtspatt` and then run cameradar and it should guess that the username is admin and the password is 12345. You can try this with any default constructor credentials (they can be found [here](dictionaries/credentials.json))

## Examples

> Running cameradar on your own machine to scan for default ports

`docker run --net=host -t ullaakut/cameradar -t localhost`

> Running cameradar with an input file, logs enabled on port 8554

`docker run -v /tmp:/tmp --net=host -t ullaakut/cameradar -t /tmp/test.txt -p 8554 -l`

## License

Copyright 2017 Ullaakut

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
DEALINGS IN THE SOFTWARE.
