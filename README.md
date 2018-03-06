# Otilio

[![Build Status](https://travis-ci.org/isalgueiro/otilio.svg?branch=master)](https://travis-ci.org/isalgueiro/otilio) [![Go Report Card](https://goreportcard.com/badge/github.com/isalgueiro/otilio)](https://goreportcard.com/report/github.com/isalgueiro/otilio)

Welcome to Otilio. [Beat](https://www.elastic.co/products/beats) to query SNMP data. This was built following [Beats developer guide](https://www.elastic.co/guide/en/beats/devguide/6.1/new-beat.html) and uses [gosnmp](https://github.com/soniah/gosnmp) library.

Example setup (see [otilo.yml](otilio.yml)):

```
otilio:
  # Defines how often an event is sent to the output
  period: 1s

  # SNMP hosts to query
  hosts: ["192.168.1.1", "192.168.1.2"]

  # SMNP version: 1, 2c or 3
  version: 2c

  # SNMP community
  community: "public"

  # oids to query
  # (the starting dot is intended)
  oids:
    - {oid: ".1.3.6.1.2.1.1.1.0", name: sysDescr}
    - {oid: ".1.3.6.1.2.1.1.3.0", name: sysUpTime}
```
This will get oids `1.3.6.1.2.1.1.1.0` and `1.3.6.1.2.1.1.3.0` from SNMP servers at 192.168.1.1 and 192.168.1.2 and store them in `otilio-YYYY.MM.DD` index in Elasticsearch in fields `sysDescr` and `sysUpTime`.

SNMP V3 configuration example

```
otilio:
  # Defines how often an event is sent to the output
  period: 1s

  # SNMP host to query
  hosts: ["127.0.0.1"]
  port: 10161

  # SMNP version
  version: 3

  # SNMP user security model parameters
  # currently only SHA auth and DES encryption supported ¯\_(ツ)_/¯
  user: "theuser"
  authpass: "theauthpassword"
  privpass: "theprivacyencryptionpassword"

  # oids to query
  # (the starting dot is intended)
  oids:
    - {oid: ".1.3.6.1.2.1.25.1", name: hrSystem}
```

## Building

Ensure that this folder is at the following location:
`${GOPATH}/github.com/isalgueiro/otilio`

### Requirements

* [Golang](https://golang.org/dl/) 1.9

### Build

To build the binary for Otilio run the command below. This will generate a binary
in the same directory with the name otilio.

```
make
```


### Run

To run Otilio with debugging output enabled, run:

```
./otilio -c otilio.yml -e -d "*"
```


### Test

To test Otilio, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/otilio.template.json and etc/otilio.asciidoc

```
make update
```
Please **check index settings**, as they may not fit your use case.

### Cleanup

To clean  Otilio source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```

## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The hole process to finish can take several minutes.
