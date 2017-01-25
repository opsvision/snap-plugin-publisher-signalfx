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
2. [Documentation](#documentation)
  * [Configuration and Usage](#configuration-and-usage)
  * [Publisher Output](#publisher-output)
3. [Issues and Roadmap](#issues-and-roadmap)
4. [Acknowledgements](#acknowledgements)

## Getting Started
Read the system requirements, supported platforms, and installation guide for obtaining and using this Snap plugin.
### System Requirements 
* [golang 1.7+](https://golang.org/dl/) (needed only for building)

### Operating systems
All OSs currently supported by snap:
* Linux/amd64
* Darwin/amd64

### Installation
The following sections provide a guide for obtaining the plugin.

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
The following file structure provides an overview of where the files exist in the source tree.

```
snap-plugin-publisher-signalfx
├── glide.yaml
├── LICENSE
├── main.go
├── metadata.yml
├── README.md
├── scripts
│   ├── load.sh
│   └── unload.sh
├── signalfx
│   └── signalfx.go
└── tasks
    └── signalfx.yaml
```

## Documentation

### Configuration and Usage
Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)

#### Load the Plugin
Once the framework is up and running, you can load the plugin.
```
$ snaptel plugin load snap-plugin-publisher-signalfx
Plugin loaded
Name: signalfx
Version: 1
Type: publisher
Signed: false
Loaded Time: Tue, 24 Jan 2017 20:45:48 UTC
```

#### Task File
You need to create or update a task file to use the SignalFx publisher plugin. We have provided an example, _tasks/awssqs.yaml_ shown below. In our example, we utilize the psutil collector so we have some data to work with. There are three (3) configuration settings you can use.

Setting|Description|Required?|
|-------|-----------|---------|
|debug_file|An absolute path to a log file - this makes debugging easier.|No|
|hostname|The hostname to use; if absent, the plugin will attempt to determine the hostname.|No|
|token|The SignalFx [API token](https://developers.signalfx.com).|Yes|


```
---
  version: 1
  schedule:
    type: "simple"
    interval: "5s"
  max-failures: 10
  workflow:
    collect:
      config:
      metrics:
        /intel/psutil/load/load1: {} 
        /intel/psutil/load/load15: {}
        /intel/psutil/load/load5: {}
        /intel/psutil/vm/available: {}
        /intel/psutil/vm/free: {}
        /intel/psutil/vm/used: {}
      publish:
        - plugin_name: "signalfx"
          config:
            token: "1234ABCD"
            debug_file: "/tmp/signalfx-debug.log"
            hostname: "spiderman"
```

Once the task file has been created, you can create and watch the task.
```
$ snaptel task create -t signalfx.yaml
Plugin loaded
Name: signalfx
Version: 1
Type: publisher
Signed: false
Loaded Time: Tue, 24 Jan 2017 20:45:48 UTC

$ snaptel task list
ID                                       NAME                                         STATE     ...
f3ad05b2-3706-4991-ab29-c96e15813893     Task-f3ad05b2-3706-4991-ab29-c96e15813893    Running   ...
```

_Note: Truncated results for brevity._

### Publisher Output
The SignalFx plugin **will only publish numeric values (int64 and float64)** using the SignalFx [Gauge and GaugeF](https://github.com/signalfx/golib/tree/master/sfxclient) respectively.  The code attempts to convert numeric values; e.g. uint --> int64.  All other metric values will be ignored (e.g. strings).  The metrics will be sent with the namespace, metric value (converted), and the hostname as a dimension. This makes it simple to identify and use the incoming values in SignalFx.

## Issues and Roadmap
* **Testing:** The testing being done is rudimentary at best. Need to improve the testing.

_Note: Please let me know if you find a bug or have feedbck on how to improve the collector._

## Acknowledgements
* Author: [@dishmael](https://github.com/dishmael/)
* Company: [OpsVision Solutions](https://github.com/opsvision)
