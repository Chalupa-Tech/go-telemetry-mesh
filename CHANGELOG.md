# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Release workflow synced to standard template: chained build-push + deploy job, removed separate build-push.yml
- CI: align git insteadOf config with peer repos (cosmetic, no behavior change)

## [0.2.1] - 2026-03-26

### Fixed
- Release workflow now correctly extracts semver from CHANGELOG, skipping `[Unreleased]`
- Dockerfile: add git, GOPRIVATE, and GITHUB_PAT for private module fetch during container build

## [0.2.0] - 2026-03-25

### Added
- `node_name` field in PingResponse proto — peers now identify themselves by K8s node name.
- OpenTelemetry SDK initialization via `go-telemetry` v0.2.0 — links gRPC traces to log lines with `trace_id`/`span_id`
- Migrated all slog calls in mesh client to context-aware variants for trace correlation

### Fixed
- Removed unused GOPRIVATE/GITHUB_PAT auth block from Dockerfile that caused build-push failure (git not in Alpine image).
- `target_node` metric label now uses node name from gRPC response instead of pod IP.

### Changed
- Upgraded Go to 1.25, OTel dependencies to v1.42.0
- Dockerfile builder image updated to golang:1.25-alpine

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
