# grid-intensity-go

A tool written in go, designed to be integrated into kubernetes, nomad, and other schedulers, to help you factor carbon intensity into decisions about where and when to run jobs.

## Background

We know that the internet runs on electricity. That electricity comes from a mix of energy sources, including wind and solar, nuclear power, biomass, fossil gas, oil and coal and so on,

We call this the fuel mix, and this fuel mix can impact on the carbon intensity of your code.

## Move your code through time and space

Because the fuel mix will be different depending when and where you run your code, you can influence the carbon intensity of the code you write by moving it through time and space - either by making it run when the grid is greener, or making it run where it's greener, like a CDN running on green power.

## Installation

If you're on a mac you can install with homebrew/

```
brew install tgwf/grid-nntensity
```

Othewise you can download the binary from the releases tab here,for your architecture.

## How to use it

Call the `grid-intensity` binary to get an idea of the current carbon intensity of electricity, on the machine you're calling it on.

You'll recieve an estimated figure in terms of either 'low', 'medium' or 'high', valid for the next 30 minutes.

```
$ grid-intensity

### estimated carbon intensity for compute will be low on this machine for next 30 mins, until HH-MM:00
```

You can also use this to get an idea of the expected intensity of the next 24 hours. This allows your job scheduler to account for this in allocating jobs your cluster(s).

```
$ grid-intensity day

### estimated carbon intensity for compute will be low on this machine for next 24 hours, but peak between 7 and 10pm today

```

Append `--json` for a machine readable version of the same data, consumable by other services in your cluster.
