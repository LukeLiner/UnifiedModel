# tests/

中文版本：[README.md](README.md)

The test suite is organized by responsibility.

| Directory | Purpose |
|---|---|
| `unit/` | Module-level unit tests. |
| `contract/` | Public contracts, storage interfaces, MCP, and SDK contract tests. |
| `integration/` | Cross-module integration paths. |
| `e2e/` | Quickstart end-to-end tests. |
| `golden/` | Query semantics, explain output, and documentation example golden tests. |
| `architecture/` | Dependency direction and forbidden API architecture guard tests. |

Common commands:

```bash
make guard
make test-service
```
