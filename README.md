<!--
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2017 OpsVision Solutions

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->
# Snap-Telemetry Plugin for SignalFx [![Build Status](https://travis-ci.org/opsvision/snap-plugin-publisher-signalfx.svg?branch=master)](https://travis-ci.org/opsvision/snap-plugin-publisher-signalfx) [![Go Report Card](https://goreportcard.com/badge/github.com/opsvision/snap-plugin-publisher-signalfx)](https://goreportcard.com/report/github.com/opsvision/snap-plugin-publisher-signalfx)
Snap-Telemetry Plugin for SignalFx sends numeric values to [SignalFx](https://signalfx.com/).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Issues and Roadmap](#issues-and-roadmap)
3. [Acknowledgements](#acknowledgements)

## Getting Started
Read the system requirements, supported platforms, and installation guide for obtaining and using this Snap plugin.
### System Requirements 
* [golang 1.7+](https://golang.org/dl/) (needed only for building)

### Operating systems
All OSs currently supported by snap:
* Linux/amd64
* Darwin/amd64

### Installation
The following sections provide a guide for obtaining the Syslog collector plugin.

#### Download
The simplest approach is to use ```go get``` to fetch and build the plugin. The following command will place the binary in your ```$GOPATH/bin``` folder where you can load it into snap.
```
$ go get github.com/opsvision/snap-plugin-publisher-signalfx
```

#### Building
The following provides instructions for building the plugin yourself if you decided to downlaod the source. We assume you already have a $GOPATH setup for [golang development](https://golang.org/doc/code.html). The repository utilizes [glide](https://github.com/Masterminds/glide) for library management.
```
$ mkdir -p $GOPATH/src/github.com/opsvision
$ cd $GOPATH/src/github.com/opsvision
$ git clone http://github.com/opsvision/snap-plugin-publisher-signalfx
$ glide up
[INFO]	Downloading dependencies. Please wait...
[INFO]	--> Fetching updates for ...
[INFO]	Resolving imports
[INFO]	--> Fetching updates for ...
[INFO]	Downloading dependencies. Please wait...
[INFO]	Setting references for remaining imports
[INFO]	Exporting resolved dependencies...
[INFO]	--> Exporting ...
[INFO]	Replacing existing vendor dependencies
[INFO]	Project relies on ... dependencies.
$ go install
```

#### Source structure
The following file structure provides an overview of where the files exist in the source tree. The [signalfx.go](https://github.com/opsvision/snap-plugin-publisher-signalfx/blob/master/signalfx/signalfx.go) file does all the work.
```
snap-plugin-collector-syslog
├── glide.yaml
├── LICENSE
├── main.go
├── README.md
└── signalfx
    ├── signalfx.go
    └── signalfx_test.go
```

### Configuration and Usage
Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)

#### Load the Plugin
Once the framework is up and running, you can load the plugin.
```
$ snaptel plugin load snap-plugin-publisher-signalfx
```

#### Task File
TBD
Once the task file has been created, you can create and watch the task.
```
$ snaptel task create -t signalfx.yaml
$ snaptel task list
ID                                       NAME                                         STATE     ...
f3ad05b2-3706-4991-ab29-c96e15813893     Task-f3ad05b2-3706-4991-ab29-c96e15813893    Running   ...
$ snaptel task watch f3ad05b2-3706-4991-ab29-c96e15813893
```

## Documentation
TBD

### Issues and Roadmap
* **Testing:** The testing being done is rudimentary at best. Need to improve the testing.

_Note: Please let me know if you find a bug or have feedbck on how to improve the collector._

## Acknowledgements
* Author: [@dishmael](https://github.com/dishmael/)
* Company: [OpsVision Solutions](https://github.com/opsvision)
