# Cameradar Changelog

This file lists all versions of the repository and precises all changes.

## v1.0.4

#### Minor changes :

* Fixed nmap package detection

## v1.0.3

#### Minor changes :

* Corrected GStreamer check

## v1.0.2

#### Minor changes :

* Fixed issues in MySQL Cache Manager
* Added useful debug logs

## v1.0.1

### Ubuntu 16.04 Release

#### Major changes :

* The Docker deployment is now done using Ubuntu 16.04 instead of Ubuntu 15.10, so that it uses more recent packages.

#### Minor changes :

* Removed useless dependencies

## v1.0.0

### First production-ready release

#### Major changes :

* Added functional testing

## v0.2.2

After doing some testing on a weirdly configured camera network in a far away Datacenter, I discovered that some Cameras needed a few tweaks to the Cameradar bruteforcing method in order to be accessed.

#### Major changes :

* Cameradar can access Cameras that are configured to always send 400 Bad Requests responses

#### Minor changes :

* Changed iterator name from `it` to `stream` in dumb cache manager to improve code readability

#### Bugfixes :

* Cameradar no longer considers a timing out Camera as an accessible stream

## v0.2.1

This package adds fixes the Docker deployment package.

#### Minor changes

* Fixed the Docker deployment package
* Updated README

## v0.2.0

### MySQL Cache Manager Release

This package adds a new cache manager using a MySQL database, that can store the results between mutiple uses.

#### Major changes

* Added a MySQL Cache Manager

#### Minor changes

* Removed legacy code
* Removed boost dependency
* Improved debugging logs

## v0.1.1

### Docker release

This package adds a way to deploy Cameradar using Docker.

#### Major changes

* Added a quick Docker deployment process
* Added automatic dependencies downloading through CMake for the manual installation
* Added CPack packaging for the Docker deployment

#### Minor changes

* Changed recommended cloning method to HTTPS
* Added lots of informations to README.md

## v0.1.0

This package was the first OpenSource version of Cameradar. It contained only a simple cache manager and had some bugs.
