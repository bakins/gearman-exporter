gearman-exporter
================

Export [gearman](http://gearman.org/) metrics in [Prometheus](https://prometheus.io/) format.

See [Releases](https://github.com/bakins/gearman-exporter/releases) for pre-built binaries.

Build
=====

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

Running
=======

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

Metrics
=======

Metrics will be exposes on `/metrics`

```
curl http://localhost:9418/metrics

# HELP gearman_status_running number of running jobs
# TYPE gearman_status_running gauge
gearman_status_running{function="bar"} 0
gearman_status_running{function="baz"} 0
gearman_status_running{function="foo"} 0
# HELP gearman_status_total number of jobs in the queue
# TYPE gearman_status_total gauge
gearman_status_total{function="bar"} 0
gearman_status_total{function="baz"} 0
gearman_status_total{function="foo"} 0
# HELP gearman_status_workers number of number of capable workers
# TYPE gearman_status_workers gauge
gearman_status_workers{function="bar"} 1
gearman_status_workers{function="baz"} 1
gearman_status_workers{function="foo"} 1
# HELP gearman_up is gearman up
# TYPE gearman_up gauge
gearman_up 1
# HELP gearman_version_info gearman version
# TYPE gearman_version_info gauge
gearman_version_info{version="1.1.18"} 1
```

LICENSE
========

See [LICENSE](./LICENSE)
