# api/

English version: [README.md](README.md)

Adapter Layer：API 定义与协议契约。

| 目录 | 说明 |
|---|---|
| `openapi/` | OpenAPI 规范 |
| `mcp/` | MCP schema 定义 |

## GraphStore Provider 信号

REST API 调用方可以通过以下响应查看当前 GraphStore provider：

- `GET /`：`graphstore.provider` 和常见 endpoint hints。
- `GET /healthz`：`graphstore.provider`。
- `POST /api/v1/query/{workspace}/explain`：`provider` 和 `storage_provider`。
- `POST /api/v1/query/{workspace}/execute`：`explain.provider` 和 `explain.storage_provider`。

已知 provider：`memory`、`file.memory`、`local.ladybug`。更多信息见 [GraphStore Providers](../docs/zh/graphstore-providers.md)。
