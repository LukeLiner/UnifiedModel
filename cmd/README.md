# cmd/

Entry Layer — 二进制入口点。

| 目录 | 说明 | Spec |
|---|---|---|
| `umodel-server/` | UModel 服务端 | spec 09 |
| `umctl/` | CLI 工具 | spec 08 |
| `umodel-mcp/` | MCP Server | spec 07 |

## GraphStore Provider 参数

`umodel-server` 和 `umodel-mcp` 都支持 `--graphstore`：

| Provider | 说明 |
|---|---|
| `memory` | 纯内存，适合本地快速验证，进程退出后数据丢失；通过纯 Go 引擎支持 Ladybug 兼容只读 Cypher。 |
| `file.memory` | 内存查询 + JSON 文件持久化，默认按 workspace 和数据类型保存到 `<--data>/graphstore/file-memory/`；通过纯 Go 引擎支持 Ladybug 兼容只读 Cypher。 |
| `local.ladybug` | Ladybug-backed provider；开启 graph-match 和 Cypher 透传，真实实现需要 `-tags ladybug` 和本地 Ladybug 运行时。 |

示例：

```bash
go run -tags ladybug ./cmd/umodel-server --addr :8080 --data data --graphstore local.ladybug
go run -tags ladybug ./cmd/umodel-mcp --data data --graphstore local.ladybug
```

更多语义见 [GraphStore Providers](../docs/zh/graphstore-providers.md)。
