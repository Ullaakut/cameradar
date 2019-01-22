# nmap

<p align="center">
    <img width="350" src="img/logo.png"/>
<p>

<p align="center">
    <a href="LICENSE">
        <img src="https://img.shields.io/badge/license-MIT-blue.svg?style=flat" />
    </a>
    <a href="https://godoc.org/github.com/Ullaakut/nmap">
        <img src="https://godoc.org/github.com/Ullaakut/cameradar?status.svg" />
    </a>
    <a href="https://goreportcard.com/report/github.com/ullaakut/nmap">
        <img src="https://goreportcard.com/badge/github.com/ullaakut/nmap">
    </a>
    <a href="https://travis-ci.org/Ullaakut/nmap">
        <img src="https://travis-ci.org/Ullaakut/nmap.svg?branch=master">
    </a>
    <a href="https://coveralls.io/github/Ullaakut/nmap?branch=master">
        <img src="https://coveralls.io/repos/github/Ullaakut/nmap/badge.svg?branch=master">
    </a>
<p>

This library aims at providing idiomatic `nmap` bindings for go developers, in order to make it easier to write security audit tools using golang.

<!-- It allows not only to parse the XML output of nmap, but also to get the output of nmap as it is running, through a channel. This can be useful for computing a scan's progress, or simply displaying live information to your users. -->

## It's currently a work in progress

This paragraph won't be removed until the library is ready to be used and properly documented.

## Supported features

- [x] All of `nmap`'s options as `WithXXX` methods.
- [x] Cancellable contexts support.
- [x] [Idiomatic go filters](examples/service_detection/main.go#L19).
- [x] Helpful enums for most nmap commands. (time templates, os families, port states, etc.)
- [x] Complete documentation of each option, mostly insipred from nmap's documentation.

## TODO

- [ ] Examples of usage - Work in progress (4/7 examples so far)
- [ ] Complete unit tests - Work in progress (95% coverage so far)
- [ ] Asynchronous scan

## Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/Ullaakut/nmap"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    // Equivalent to `/usr/local/bin/nmap -p 80,443,843 google.com facebook.com youtube.com`,
    // with a 5 minute timeout.
    scanner, err := nmap.NewScanner(
        nmap.WithTargets("google.com", "facebook.com", "youtube.com"),
        nmap.WithPorts("80,443,843"),
        nmap.WithContext(ctx),
    )
    if err != nil {
        log.Fatalf("unable to create nmap scanner: %v", err)
    }

    result, err := scanner.Run()
    if err != nil {
        log.Fatalf("unable to run nmap scan: %v", err)
    }

    // Use the results to print an example output
    for _, host := range result.Hosts {
        if len(host.Ports) == 0 || len(host.Addresses) == 0 {
            continue
        }

        fmt.Printf("Host %q:\n", host.Addresses[0])

        for _, port := range host.Ports {
            fmt.Printf("\tPort %d/%s %s %s\n", port.ID, port.Protocol, port.State, port.Service.Name)
        }
    }

    fmt.Printf("Nmap done: %d hosts up scanned in %3f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed)
}
```

The program above outputs:

```bash
Host "172.217.16.46":
    Port 80/tcp open http
    Port 443/tcp open https
    Port 843/tcp filtered unknown
Host "31.13.81.36":
    Port 80/tcp open http
    Port 443/tcp open https
    Port 843/tcp open unknown
Host "216.58.215.110":
    Port 80/tcp open http
    Port 443/tcp open https
    Port 843/tcp filtered unknown
Nmap done: 3 hosts up scanned in 1.29 seconds
```
