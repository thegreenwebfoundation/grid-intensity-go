# Overview

For the Ember provider we embed their dataset in the grid-intensity CLI. The data
uses ISO 3 char country codes. However we want to also support 2 char ISO codes.

Ideally this will be added by Ember in a future release of their data. Until then
we can use a simple Go program in the `hack` directory to map the country codes.

## Processing the data

- From the root of this repo call the program.

```
go run hack/country_codes.go ember-input.csv > ember-output.csv
```

- Update the data file in the `ember` directory. e.g. co2-intensities-ember-2021.csv
- If the data has already been processed an error will be returned. 

```
go run hack/country_codes.go /tmp/ember-output.csv
panic: data already processed - `country_code_iso_2` should not be present
```
