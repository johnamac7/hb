## hb

[![GitHub release](https://img.shields.io/github/v/release/damianoneill/hb.svg)](https://GitHub.com/damianoneill/hb/releases/)

Healthbot Command Line Interface

### Synopsis

A tool for interacting with Healthbot over the REST API.

The intent with this tool is to provide bulk or aggregate functions, that
simplify interacting with Healthbot.

### Options

```
      --config string     config file (default is $HOME/.hb.yaml)
  -h, --help              help for hb
  -p, --password string   Healthbot Password (default "****")
  -r, --resource string   Healthbot Resource Name (default "localhost:8080")
  -t, --toggle            Help message for toggle
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

### Example

See below for a common set of example commands.

#### Devices

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

#### Device Groups

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

More complete examples can be viewed in the [tests folder](./cmd/provision/testdata/).
