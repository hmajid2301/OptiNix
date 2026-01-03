# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 03-01-2026

### Added

- Logs file so users can share errors.
- Various TUI improvements such as progress bar.

### Changed

- Use flake to get options instead of channels.

## [0.1.4] - 08-01-2025

### Fixed

- Getting an error requiring a version when trying to fetch darwin options.

## [0.1.3] - 13-07-2024

### Fixed

- Duplicate options with same name in the database, it will now replace the existing option.


## [0.1.2] - 12-07-2024

### Changed

- Throw an error if user tries to pass no option name while using TUI.

## [0.1.1] - 11-07-2024

### Fixed

- Nix build not adding shell completions properly due to needing a DB setup.

## [0.1.0] - 09-07-2024

### Added

- Initial version released.

[unreleased]: https://gitlab.com/hmajid2301/optinix/compare/v0.2.0...HEAD
[0.2.0]: https://gitlab.com/hmajid2301/optinix/-/compare/v0.2.0...v0.1.4
[0.1.4]: https://gitlab.com/hmajid2301/optinix/-/compare/v0.1.4...v0.1.3
[0.1.3]: https://gitlab.com/hmajid2301/optinix/-/compare/v0.1.3...v0.1.2
[0.1.2]: https://gitlab.com/hmajid2301/optinix/-/compare/v0.1.2...v0.1.1
[0.1.1]: https://gitlab.com/hmajid2301/optinix/-/compare/v0.1.1...v0.1.0
[0.1.0]: https://gitlab.com/hmajid2301/optinix/releases/tag/v0.1.0
