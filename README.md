# Otilio

[![Build Status](https://travis-ci.org/isalgueiro/otilio.svg?branch=master)](https://travis-ci.org/isalgueiro/otilio)

Welcome to Otilio. [Beat](https://www.elastic.co/products/beats) to query SNMP data. This was built following [Beats developer guide](https://www.elastic.co/guide/en/beats/devguide/6.1/new-beat.html) and uses [gosnmp](https://github.com/soniah/gosnmp) library.

Example setup (see `otilo.yml`):

```
otilio:
  # Defines how often an event is sent to the output
  period: 1s

  # SNMP host to query
  host: "127.0.0.1"

  # SMNP version
  version: 2c

  # SNMP community
  community: "public"

  # oids to query
  # (the starting dot is intended)
  oids:
    - {oid: ".1.3.6.1.2.1.1.1.0", name: sysDescr}
    - {oid: ".1.3.6.1.2.1.1.3.0", name: sysUpTime}
```
This will get oids `1.3.6.1.2.1.1.1.0` and `1.3.6.1.2.1.1.3.0` from SNMP server at localhost and store them in `otilio-YYYY.MM.DD` index in Elasticsearch in fields `sysDescr` and `sysUpTime`.

## Building

Ensure that this folder is at the following location:
`${GOPATH}/github.com/isalgueiro/otilio`

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Init Project
To get running with Otilio and also install the
dependencies, run the following command:

```
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Otilio in the git repository, run the following commands:

```
git remote set-url origin https://github.com/isalgueiro/otilio
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).
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
Please check index settings, as they may not fit your use case.

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
