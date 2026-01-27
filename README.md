## Cameradar

<p align="center">
    <a href="#license">
        <img src="https://img.shields.io/badge/license-Apache-blue.svg?style=flat" />
    </a>
    <a href="https://hub.docker.com/r/ullaakut/cameradar/">
        <img src="https://img.shields.io/docker/pulls/ullaakut/cameradar.svg?style=flat" />
    </a>
    <a href="https://github.com/Ullaakut/cameradar/actions">
        <img src="https://img.shields.io/github/actions/workflow/status/Ullaakut/cameradar/build.yaml" />
    </a>
    <a href='https://coveralls.io/github/Ullaakut/cameradar?branch=master'>
        <img src='https://coveralls.io/repos/github/Ullaakut/cameradar/badge.svg?branch=master' alt='Coverage Status' />
    </a>
    <a href="https://goreportcard.com/report/github.com/ullaakut/cameradar">
        <img src="https://goreportcard.com/badge/github.com/ullaakut/cameradar" />
    </a>
    <a href="https://github.com/ullaakut/cameradar/releases/latest">
        <img src="https://img.shields.io/github/release/Ullaakut/cameradar.svg?style=flat" />
    </a>
    <a href="https://pkg.go.dev/github.com/ullaakut/cameradar">
        <img src="https://godoc.org/github.com/ullaakut/cameradar?status.svg" />
    </a>
</p>

## RTSP stream access tool

Cameradar scans RTSP endpoints on authorized targets, and uses dictionary attacks to bruteforce their credentials and routes.

### What Cameradar does

- Detects open RTSP hosts on accessible targets.
- Detects the device model that streams the RTSP feed.
- Attempts dictionary-based discovery of stream routes (for example, `/live.sdp`).
- Attempts dictionary-based discovery of camera credentials.
- Produces a report of findings.

<p align="center"><img src="images/Cameradar.png" width="250"/></p>

## Table of contents

