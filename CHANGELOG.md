# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com)
and this project adheres to [Semantic Versioning](https://semver.org).

---

## [0.0.2] - 2025-06-30


### Added
- Added basic status messages to improve CLI feedback.
- Implemented a new output text processing system with centralized message handling.
- Implemented quiet to disable banners on console

### Changed
- Changed Path for temp Files from `src/koyaneframework/tmp` to /tmp/koyane-framework


### Fixed
- Fixed a bug where the `analyze` command printed incorrect word lengths.
- Fixed an issue in mask-based wordlist generation that caused duplicate entries.

