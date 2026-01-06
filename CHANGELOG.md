# Changelog

All notable changes to the GizmoSQL Grafana Data Source Plugin will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.5.5] - 2026-01-06

### Fixed
- Use ts-loader with inlineSources for accurate source maps

## [1.5.4] - 2026-01-06

### Fixed
- Use resourcePath for source map paths (strips webpack:// prefix)

## [1.5.3] - 2026-01-06

### Fixed
- Remove custom devtoolModuleFilenameTemplate to use default source map paths

## [1.5.2] - 2026-01-06

### Fixed
- Fix LICENSE file being copied as directory instead of file
- Include CHANGELOG.md in plugin archive
- Switch back to swc-loader with proper configuration for source maps

## [1.5.1] - 2026-01-06

### Fixed
- Specify Go 1.22 in release workflow to match go.mod requirements

## [1.5.0] - 2026-01-06

### Changed
- Switch to official Grafana plugin build system (grafana/plugin-actions/build-plugin)
- Update Magefile.go to use Grafana SDK build targets

### Fixed
- Fix no-go-manifest error by using SDK build which generates go_plugin_build_manifest
- Fix js-map-no-match error by using official Grafana build tooling

## [1.4.9] - 2025-01-06

### Fixed
- Simplify source map paths using devtoolModuleFilenameTemplate

## [1.4.8] - 2025-01-06

### Fixed
- Configure ts-loader with sourceMap for proper TypeScript source maps

## [1.4.7] - 2025-01-06

### Fixed
- Switch from swc-loader to ts-loader for source map compatibility

## [1.4.6] - 2025-01-06

### Fixed
- Add trailing newlines to source files to match webpack output
- Remove go.mod toolchain directive for better compatibility

## [1.4.5] - 2025-01-06

### Fixed
- Fix source map content mismatch by using nosources-source-map

## [1.4.4] - 2025-01-06

### Fixed
- Update Go module path to match repository name (grafana-gizmosql-datasource)

## [1.4.3] - 2025-01-06

### Fixed
- Fix source map mismatch by removing dynamic date update from build workflow

## [1.4.2] - 2025-01-06

### Changed
- Update plugin ID to gizmodata-gizmosql-datasource

## [1.4.1] - 2025-01-06

### Fixed
- Fix release workflow LICENSE directory conflict

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

[1.5.5]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.5.5
[1.5.4]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.5.4
[1.5.3]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.5.3
[1.5.2]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.5.2
[1.5.1]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.5.1
[1.5.0]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.5.0
[1.4.9]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.4.9
[1.4.8]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.4.8
[1.4.7]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.4.7
[1.4.6]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.4.6
[1.4.5]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.4.5
[1.4.4]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.4.4
[1.4.3]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.4.3
[1.4.2]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.4.2
[1.4.1]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.4.1
[1.4.0]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.4.0
[1.3.0]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.3.0
[1.2.0]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.2.0
[1.1.0]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.1.0
[1.0.0]: https://github.com/gizmodata/grafana-gizmosql-datasource/releases/tag/v1.0.0
