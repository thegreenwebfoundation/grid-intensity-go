[![GoDoc](https://godoc.org/github.com/thegreenwebfoundation/grid-intensity-go?status.svg)](http://godoc.org/github.com/thegreenwebfoundation/grid-intensity-go) ![go-unit-test](https://github.com/thegreenwebfoundation/grid-intensity-go/workflows/go-unit-test/badge.svg)

# grid-intensity-go

A tool written in go, designed to be integrated into kubernetes, nomad, and other schedulers, to help you factor carbon intensity into decisions about where and when to run jobs.


## Usage


## Background

We know that the internet runs on electricity. That electricity comes from a mix of energy sources, including wind and solar, nuclear power, biomass, fossil gas, oil and coal and so on,

We call this the fuel mix, and this fuel mix can impact on the carbon intensity of your code.

## Move your code through time and space

Because the fuel mix will be different depending when and where you run your code, you can influence the carbon intensity of the code you write by moving it through time and space - either by making it run when the grid is greener, or making it run where it's greener, like a CDN running on green power.
