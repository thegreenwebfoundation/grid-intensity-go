# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

> **Added** for new features.
> **Changed** for changes in existing functionality.
> **Deprecated** for soon-to-be removed features.
> **Removed** for now removed features.
> **Fixed** for any bug fixes.
> **Security** in case of vulnerabilities.

## Unreleased

### Added

- Push container image to scaleway registry.

## 0.4.1 2023-05-05

- Fix: Add `auth-token` header to Electricity Maps request.

## 0.4.0 2022-08-31

### Added

- Add node and region labels to prometheus metrics for carbon aware scheduling.

### Changed

- Breaking change to refactor interface to return the same JSON format for all
providers.
- Breaking change to rename region to location since region has a different meaning
in Nomad and Kubernetes.

### Fixed

- Fix link to install.sh in readme.

## 0.3.0 2022-07-15

### Added

- Add WattTime provider.
- Add Carbon Intensity Org UK support to CLI.

## 0.2.1 2022-07-01

### Fixed

- Generate installation token for GoReleaser.

## 0.2.0 2022-06-23

### Added

- Add install script.
- Add Ember support to exporter.

## 0.1.0 2022-06-21

### Added

- Initial release.
