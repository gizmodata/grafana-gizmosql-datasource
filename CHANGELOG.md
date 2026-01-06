# Changelog

All notable changes to the GizmoSQL Grafana Data Source Plugin will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.4.0] - 2025-01-06

### Fixed
- Include LICENSE file in plugin archive for Grafana submission compliance
- Include Go manifest (go.mod, go.sum) in plugin archive
- Version consistency between source files and release tags
- Updated dependencies to address security vulnerabilities

### Added
- Build provenance attestation for release artifacts
- CHANGELOG.md for release history

## [1.3.0] - 2025-01-06

### Added
- Screenshots for plugin catalog display
- GizmoSQL logo and branding assets
- Plugin catalog metadata for discoverability

### Changed
- Release workflow now dynamically injects version from git tags

## [1.2.0] - 2025-01-06

### Added
- Provisioning configurations for easier setup
- Sample dashboards demonstrating plugin capabilities

## [1.1.0] - 2025-01-06

### Added
- GitHub Actions CI/CD pipelines for automated builds and releases
- Multi-platform binary builds (Linux, macOS, Windows)
- SHA256 checksums for release artifacts

## [1.0.0] - 2025-01-06

### Added
- Initial release of GizmoSQL Grafana Data Source Plugin
- Go backend with Arrow Flight SQL connectivity
- TypeScript frontend with query editor
- TLS/SSL support for secure connections
- Username/password authentication
- Template variable support for dynamic dashboards
- Support for DECIMAL128/256 Arrow data types
- Alerting support
- Annotations support

[1.4.0]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.4.0
[1.3.0]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.3.0
[1.2.0]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.2.0
[1.1.0]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.1.0
[1.0.0]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.0.0
