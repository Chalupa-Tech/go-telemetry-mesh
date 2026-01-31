# Contributing

We welcome contributions! Please follow these steps to ensure a smooth process.

## 🛠️ Workflow

1.  **Fork & Branch**: Create a feature branch (`feat/usage-update`, `fix/login-bug`).
2.  **Commit**: Use clear, conventional commit messages.
3.  **Verify**: Run `make lint` and `make test` locally.
4.  **Pull Request**: Open a PR using the provided template.

## 🧪 Testing

All changes must include tests. Run `make test` to verify.

## 📝 Standards

-   **Code Style**: We use `gofumpt` and standard Go conventions.
-   **Linting**: We use `golangci-lint` (and `revive`).
-   **Changelog**: meaningful changes must be logged in `CHANGELOG.md`.
