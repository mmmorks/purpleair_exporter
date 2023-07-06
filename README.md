# `purpleair_exporter` [![source](https://img.shields.io/badge/source-gray?logo=github)](https://github.com/willglynn/purpleair_exporter)

This is a Prometheus/OpenMetrics exporter for [PurpleAir air quality sensors](https://www.purpleair.com).

This tool runs an HTTP server which responds to requests for `/purpleair?target=…` by making HTTP request(s) to the
target PurpleAir sensor over your local network. There are no API keys needed, and it works fully offline. This
strategy is similar to the [SNMP exporter](https://github.com/prometheus/snmp_exporter), essentially acting as a
Prometheus-to-PurpleAir proxy. One instance of this exporter can easily support multiple sensors even on a
resource-constrained host.

## Quickstart

Container images are available at [Docker Hub](https://hub.docker.com/r/willglynn/purpleair_exporter) and [GitHub 
container registry](https://github.com/willglynn/purpleair_exporter/pkgs/container/purpleair_exporter). 

```shell
$ docker run -it --rm -p 2020:2020 willglynn/purpleair_exporter
# or
$ docker run -it --rm -p 2020:2020 ghcr.io/willglynn/purpleair_exporter
level=info msg="Starting HTTP server" addr=:2020
```

Once it's running, fetch
[http://localhost:2020/purpleair?target=0.0.0.0](http://localhost:2020/purpleair?target=0.0.0.0), replacing `0.0.0.0`
with the IP address of the sensor on your LAN.

## Prometheus configuration

Scrape one or more sensor target(s) via an instance of `purpleair_exporter`:

```yaml
scrape_configs:
  - job_name: 'purpleair'
    metrics_path: /purpleair
    scrape_interval: 10s    # 1s if you dare
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:2020  # this exporter's hostname:port
    static_configs:
      - targets:
          - 172.16.4.40   # address of a sensor
          - 172.16.4.41
          - 172.16.4.42
```

## Exporter configuration

Minimal, via environment variables:

* `LISTEN`: a host and port on which the web server should bind
* `PORT`: a port on which the web server should bind

If neither of these are set, the exporter runs on `:2020`. There is no authentication for this service, no sensor
configuration, and no TLS, just like the PurpleAir sensors themselves.

## Metrics endpoint

`GET /purpleair` supports the following URL parameters:

* `target`: the IP address of the sensor (required)
* `period`: `1s` if you only want 1-second readings, `2m` if you only want 2-minute averages, empty if you want both

## Status

This works for me and my [PurpleAir Flex](https://www2.purpleair.com/products/purpleair-flex) setup. Feel free to open
pull requests with proposed changes.
