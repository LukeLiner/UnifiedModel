# Query Service 指南

English: [Query Service Guide](../../en/guides/query-service.md)

Query Service 是 UModel 定义、实体、关系和拓扑的唯一公共读取路径。它接受以 `.umodel`、`.entity` 或 `.topo` 开头的 SPL 字符串。


## 为什么读取统一走 Query Service

UModel 不暴露分散的公共读取 API，例如 entity lookup、relation lookup、graph traversal 或 model search endpoint。统一读取面让 CLI、Web UI、REST API、MCP tools 和 SDK 保持一致。

## 入口

REST：

```http
POST /api/v1/query/{workspace}/execute
POST /api/v1/query/{workspace}/explain
```

CLI：

```bash
go run ./cmd/umctl --addr http://localhost:8080 query run demo ".umodel | limit 5"
go run ./cmd/umctl --addr http://localhost:8080 query explain demo ".umodel | limit 5"
```

Agent tool：

```bash
go run ./cmd/umctl --addr http://localhost:8080 agent tool demo query_spl_execute '{"query":".umodel | limit 5"}'
```

## `.umodel`

读取 UModel 定义：

```bash
go run ./cmd/umctl --addr http://localhost:8080 query run demo ".umodel with(kind='entity_set') | sort name | limit 20"
```

常见读取：

- 列出 EntitySet。
- 查看 metric、log、trace、event、storage、link 定义。
- 支撑 Web UI Explorer 的图/表视图。

## `.entity`

读取运行时实体：

```bash
go run ./cmd/umctl --addr http://localhost:8080 query run demo ".entity with(domain='devops', name='devops.service', query='checkout') | project __entity_id__,display_name | limit 20"
```

## `.topo`

读取运行时拓扑关系：

```bash
go run ./cmd/umctl --addr http://localhost:8080 query run demo ".topo | graph-call getDirectRelations([(:\"devops@devops.service\" {__entity_id__: '10000000000000000000000000000101'})]) | project src,relation,dest | limit 20"
```

## 常用管道操作

- `with(...)`：source-specific 过滤。
- `project`：选择字段。
- `sort`：排序。
- `limit`：限制输出。
- `graph-call`：拓扑函数。

查看内置示例：

```bash
go run ./cmd/umctl --addr http://localhost:8080 query examples
```

## Explain

```bash
go run ./cmd/umctl --addr http://localhost:8080 query explain demo ".entity with(domain='devops', name='devops.service') | limit 5"
```

Explain 输出包含 source、provider、storage provider、filters 和 limits。

## 边界规则

- 不新增 Query Service 之外的公共 entity/relation/topology 读取 endpoint。
- CLI 领域读取保持在 `query run` 和 `query explain` 后面。
- AgentGateway resources 保持 metadata-only，运行时 rows 通过 tools 返回。
