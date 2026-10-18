# Changelog

All notable changes to this project are documented in this file.

## [Unreleased]

## [v1.0.3] - 2026-10-18

### Fixed

- `scripts/check-tools.sh` install hints no longer embed URLs. Curator rc.7 refused the package as a profile
  member (`audit.capability.network-undeclared` for `git-scm.com` and `go.dev`) although the script only prints
  the hints and makes no network request. No behaviour change and no capability declaration.

## [v1.0.2] - 2026-07-14

### Added

- Added a cross-platform setup and credentials reference for reusable skill installation patterns.

### Changed

- Refreshed the README with clearer Relux stack framing and ecosystem background.

### Notes

- Top-level `v*` tags describe the skill repository release.
- `tuitestkit/*` tags continue to version the Go library separately.