- [Quick start with Docker](#quick-start-with-docker)
- [Install the binary](#install-the-binary)
- [Install on Android (Termux)](#install-on-android-termux)
- [Configuration](#configuration)
- [Security and responsible use](#security-and-responsible-use)
- [Output](#output)
- [Check camera access](#check-camera-access)
- [Command-line options](#command-line-options)
- [Input file format](#input-file-format)
- [Environment variables](#environment-variables)
- [Build and contribute](#build-and-contribute)
- [Frequently asked questions](#frequently-asked-questions)
- [Examples](#examples)
- [License](#license)

---

<p align="center"><img src="images/example.gif"/></p>

## Quick start with Docker

Install [Docker](https://docs.docker.com/engine/installation/) and run:

```bash
docker run --rm -t --net=host ullaakut/cameradar --targets <target>
```

Example:

```bash
docker run --rm -t --net=host ullaakut/cameradar --targets 192.168.100.0/24
```

This scans ports 554, 5554, and 8554 on the target subnet.
It attempts to enumerate RTSP streams.
For all options, see [command-line options](#command-line-options).

- Targets can be CIDRs, IPs, IP ranges or a hostname.
    - Subnet: `172.16.100.0/24`
    - IP: `172.16.100.10`
    - Host: `localhost`
    - Range: `172.16.100.10-20`

- To use custom dictionaries, mount them and pass both flags:

    ```bash
    docker run --rm -t --net=host \
        -v /path/to/dictionaries:/tmp/dictionaries \
        ullaakut/cameradar \
        --custom-routes /tmp/dictionaries/my_routes \
        --custom-credentials /tmp/dictionaries/my_credentials.json \
        --targets 192.168.100.0/24
    ```

## Install the binary

Use this option if Docker is not available or if you want a local build.

### Dependencies

- Go 1.25 or later

### Steps

1. `go install github.com/Ullaakut/cameradar/v6/cmd/cameradar@latest`

The `cameradar` binary is now in your `$GOPATH/bin`.
For available flags, see [command-line options](#command-line-options).

## Install on Android (Termux)

These steps summarize a working Termux setup for Android.
Use Termux 117 from F-Droid or the official Termux site, not Google Play.

### 1) Set up Termux and Alpine

Install the required packages in Termux:

```bash
pkg update
pkg install mc wget git nmap proot-distro
```

Install Alpine and log in:

```bash
proot-distro install alpine
proot-distro login alpine
```

### 2) Install build tools in Alpine

```bash
apk add wget git go gcc clang musl-dev make
```

### 3) Build Cameradar

Create a module path and clone the repo:

```bash
mkdir -p go/pkg/mod/github.com/Ullaakut
cd go/pkg/mod/github.com/Ullaakut
git clone https://github.com/Ullaakut/cameradar.git
cd cameradar/cmd/cameradar
go install
```

### 4) Run Cameradar

Copy dictionaries and run the binary:

```bash
mkdir -p /tmp
cp -r ../../dictionaries /tmp/dictionaries
/go/bin/cameradar --targets=<target> --custom-credentials=/tmp/dictionaries/credentials.json --custom-routes=/tmp/dictionaries/routes --ui=plain --debug 
```

Replace `<target>` with an IP, range, host or subnet you are authorized to test.

## Configuration

The default RTSP ports are `554`, `5554`, `8554`.
If you do not specify ports, Cameradar uses those.

Example of scanning custom ports:

```bash
docker run --rm -t --net=host \
    ullaakut/cameradar \
    --ports "18554,19000-19010" \
    --targets localhost
```

You can replace the default dictionaries with your own routes and credentials files.
The repository provides baseline dictionaries in the `dictionaries` folder.

```bash
docker run --rm -t --net=host \
    -v /my/folder/with/dictionaries:/tmp/dictionaries \
    ullaakut/cameradar \
    --custom-routes /tmp/dictionaries/my_routes \
    --custom-credentials /tmp/dictionaries/my_credentials.json \
    --targets 172.19.124.0/24
```

### Skip discovery with `--skip-scan`

If you already know the RTSP endpoints, you can skip discovery and treat each
target and port as a stream candidate. This mode does not run nmap and can be
useful on restricted networks or when you want to attack a known inventory.

Skipping discovery means:

- Cameradar does not run nmap and does not detect device models.
- Targets resolve to IP addresses. Hostnames resolve via DNS.
- CIDR blocks and IPv4 ranges expand to every address in the range.
- Large ranges create many targets, so use them carefully.

Example:

```bash
docker run --rm -t --net=host \
    ullaakut/cameradar \
    --skip-scan \
    --ports "554,8554" \
    --targets 192.168.1.10
```

In this example, Cameradar attempts dictionary attacks against
ports 554 and 8554 of `192.168.1.10`.

## Security and responsible use

Cameradar is a penetration testing tool.
Only scan networks and devices you own or have explicit permission to test.
Do not use this tool to access unauthorized systems or streams.
If you are unsure, stop and get written approval before scanning.

## Output

Cameradar presents results in a readable terminal UI.
It logs findings to the console.
The report includes discovered hosts, identified device models, and valid routes or credentials.
If you specify a path for the `--output` flag, Cameradar also writes an M3U playlist with the discovered streams.

## Check camera access

Use [VLC Media Player](http://www.videolan.org/vlc/) to connect to a stream:

`rtsp://username:password@address:port/route`

## Input file format

The file can contain IPs, hostnames, IP ranges, and subnets.
Separate entries with newlines.
Example:

```text
0.0.0.0
localhost
192.17.0.0/16
192.168.1.140-255
192.168.2-3.0-255
```

When you use `--skip-scan`, Cameradar expands each entry into explicit IP
addresses before building the target list.

## Options

### `TARGETS` / `--targets` / `-t`

This variable is required.
It specifies the target that Cameradar scans and attempts to access.

Examples:

* `172.16.100.0/24`
* `192.168.1.1`
* `localhost`
* `192.168.1.140-255`
* `192.168.2-3.0-255`

### `PORTS` / `--ports` / `-p`

This variable is optional and allows you to specify the ports to scan.

Default value: `554,5554,8554`

Change these only if you are sure cameras stream over different ports.
Most cameras use these defaults.

### `CUSTOM_ROUTES` / `--custom-routes` / `-r`

This option is optional.
It replaces the default routes dictionary used for the dictionary attack.

If unset, Cameradar uses the built-in routes dictionary.

### `CUSTOM_CREDENTIALS` / `--custom-credentials` / `-c`

This option is optional.
It replaces the default credentials dictionary used for the dictionary attack.

If unset, Cameradar uses the built-in credentials dictionary.

### `SCAN_SPEED` / `--scan-speed` / `-s`

This optional variable sets nmap discovery presets for speed or accuracy.
Lower it on slow networks and raise it on fast networks.
See [nmap timing templates](https://nmap.org/book/man-performance.html).

Default value: `4`

### `SKIP_SCAN` / `--skip-scan`

This optional flag skips network discovery and assumes every target and port
pair is an RTSP stream.

Use it when you already know the RTSP endpoints or when discovery is blocked.
For best results, specify only RTSP ports.

Default value: `false`

### `ATTACK_INTERVAL` / `--attack-interval` / `-I`

This optional variable sets a delay between attacks.
Increase it for networks that may block brute-force attempts.
Default: no delay.

Default value: `0ms`

### `TIMEOUT` / `--timeout` / `-T`

This optional variable sets the timeout for requests sent to the cameras.
Increase it for slow networks and decrease it for fast networks.

Default value: `2000ms`

### `DEBUG` / `--debug` / `-d`

This optional variable enables more verbose output.

It outputs nmap results, cURL requests, and more.

Default: `false`

### `UI` / `--ui`

This option selects the UI mode.

* `auto` selects `tui` if your terminal is interactive, `plain` otherwise
* `tui` shows a fullscreen interface with a progress bar and shows the results in a table
* `plain` logs the steps taken by cameradar as plain text and is meant to be used by non-interactive terminals

Supported values: `auto`, `tui`, `plain`

Default: `auto`

### `OUTPUT` / `--output`

This optional variable writes an M3U playlist of the discovered streams to the given file path.

Example: `/tmp/cameradar.m3u`

## Build and contribute

### Docker build

Run the following command in the repository root:

`docker build . -t cameradar`

The resulting image is named `cameradar`.

### Go build

1. `go install github.com/Ullaakut/cameradar/v6/cmd/cameradar@latest`

The `cameradar` binary is now in `$GOPATH/bin/cameradar`.

## Frequently asked questions

> Cameradar does not detect any camera!

This usually means the cameras are not streaming over RTSP.
It can also mean the targets are not in your scan range.
CCTV cameras are often on private subnets.
Use `-t` to set the correct targets.
If you still see no results, open an issue with device details.

> Cameradar detects my cameras, but does not manage to access them!

The camera configuration may have changed, so defaults do not match.
Cameradar uses defaults unless you provide custom dictionaries.
Add your credentials and routes, then follow the [configuration](#configuration) section.

> What happened to the C++ version?

The 1.1.4 tag contains the legacy C++ implementation.
It is slower and less stable than the Go version, so it is not recommended to use.

> I want to scan my local network or my own machine, and it does not work! What's going on?

Use `--net=host` when running the Docker image, or use the installed binary.

> I don't have a camera, but I'd like to try Cameradar!

Run the following container, then run Cameradar against it:

`docker run -p 8554:8554 -e RTSP_USERNAME=admin -e RTSP_PASSWORD=12345 -e RTSP_PORT=8554 ullaakut/rtspatt`

Cameradar should discover the `admin` / `12345` credentials.
You can try other default credentials listed in the dictionaries.

> What authentication types does Cameradar support?

Cameradar supports both basic and digest authentication.

## Examples

> Running cameradar on your own machine to scan for default ports

`docker run --rm -t --net=host ullaakut/cameradar --targets localhost`

> Running cameradar with an input file, logs enabled on port 8554

`docker run --rm -t --net=host -v /tmp:/tmp ullaakut/cameradar --targets /tmp/test.txt --ports 8554`

> Running cameradar on a subnetwork with custom dictionaries, on ports 554, 5554 and 8554

`docker run --rm -t --net=host -v /tmp:/tmp ullaakut/cameradar --targets 192.168.0.0/24 --custom-credentials "/tmp/dictionaries/credentials.json" --custom-routes "/tmp/dictionaries/routes" --ports 554,5554,8554`

## License

Copyright 2026 Ullaakut

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
