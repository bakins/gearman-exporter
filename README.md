gearman-exporter
================

Export [gearman](http://gearman.org/) metrics in [Prometheus](https://prometheus.io/) format.

[![Build Status](https://img.shields.io/travis/bakins/gearman-exporter/master.svg)](https://travis-ci.org/bakins/gearman-exporter)
[![Docker Image](https://img.shields.io/docker/pulls/gearmanexporter/gearman-exporter.svg)](https://hub.docker.com/r/gearmanexporter/gearman-exporter)


Usage
=====

See [Releases](https://github.com/bakins/gearman-exporter/releases) for pre-built binaries.

Running
-------

```
./gearman-exporter --help
Gearman metrics exporter

Usage:
  gearman-exporter [flags]

Flags:
      --addr string       listen address for metrics handler (default "127.0.0.1:9418")
      --gearmand string   address of gearmand (default "127.0.0.1:4730")
```

When running, a simple healthcheck is availible on `/healthz`

Docker
------

A docker image is published from the Travis build to [Docker Hub](https://hub.docker.com/r/gearmanexporter/gearman-exporter).
```
docker run -p9418:9418 gearmanexporter/gearman-exporter --addr 0.0.0.0:9418
```

Metrics
-------

Metrics will be exposed on `/metrics`

```
curl http://localhost:9418/metrics

# HELP gearman_jobs number of jobs queued or running
# TYPE gearman_jobs gauge
gearman_jobs{function="bar"} 0
gearman_jobs{function="foo"} 0

# HELP gearman_jobs_running number of running jobs
# TYPE gearman_jobs_running gauge
gearman_jobs_running{function="bar"} 0
gearman_jobs_running{function="foo"} 0

# HELP gearman_jobs_waiting number of jobs waiting for an available worker
# TYPE gearman_jobs_waiting gauge
gearman_jobs_waiting{function="bar"} 0
gearman_jobs_waiting{function="foo"} 0

# HELP gearman_workers number of capable workers
# TYPE gearman_workers gauge
gearman_workers{function="bar"} 1
gearman_workers{function="foo"} 1

# HELP gearman_up is gearman up
# TYPE gearman_up gauge
gearman_up 1

# HELP gearman_version_info gearman version
# TYPE gearman_version_info gauge
gearman_version_info{version="1.1.18"} 1
```

Grafana
-------
https://grafana.com/grafana/dashboards/11423

Development
===========

Build
-----

Requires [Go](https://golang.org/doc/install). Tested with Go 1.8+.

Clone this repo into your `GOPATH` (`$HOME/go` by default) and run build:

```
mkdir -p $HOME/go/src/github.com/bakins
cd $HOME/go/src/github.com/bakins
git clone https://github.com/bakins/gearman-exporter
cd gearman-exporter
make build
```

You should then have two executables: gearman-exporter.linux.amd64 and gearman-exporter.darwin.amd64

You may want to rename for your local OS, ie `mv gearman-exporter.darwin.amd64 gearman-exporter`

Run Gearman server for testing
------------------------------
For testing you might want to run a gearman server with Docker:
```
docker run -p 4730:4730 cargomedia/gearman
```

While the server is running you could attach a worker function "foo" like this:
```
docker exec $(docker ps -qf ancestor=cargomedia/gearman) gearman -t9999 -wnf foo
```

Release new version
-------------------
1. Push a tag `vX.Y.Z` to Github
2. Travis will build the program and create a *Github release*


LICENSE
========

See [LICENSE](./LICENSE)
