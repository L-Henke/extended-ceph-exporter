# extended-ceph-exporter

A Prometheus exporter to provide "extended" metrics about a Ceph cluster's running components (e.g., RGW).

Due to the closure of Koor Technologies, Inc. this repository has been made to continue the work on the extended-ceph-exporter project.

[![Ceph - RGW Bucket Usage Overview Grafana Dashboard Screenshot](grafana/ceph-rgw-bucket-usage-overview.png)](grafana/)

## Requirements

* Needs a Ceph cluster up and running.

* Needs an admin user

    ```
    radosgw-admin user create --uid extended-ceph-exporter --display-name "extended-ceph-exporter admin user" --caps "buckets=read;users=read;usage=read;metadata=read;zone=read"
    # Access key / "Username"
    radosgw-admin user info --uid extended-ceph-exporter | jq '.keys[0].access_key'
    # Secret key / "Password
    radosgw-admin user info --uid extended-ceph-exporter | jq '.keys[0].secret_key'
    ```

## Rook

If using Rook to manage RGW, the admin user may also be created using a `CephOjectStoreUser` resource:

```yaml
apiVersion: ceph.rook.io/v1
kind: CephObjectStoreUser
metadata:
  name: extended-ceph-exporter
  namespace: rook-ceph
spec:
  store: <objectstore-name>
  clusterNamespace: rook-ceph
  displayName: extended-ceph-exporter
  capabilities:
    buckets: read
    users: read
    usage: read
    metadata: read
    zone: read
```

Applying this will create an user with all permissions needed.

## Quickstart

* Clone the repository, download release binary or pull the container image:
  ```console
  git clone https://github.com/galexrt/extended-ceph-exporter
  cd extended-ceph-exporter
  ```

* Create a copy of the `.env.example` file and name it `.env`. Configure your RGW credentials and endpoint in the `.env` file.

* Configure Prometheus to collect metrics from the exporter from `:9138/metrics` endpoint using a static configuration, here's a sample scrape job from the `prometheus.yml`:

  ```yaml
  # For more information on Prometheus scrape_configs:
  # https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config
  scrape_configs:

    - job_name: "extended-ceph-metrics"

      # Override the global default and scrape targets from this job every 30 seconds.
      scrape_interval: 30s

      static_configs:
        # Please change the ip address `127.0.0.1` to target the exporter is running
        - targets: ['127.0.0.1:9138']
  ```

* To run the exporter locally, run `go run .`

* Should you have Grafana running for metrics visulization, check out the [Grafana dashboards](grafana/).

### Helm

To install the exporter to Kubernetes using Helm, check out the [extended-ceph-exporter Helm Chart](charts/extended-ceph-exporter/).

## Collectors

There is varying support for collectors. The tables
below list all existing collectors and the required Ceph components.

### Enabled by default

| Name             |                            Description                            | Ceph Component |
| :--------------- | :---------------------------------------------------------------: | -------------- |
| `rgw_buckets`    | Exposes RGW Bucket Usage and Quota metrics from the Ceph cluster. | RGW            |
| `rgw_user_quota` |       Exposes RGW User Quota metrics from the Ceph cluster.       | RGW            |

### Disabled by default

| Name          |                                                Description                                                 | Ceph Component |
| :------------ | :--------------------------------------------------------------------------------------------------------: | -------------- |
| `rbd_volumes` | Exposes RBD volumes size (volume pool, id, and name are available as labels). Not available at the moment. | RBD            |

## RGW: Multi-Realm Mode

You can use the exporter to scrape metrics from multiple RGW realms by enabling the "multi realm mode" and providing a "multi realm config" file.

An example multi realm config file can be found here [`realms.example.yaml`](realms.example.yaml).

Please note that if the multi realm mode is enabled, the RGW flags (e.g., `--rgw-host`, `--rgw-access-key`, `--rgw-secret-key`) are ignored as the `realms.yaml` (flag `--multi-realm-config`) takes over.

## Flags

```console
$ extended-ceph-exporter --help
Usage of exporter:
      --cache-duration duration      Cache duration in seconds (default 20s)
      --cache-enabled                Enable metrics caching to reduce load
      --collectors-enabled strings   List of enabled collectors (please refer to the readme for a list of all available collectors) (default [rgw_user_quota,rgw_buckets])
      --context-timeout duration     Context timeout for collecting metrics per collector (default 1m0s)
      --http-timeout duration        HTTP request timeout for collecting metrics for RGW API HTTP client (default 55s)
      --listen-host string           Exporter listen host (default ":9138")
      --log-level string             Set log level (default "INFO")
      --metrics-path string          Set the metrics endpoint path (default "/metrics")
      --multi-realm                  Enable multi realm mode (requires realms.yaml config, see --multi-realm-config flag)
      --multi-realm-config string    Path to your realms.yaml config file (default "realms.yaml")
      --rgw-access-key string        RGW Access Key
      --rgw-host string              RGW Host URL
      --rgw-secret-key string        RGW Secret Key
      --skip-tls-verify              Skip TLS cert verification
      --version                      Show version info and exit
pflag: help requested
exit status 2
```

## Development

### Requirements

* Golang 1.23.x (or higher should work)
* Ceph development files (`librados`, `librdb`)
    * If you are using `nix`, the `flake.nix` should be satisfy these lib dependencies.
* `helm`

### Making Changes to the Helm Chart

When changing anything in the Helm Chart, the version in the `Chart.yaml` needs to be increased according to [Semver](https://semver.org/).
Additionally `make helm-doc` must be run afterwards and the changes to the Helm Chart's `README.md` must be commited as well.

### Debugging

A VSCode debug config is available to run and debug the project.

To make the exporter talk with a Ceph RGW S3 endpoint, create a copy of the `.env.example` file and name it `.env`.
Be sure ot add your Ceph RGW S3 endpoint and credentials in it.
