[![GoDoc](https://godoc.org/github.com/thegreenwebfoundation/grid-intensity-go?status.svg)](http://godoc.org/github.com/thegreenwebfoundation/grid-intensity-go) ![go-unit-test](https://github.com/thegreenwebfoundation/grid-intensity-go/workflows/go-unit-test/badge.svg) ![docker](https://github.com/thegreenwebfoundation/grid-intensity-go/workflows/docker-integration-test/badge.svg) ![kubernetes](https://github.com/thegreenwebfoundation/grid-intensity-go/workflows/kubernetes-integration-test/badge.svg) ![nomad](https://github.com/thegreenwebfoundation/grid-intensity-go/workflows/nomad-integration-test/badge.svg)

# grid-intensity-go

A tool written in Go, designed to be integrated into Kubernetes, Nomad, and other schedulers, to help you factor carbon intensity into decisions about where and when to run jobs.

The tool has 3 components.

- The `grid-intensity` CLI for interacting with carbon intensity data.
- A [Prometheus](https://prometheus.io/) exporter with carbon intensity metrics that can be deployed via
Docker, Nomad, or Kubernetes.
- A Go library that can be integrated into your Go code.

## Changelog

See [CHANGELOG.md](/CHANGELOG.md).

## Background

We know that the internet runs on electricity. That electricity comes from a mix of energy sources, including wind and solar, nuclear power, biomass, fossil gas, oil and coal and so on,

We call this the fuel mix, and this fuel mix can impact on the carbon intensity of your code.

## Move your code through time and space

Because the fuel mix will be different depending when and where you run your code, you can influence the carbon intensity of the code you write by moving it through time and space - either by making it run when the grid is greener, or making it run where it's greener, like a CDN running on green power.

## Inspired By

This tool builds on research and tools developed from across the sustainable software community. 

### Articles

- A carbon aware internet - Branch magazine - https://branch.climateaction.tech/issues/issue-2/a-carbon-aware-internet/
- Carbon Aware Kubernetes - https://devblogs.microsoft.com/sustainable-software/carbon-aware-kubernetes/
- Clean energy technologies threaten to overwhelm the grid. Here’s how it can adapt. - https://www.vox.com/energy-and-environment/2018/11/30/17868620/renewable-energy-power-grid-architecture

### Papers

- A Tale of Two Visions: Designing a Decentralized Transactive Electric System - https://ieeexplore.ieee.org/document/7452738
- Carbon Explorer - https://github.com/facebookresearch/CarbonExplorer/
- Cucumber: Renewable-Aware Admission Control for Delay-Tolerant Cloud and Edge Workloads - https://arxiv.org/abs/2205.02895 
- Let's Wait Awhile: How Temporal Workload Shifting Can Reduce Carbon Emissions in the Cloud - https://arxiv.org/abs/2110.13234

### Tools

- Carbon Aware Nomad - experimental branch - https://github.com/hashicorp/nomad/blob/h-carbon-meta/CARBON.md
- Cloud Carbon Footprint - https://www.cloudcarbonfootprint.org/
- Scaphandre - https://github.com/hubblo-org/scaphandre
- Solar Protocol - http://solarprotocol.net/
- The carbon aware scheduler - https://pypi.org/project/carbon-aware-scheduler/

## Installing

- Install via [brew](https://brew.sh/).

```sh
brew install thegreenwebfoundation/carbon-aware-tools/grid-intensity
```

- Install via curl (feel free to do due diligence and check the [script](https://github.com/thegreenwebfoundation/grid-intensity-go/blob/main/install.sh) first).

```sh
curl -fsSL https://raw.githubusercontent.com/thegreenwebfoundation/grid-intensity-go/main/install.sh | sudo sh 
```

- Fetch a binary release from the [releases](https://github.com/thegreenwebfoundation/grid-intensity-go/releases) page.

## grid-intensity CLI

The CLI allows you to interact with carbon intensity data from multiple providers.

```sh
$ grid-intensity
Provider ember-climate.org needs an ISO country code as a location parameter.
ESP detected from your locale.
ESP
[
	{
		"emissions_type": "average",
		"metric_type": "absolute",
		"provider": "Ember",
		"location": "ESP",
		"units": "gCO2e per kWh",
		"valid_from": "2021-01-01T00:00:00Z",
		"valid_to": "2021-12-31T23:59:00Z",
		"value": 193.737
	}
]
```

The `--provider` and `--location` flags allow you to select other providers and locations.
You can also set the `GRID_INTENSITY_PROVIDER` and `GRID_INTENSITY_LOCATION` environment
variables or edit the config file at `~/.config/grid-intensity/config.yaml`.

```sh
$ grid-intensity --provider CarbonIntensityOrgUK --location UK
{
	"from": "2022-07-14T14:30Z",
	"to": "2022-07-14T15:00Z",
	"intensity": {
		"forecast": 184,
		"actual": 194,
		"index": "moderate"
	}
}
```

The [providers](#providers) section shows how to configure other providers.

## grid-intensity exporter

The `exporter` subcommand starts the prometheus exporter on port 8000.

```sh
$ grid-intensity exporter --provider Ember --location FR
Using provider "Ember" with location "FR"
Metrics available at :8000/metrics
```

View the metrics with curl.

```
$ curl -s http://localhost:8000/metrics | grep grid
# HELP grid_intensity_carbon_average Average carbon intensity for the electricity grid in this location.
# TYPE grid_intensity_carbon_average gauge
grid_intensity_carbon_average{provider="Ember",location="FR",units="gCO2 per kWh"} 67.781 1718258400000
```

**Note about Prometheus and samples in the past**

If you are using the exporter with the ElectricityMaps provider, it will return a value for estimated, which will be the most recent one, and another value for the real value, which can be a few hours in the past. Depending on your Prometheus installation, it could be that the metrics that have a timestamp in the past are not accepted, with an error such as this:

`Error on ingesting samples that are too old or are too far into the future`

In that case, you can configure the property `tsdb.outOfOrderTimeWindow` to extend the time window accepted, for example to `3h`.


### Docker Image

Build the docker image to deploy the exporter.

```sh
CGO_ENABLED=0 GOOS=linux go build -o grid-intensity .
docker build -t thegreenwebfoundation/grid-intensity:latest .
```

### Kubernetes

Install the [helm](https://helm.sh/) chart in [/helm/grid-intensity-exporter](https://github.com/thegreenwebfoundation/grid-intensity-go/tree/main/helm/grid-intensity-exporter).
Needs the Docker image to be available in the cluster.

```sh
helm install --set gridIntensity.location=FR grid-intensity-exporter helm/grid-intensity-exporter
```

### Nomad

Edit the Nomad job in [/nomad/grid-intensity-exporter.nomad](https://github.com/thegreenwebfoundation/grid-intensity-go/blob/main/nomad/grid-intensity-exporter.nomad) to set the
env vars `GRID_INTENSITY_LOCATION` and `GRID_INTENSITY_PROVIDER`

Start the Nomad job. Needs the Docker image to be available in the cluster.

```sh
nomad run ./nomad/grid-intensity-exporter.nomad
```

## grid-intensity-go library

See the [/examples/](https://github.com/thegreenwebfoundation/grid-intensity-go/tree/main/examples) 
directory for examples of how to integrate each provider.

## Providers

Currently these providers of carbon intensity data are integrated. If you would like
us to integrate more providers please open an [issue](https://github.com/thegreenwebfoundation/grid-intensity-go/issues).

### Electricity Maps

[Electricity Maps](https://app.electricitymaps.com/map) have carbon intensity data
from multiple sources. You need to get an API token and URL from their
[API portal](https://api-portal.electricitymaps.com/) to use the API. You can use
their free tier for non-commercial use or sign up for a 30 day trial.

The `location` parameter needs to be set to a zone present in the public [zones](https://static.electricitymaps.com/api/docs/index.html#zones) endpoint.

```sh
ELECTRICITY_MAPS_API_TOKEN=your-token \
ELECTRICITY_MAPS_API_URL=https://api-access.electricitymaps.com/free-tier/ \
grid-intensity --provider=ElectricityMaps --location=IN-KA
```

### WattTime

[WattTime](https://www.watttime.org/) have carbon intensity data from multiple sources.
You need to [register](https://www.watttime.org/api-documentation/#authentication) to use the API.

The `location` parameter should be set to a supported location. The `/ba-from-loc`
endpoint allows you to provide a latitude and longitude. See the [docs](https://www.watttime.org/api-documentation/#determine-grid-region) for more details.

```sh
WATT_TIME_USER=your-user \
WATT_TIME_PASSWORD=your-password \
grid-intensity --provider=WattTime --location=CAISO_NORTH
```

### Ember

Carbon intensity data from [Ember](https://ember-climate.org/), is embedded in the binary
in accordance with their licensing - [CC-BY-SA 4.0](https://ember-climate.org/creative-commons/)

```sh
grid-intensity --provider=Ember --location=DE
```

The `location` parameter should be set to a 2 or 3 char ISO country code.

### UK Carbon Intensity API

UK Carbon Intensity API https://carbonintensity.org.uk/ this is a public API
and the only location supported is `UK`.

```sh
grid-intensity --provider=CarbonIntensityOrgUK --location=UK
```
