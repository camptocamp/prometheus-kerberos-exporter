Prometheus Kerberos exporter
============================

## Usage

```
Usage:
  prometheus-kerberos-exporter [OPTIONS]

Application Options:
  -v, --version          Show version.
  -u, --username=        A username to use to connect to kerberos server. [$KERBEROS_USER]
  -r, --realm=           A realm to use to connect to kerberos server. [$KERBEROS_REALM]
  -k, --keytab=          A keytab file to use to connect to kerberos server. [$KERBEROS_KEYTAB_FILE]
  -s, --server=          A list of servers to connect to. (separated by commas) [$KERBEROS_SERVERS]
      --scrape-interval= Duration between two scrapes. (default: 5s) [$KERBEROS_SCRAPE_INTERVAL]
      --listen-address=  Address to listen on for web interface and telemetry. (default: 0.0.0.0:9259) [$KERBEROS_LISTEN_ADDRESS]
      --metric-path=     Path under which to expose metrics. (default: /metrics) [$KERBEROS_METRIC_PATH]
      --verbose          Enable debug mode [$KERBEROS_VERBOSE]

Help Options:
  -h, --help             Show this help message
```

## Metrics

```
# HELP kerberos_exporter_build_info Kerberos exporter build informations
# TYPE kerberos_exporter_build_info gauge
kerberos_exporter_build_info{build_date="2019-11-18",commit_sha="XXXXXXXXXX",golang_version="go1.12.7",version="1.0.0"} 1
# HELP kerberos_status_available Kerberos server availability
# TYPE kerberos_status_available gauge
kerberos_status_available{kdc="kdc1.example.com"} 1
kerberos_status_available{kdc="kdc2.example.com"} 1
kerberos_status_available{kdc="kdc3.example.com"} 0
```
