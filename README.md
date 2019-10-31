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

### Example

The example below will generate a request against hb-server to provision device defined in yml or json files in the /tmp/devices directory.

```sh
hb provision -r hb-server:8080 -u root -p changeme devices -d /tmp/devices/
```

An example of a configuration (mx.yml) can be seen below.

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
    authentication:
      password:
        username: doneill
        password: "$9$3eJsnAuvMX-dsSreWXxwsmf5z/CuO1"
```
