Use this dockerfile to start a testing SNMP v3 server

```
$ docker build .
[...]
Successfully built <image id>
$ docker run -p 10161:10161/udp <image id>
```
