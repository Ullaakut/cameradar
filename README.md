## Cameradar

<p align="center">
    <a href="#license">
        <img src="https://img.shields.io/badge/license-MIT-blue.svg?style=flat" />
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
- [Command-line options and environment variables](#command-line-options-and-environment-variables)
- [Input file format](#input-file-format)
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
For all options, see [Configuration reference](https://github.com/Ullaakut/cameradar/wiki/Configuration-Reference).

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
For available flags, see [Configuration reference](https://github.com/Ullaakut/cameradar/wiki/Configuration-Reference).

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
target and port as a stream candidate. This mode does not run discovery and can be
useful on restricted networks or when you want to attack a known inventory.

Skipping discovery means:

- Cameradar does not run discovery and does not detect device models.
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

### Choose the discovery scanner with `--scanner`

Cameradar supports two discovery backends:

- `nmap` (default)
- `masscan`

Use `nmap` when you want more reliable RTSP discovery: it performs service
identification and can better distinguish RTSP from other open ports.

Use `masscan` when scanning very large networks: it is generally faster and
more efficient at scale, but it does not provide service discovery.

```bash
docker run --rm -t --net=host \
    ullaakut/cameradar \
    --scanner masscan \
    --ports "554,8554" \
    --targets 192.168.1.0/24
```

> [!WARNING]  
> `--scan-speed` only applies to the `nmap` scanner.

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

## Command-line options and environment variables

The complete CLI and environment variable reference is maintained in [Configuration reference](https://github.com/Ullaakut/cameradar/wiki/Configuration-Reference).

This includes all supported flags, defaults, accepted values, and env var mapping.

## Build and contribute

### Docker build

Run the following command in the repository root:

`docker build . -t cameradar`

The resulting image is named `cameradar`.

### Go build

1. `go install github.com/Ullaakut/cameradar/v6/cmd/cameradar@latest`

The `cameradar` binary is now in `$GOPATH/bin/cameradar`.

## Frequently asked questions

See [Troubleshooting & FAQ](https://github.com/Ullaakut/cameradar/wiki/Troubleshooting-%26-FAQ)

## Examples

> Running cameradar on your own machine to scan for default ports

`docker run --rm -t --net=host ullaakut/cameradar --targets localhost`

> Running cameradar with an input file, logs enabled on port 8554

`docker run --rm -t --net=host -v /tmp:/tmp ullaakut/cameradar --targets /tmp/test.txt --ports 8554`

> Running cameradar on a subnetwork with custom dictionaries, on ports 554, 5554 and 8554

`docker run --rm -t --net=host -v /tmp:/tmp ullaakut/cameradar --targets 192.168.0.0/24 --custom-credentials "/tmp/dictionaries/credentials.json" --custom-routes "/tmp/dictionaries/routes" --ports 554,5554,8554`

> Running cameradar with masscan discovery

`docker run --rm -t --net=host ullaakut/cameradar --scanner masscan --targets 192.168.0.0/24 --ports 554,8554`

## License

Copyright 2026 Ullaakut

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
