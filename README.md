# vmware_exporter

VMWare Exporter is implemented via the multi-target exporter pattern.
By multi-target exporter pattern we refer to a specific design, in which:

- the exporter will get the target’s metrics via a network protocol.
- the exporter does not have to run on the machine the metrics are taken from.
- the exporter gets the targets and a query config string as parameters of Prometheus’ GET request.
- the exporter subsequently starts the scrape after getting Prometheus’ GET requests and once it is done with scraping.

When the exporter starts the scarape, it is performing following actions:

- POST /api/session
- GET /api/vcenter/datacenter?datacenters={someDatacenterID}
- GET /api/vcenter/datastore?datacenters={someDatacenterID}
- GET /api/vcenter/host?datacenters={someDatacenterID}
- DELETE POST /api/session
- 
The exporter exposes following metrics:

- vmware_datastore_capacity_size
- vmware_datastore_free_size
- vmware_host_connection_state
- vmware_host_power_state

The exporter filters all queries on the datacenter ID, which MUST be transferred via the query param "dc"

## Getting Started
The project is developed in Go (1.23+).\
The repository is formatted for use in GoLand.

NOTE: The rest of this README assumes you are using GoLand.

## Prerequisites
Development requirements:
* GoLand.

## How to start
* Install GoLand .
* Open GoLand - Clone  project from Github

## Run System
* make start
* Open a web browser and navigate to `http://localhost:9141/probe?target=some.vcenterhost.com&dc=datacenter-666`

## Build System
* make build

## Push System to repository
* make deploy


## Enviroment
    HOST                (default binds to 0.0.0.0)
    PORT                (listening port, default 9141)
    VCENTER_USER
    VCENTER_PASS

## Prometheus configuration

```yaml
  - job_name: 'vmware_exporter'
    metrics_path: /probe
    params:
      dc: [datacenter-666]
    static_configs:
      - targets:
        - some.vcenterhost.com
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:9141
```
