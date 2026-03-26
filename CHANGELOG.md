# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `node_name` field in PingResponse proto — peers now identify themselves by K8s node name.

### Fixed
- Removed unused GOPRIVATE/GITHUB_PAT auth block from Dockerfile that caused build-push failure (git not in Alpine image).
- `target_node` metric label now uses node name from gRPC response instead of pod IP.

## [0.1.0] - 2026-03-25

### Added
- gRPC mesh probe server with DNS-based peer discovery via headless Service.
- HTTP server: `/metrics` (Prometheus), `/healthz`, `/readyz`.
- Prometheus metrics: probe duration histogram, probe/error counters, peer gauge, cluster health gauge.
- OpenTelemetry gRPC instrumentation.
- `build-push.yml` workflow for Gitea container image CD via chalupa-infra reusable workflow.
- Dockerfile with multi-arch buildx support and GOPRIVATE config.

### Changed
- `release.yml` updated to template standard: module tidy verification, test gate, release notes extraction.

### Fixed
- Initial scaffold created from `go-repo-template`.
- Renamed from `gemini-mesh` to `go-telemetry-mesh` per ADR 019 update.
