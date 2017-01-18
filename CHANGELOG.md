# Cameradar Changelog

This file lists all versions of the repository and precises all changes.

## v1.1.4

#### Minor changes :
* Simplified use of Docker image
* Renamed MySQL table name to be more explicit
* Refactoring of the Golang functional tester done
* The output was made more human readable
* Added automatic code quality checks for pull requests
* Added contribution documentation
* Updated dictionaries to add user suggestions for Chinese cameras
* Enhanced `result.json` file's format

#### Bugfixes :
* Fixed a bug in the functional testing in which if the `result.json` file was not formatted correctly, the test failed but was still considered a success.

## v1.1.3

#### Minor changes :
* Added automatic pushes to DockerHub to the travis integration
* Made travis configuration file better
* Changed the package generation scripts to make them report errors
* Removed old etix_rtsp_server binary from the test folder

#### Bugfixes :
* Fixed an issue that made it mandatory to launch tests at least once so that they can work the second time
* Fixed an issue that made the golang testing tool not compile in the testing script
* Fixed an issue that made the golang testing tool sometimes ignore some tests
* The previous known issue has been investigated and we don't know where it came from. However after a night of testing I have been unable to reproduce it, so I will consider it closed

## v1.1.2

#### Minor changes :
* Added travis integration
* Added default environment value for Docker deployment
* Updated docker image description with new easy usage
* Updated README badges style (replaced flat with square-flat)
* Build last package can now also generate a debug package if given the `Debug` command-line argument

#### Known issues :
* There is still the issue with Camera Emulation Server, see the [previous version's patchnote](#v1.1.1) for more information.

## v1.1.1

#### Minor changes :
* Removed unnecessary null pointer checks (thanks to https://github.com/elfring)
* Updated package description
* Removed debug message in CMake build
* Added `/ch01.264` to the URL dictionary in the deployment (Comelit default RTSP URL)
* Updated tests partially (still needs work to make the code cleaner)
  * Variable names are now compliant with Golang best practices
  * JSON variable names are back to normal
  * Functions have been moved in more appropriate source files
  * Structure definitions have been moved in more appropriate source files
  * Source files have been renamed to be more relevant
  * JUnit output now considers each camera as a test case
  * JUnit output now contains errors which makes debugging much easier
* Added header files where it was forgotten

#### Bugfixes :
* Fixed an issue where if you loose your internet connection during thumbnail generation, FFMpeg would get stuck forever and thus Cameradar would never finish
* Fixed an issue where multithreading could cause crashes
* Fixed an issue where the routes dictionary was mistaken for the credentials dictionary
* Fixed issues with the golang testing tool
  * Fixed automated camera generation
  * Fixed docker IP address resolution

#### Known issues :
* There is an issue with Camera Emulation Server that makes it impossible for Cameradar to generate thumbnails, which is why right now the verification of the thumbnails presence is commented and it is assumed correct. It is probably an issue with GST-RTSP-Server but requires investigation.

## v1.1.0

#### Major changes :
* There are more command line options
  * Port can now be overridden in the command line
  * Subnet can now be overridden in the command line
* Bruteforce is now multithreaded and will use as many threads as there are discovered cameras
* Thumbnail generation is now multithreaded and will use as many threads as there are discovered cameras
* There are now default configuration values in order to make cameradar easier to use

#### Minor changes :
* The algorithms take external input into account (so that a 3rd party can change the DB to help Cameradar in real-time) and thus check the persistent data at each iteration
* The default log level is now DEBUG instead of INFO
* The bruteforce logs are now INFO instead of DEBUG
* The thumbnail generation logs are now INFO instead of DEBUG

#### Bugs fixed
* Fixed a bug in which the MySQL cache manager would consider a camera with known ids as having a valid path even if it weren't
* Fixed a bug in which TCP RTSP streams would not generate thumbnails

## v1.0.5

* Fixed error in MySQL Cache Manager in which thumbnail generation on valid streams could not be done
* Fixed potential crash in the case the machine running cameradar has no memory left to allocate space for the dynamic cache manager

## v1.0.4

#### Bugs fixed :

* Fixed nmap package detection

## v1.0.3

#### Bugs fixed :

* Corrected GStreamer check

## v1.0.2

#### Bugs fixed :

* Fixed issues in MySQL Cache Manager

#### Minor changes :

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
