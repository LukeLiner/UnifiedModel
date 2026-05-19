# api/

Adapter Layer — API 定义与协议契约。

| 目录 | 说明 | Spec |
|---|---|---|
| `openapi/` | OpenAPI 规范 | spec 08 |
| `mcp/` | MCP schema 定义 | spec 07 |

## GraphStore Provider Signals

REST API callers can inspect the active GraphStore provider through:

- `GET /`: `graphstore.provider` plus common endpoint hints
- `GET /healthz`: `graphstore.provider`
- `POST /api/v1/query/{workspace}/explain`: `provider` and `storage_provider`
- `POST /api/v1/query/{workspace}/execute`: `explain.provider` and `explain.storage_provider`

Known provider values are `memory`, `file.memory`, and `local.ladybug`. See
[GraphStore Providers](../docs/en/graphstore-providers.md) for persistence behavior
and startup examples.
