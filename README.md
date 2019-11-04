## hb

[![GitHub release](https://img.shields.io/github/v/release/damianoneill/hb.svg)](https://GitHub.com/damianoneill/hb/releases/)

Healthbot Command Line Interface

## Synopsis

A tool for interacting with Healthbot over the REST API.

The intent with this tool is to provide bulk or aggregate functions, that
simplify interacting with Healthbot.

## Options

```
      --config string     config file (default is $HOME/.hb.yaml)
      --debug             Enable REST debugging
  -h, --help              help for hb
  -p, --password string   Healthbot Password (default "****")
  -r, --resource string   Healthbot Resource Name (default "localhost:8080")
  -u, --username string   Healthbot Username (default "admin")
```

A full list of the options available with the tool is described in the [docs](./docs/hb.md).

The config file (~/.hb.yaml) mentioned in the option above allows persistent setting of the following options:

```yaml
---
resource: "hb-server:8080"
username: root
password: changeme
```

## Examples

See below for a common set of example commands.

### Summary

Provides an overview of the Healthbot Installation

```sh
$ hb summary -r foyle:8080
Using config file: /Users/doneill/.hb.yaml

Healthbot Version: HealthBot 2.1.0-beta
Healthbot Time: 2019-11-01T18:41:59Z

No of Managed Devices: 12

  Device Id  Platform    Release                           Serial Number

  MX34
  Ex10       EX4200-48P  15.1R5.5                          BQ0208189292
  Ex11       EX4200-48P  15.1R5.5                          BQ0210090233
  vMX-210    VMX         17.4R2.4                          VM5C4052633B
  vMX-212    VMX         17.4R2.4                          VM5C4053F72F
  flames     MX240       19.4I-20191005.0.2307             JN1263A8BAFC
  Capella38  ACX6360-OR  19.2I-20190228_dev_common.0.2316  DX008
  Capella39  ACX6360-OR  19.2I-20190228_dev_common.0.2316  DX004
  Mx35       MX240       19.3R1.8                          JN11AC665AFC
  Mx34       MX240       17.4R1.16                         JN1261DB3AFC
  mx960-1    MX960       19.3R1.8                          JN1232C39AFA
  mx960-3    MX960       19.3R1.8                          JN1233EF1AFA

No of Device Groups: 6

  Device Group    No of Devices

  Cappella_Group              2
  Real_Mx_Group               2
  ptp-test-group              2
  TT-SNMP-Group               1
  Switch_Group                2
  Test-Group                  3
```

### Scaffold

The scaffold command will read the configuration from a Healthbot installation and create the config directories and learned configuration. The example below assumes your in the directory where the config should be written too and that a valid .hb.yaml exists for the Healthbot installation you want to learn from.

```sh
hb --config .hb.yaml scaffold  .
```

```console
$ tree .
.
├── device-groups
│   └── device-groups.yml
├── devices
    └── devices.yml


2 directories, 2 files
```

### Devices

The example below will generate a request against hb-server to provision device defined in yml or json files in the /tmp/devices directory.

```sh
hb provision -r hb-server:8080 -u root -p changeme devices -d /tmp/devices/
```

An example of a configuration (/tmp/devices/mx.yml) can be seen below.

```yaml
---
device:
  - device-id: mx960-1
    host: 172.30.177.102
    authentication:
      password:
        username: doneill
        password: "$9$.mQ3EhrvMX0BIcrlLXGDjkfT369"
  - device-id: mx960-3
    host: 172.30.177.113
```

To delete the Devices using the configuration, you can pass the '-e' flag.

```sh
 hb provision -r hb-server:8080 -u root -p changeme devices -d /tmp/devices/mx.yml -e
```

> You can only delete Devices that are not associated with any Device Groups.

### Device Groups

The example below will generate a request against the HB Server with Username and Password defined in ~/.hb.yaml to provision Device Groups defined in yml or json files in the /tmp/device-groups directory.

```sh
hb provision device-groups -d /tmp/device-groups/
```

An example of a configuration (/tmp/device-groups/l2-swtiches.yml) can be seen below.

```yaml
---
device-group:
  - device-group-name: l2-test-group
    devices:
      - mx960-1
      - mx960-3
    authentication:
      password:
        username: root
        password: "$9$VgY2akqfTQnGDPQFnpuevWLxd"
```

#### Helper Files

The example below will generate a request against the HB Server with Username and Password defined in the local .hb.yaml to upload files in the /tmp/helper-files directory.

```sh
hb --config .hb.yaml provision helper-files -d /tmp/helper-files/
Using config file: .hb.yaml
Using directory: /tmp/helper-files/
Using files: [bps.py]
Successfully uploaded 1 Files
```

More complete examples can be viewed in the [tests folder](./cmd/provision/testdata/).

## TODO

- Commands
  - ~~version~~ - verison of hb tool
  - ~~completion~~ - bash completion for hb
  - ~~Summary~~ - high level info on the Healthbot installation
  - Provision
    - ~~Devices~~
    - ~~DeviceGroups~~
    - ~~Helper Files~~
    - Playbook Instances
    - All
  - Scaffold - generate hb configuration from an existing Healthbot deployment (round trip)
- Refactor common code across commands
- UT
- Move types into their own package
