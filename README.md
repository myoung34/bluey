Bluey
=====

[![Coverage Status](https://coveralls.io/repos/github/myoung34/bluey/badge.svg)](https://coveralls.io/github/myoung34/bluey)
[![Docker Pulls](https://img.shields.io/docker/pulls/myoung34/bluey.svg)](https://hub.docker.com/r/myoung34/bluey)


## Supported Emitters ##

* stdout
* Http(s) Webhooks
* datadog (dogstatsd)
* sqlite

## Usage ##

### Generate Config ###

```
$ cat <<EOF >config.toml
[general]
sleep_interval = 2 # defaults to 1
logging_level = DEBUG # defaults to INFO
logfile = /var/log/foo.log # defaults to stdout

# stdout example
[stdout]

# Generic application/json example
[webhook]
enabled      = true
device_uuids = "61020304620663030064086507660367=puck"
method       = "POST"
url          = "http://www.foo.com"
headers      = "{\"Content-Type\": \"application/json\"}"
template     = "{\"color\": \"{{ .Nickname }}\", \"gravity\": {{ .Minor }}, \"mac\": \"{{ .Mac }}\", \"temp\": {{ .Major }}, \"timestamp\": \"{{ .Timestamp }}\"}"

# Note: make sure that the dd agent has DD_DOGSTATSD_NON_LOCAL_TRAFFIC=true
[datadog]
enabled     = true
statsd_host = "statsdhost.corp.com"
statsd_port = 8125
tags        = "color:{{.Nickname}},mac:{{.Mac}}"
gauges      = "bluey.temperature={{.Major}},bluey.gravity={{.Minor}}"

[sqlite]
enabled      = true
device_uuids = "1234=foo,61020304620663030064086507660367=puck,5678=bar"
file         = "/home/myoung/repos/github/bluey/bluey.db"

[mqtt]
enabled      = true
device_uuids = "1234=foo,61020304620663030064086507660367=puck,5678=bar"
url          = "tcp://mosquitto.service.kube:1883"
client_id    = "foo"
username     = "wut"
topic        = "bar"
template     = "{\"color\": \"{{ .Nickname }}\", \"gravity\": {{ .Minor }}, \"mac\": \"{{ .Mac }}\", \"temp\": {{ .Major }}, \"timestamp\": \"{{ .Timestamp }}\"}"
```

### Run ###

```
$ bluey
$ # Or from docker ( generate config into $cwd/config/config.toml )
$ docker run -it \
  -v $(pwd)/config:/etc/bluey \
  --privileged \
  --net=host \
  myoung34/bluey:latest \
  -c /etc/bluey/config.toml
```
