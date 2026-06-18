# REST API 参考

English: [REST API Reference](../../en/reference/rest-api.md)

UModel 完整的 REST API 契约。所有接口接受并返回 `application/json`。

OpenAPI 3.1 规范：[api/openapi/openapi.yaml](../../../api/openapi/openapi.yaml)

---

## Base URL

```
http://localhost:8080
```

---

## 通用规则

### 请求

- 所有 `POST`/`PUT` 请求需要 `Content-Type: application/json`。
- JSON body 必须严格匹配 schema。**未知字段会被拒绝**，返回 `400 INVALID_ARGUMENT "invalid json body"`。
- 路径参数必须**非空**且严格匹配 URL 模板。缺少或多余的段会返回 `400` 或 `404`。

### 响应

所有错误响应使用统一的错误信封：

```json
{
  "error": {
    "code": "INVALID_ARGUMENT",
    "message": "人类可读的描述",
    "retryable": false
  }
}
```

### URL 格式

对于带 scope 的接口（`/api/v1/umodel/`、`/api/v1/entitystore/`、`/api/v1/samples/`、`/api/v1/query/`、`/api/v1/agent/`），URL 格式严格为：

```
/api/v1/{prefix}/{workspace}/{action}
```

**前缀之后只能有且仅有两个段**。不能有多余的尾部斜杠或额外段。

---

## 服务

### `GET /`

服务索引。

**响应 `200`：**

```json
{
  "service": "umodel-server",
  "status": "ok",
  "graphstore": { "provider": "memory", "status": "healthy" },
  "endpoints": {
    "health": "/healthz",
    "workspaces": "/api/v1/workspaces",
    "samples": "/api/v1/samples/{workspace}/multi-domain-quickstart:import",
    "query": "/api/v1/query/{workspace}/execute",
    "queryExplain": "/api/v1/query/{workspace}/explain",
    "agent": "/api/v1/agent/{workspace}/discover"
  }
}
```

### `GET /healthz`

健康检查。

**响应 `200`：**

```json
{
  "status": "ok",
  "graphstore": { "provider": "memory", "status": "healthy" }
}
```

---

## Workspace

### `POST /api/v1/workspaces`

创建 workspace。

**请求：**

```json
{
  "id": "my-workspace",
  "name": "My Workspace",
  "description": "可选描述",
  "labels": { "env": "production" }
}
```

| 字段 | 必填 | 类型 | 说明 |
|---|---|---|---|
| `id` | 是 | string | 正则：`^[a-z0-9](?:[a-z0-9_-]{0,62}[a-z0-9])?$` |
| `name` | 否 | string | |
| `description` | 否 | string | |
| `labels` | 否 | map[string]string | |
| `config` | 否 | map[string]map[string]any | |

**响应 `201`：** `WorkspaceMetadata`

### `GET /api/v1/workspaces`

列出 workspace。

**查询参数：**

| 参数 | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `page_size` | integer | 100 | 1–100 |
| `page_token` | string | | 分页令牌 |
| `include_deleted` | boolean | false | |
| `include_conflicts` | boolean | false | |

**响应 `200`：** `WorkspacePage`

### `GET /api/v1/workspaces/{workspace}`

获取单个 workspace。

**响应 `200`：** `WorkspaceMetadata`

### `PUT /api/v1/workspaces/{workspace}`

更新 workspace。

**请求：**

```json
{
  "name": "New Name",
  "labels": { "env": "staging" },
  "replace_labels": true,
  "if_match_version": 3
}
```

**响应 `200`：** `WorkspaceMetadata`

### `DELETE /api/v1/workspaces/{workspace}`

软删除 workspace。

**响应 `200`：** `WorkspaceMetadata`（status: `deleted`）

---

## EntityStore — 写入运行时数据

所有 EntityStore 接口仅使用 `POST` 方法。

### `POST /api/v1/entitystore/{workspace}/entities:write`

向 workspace 写入实体。

**请求：**

```json
{
  "workspace": "my-workspace",
  "idempotency_key": "可选的唯一键",
  "partial_success": false,
  "entities": [
    {
      "__domain__": "my-domain",
      "__entity_type__": "my.entity_type",
      "__entity_id__": "entity-001",
      "__category__": "entity",
      "__method__": "Update",
      "__first_observed_time__": 1700000000,
      "__last_observed_time__": 1800000000,
      "__keep_alive_seconds__": 3600,
      "display_name": "My Entity",
      "status": "active",
      "...自定义字段...": "..."
    }
  ]
}
```

**重要规则：**

| 规则 | 详情 |
|---|---|
| 顶层结构 | 必须是包含 `entities` 数组的对象。裸数组 `[{...}]` 会被拒绝，返回 `400 "invalid json body"`。 |
| 禁止未知顶层字段 | 顶层只允许 `workspace`、`idempotency_key`、`partial_success`、`entities` 这四个字段。 |
| `entities` | 必填。entity payload 对象的数组。 |
| `workspace` | body 中可选（服务端会用 URL 中的值覆盖）。建议包含以避免未知字段拒绝。 |

**Entity payload 字段：**

| 字段 | 必填 | 类型 | 说明 |
|---|---|---|---|
| `__domain__` | 是 | string | 匹配已导入 EntitySet 的 domain 名称 |
| `__entity_type__` | 是 | string | 匹配已导入 EntitySet 的 entity 类型 |
| `__entity_id__` | 是 | string | 唯一实体标识 |
| `__category__` | 否 | string | `"entity"`（标准值） |
| `__method__` | 否 | string | `"Update"`、`"Upsert"`、`"Delete"`、`"Expire"` |
| `__first_observed_time__` | 否 | integer | Unix 时间戳（秒） |
| `__last_observed_time__` | 否 | integer | Unix 时间戳（秒） |
| `__keep_alive_seconds__` | 否 | integer | TTL 提示 |
| *(自定义字段)* | 否 | any | EntitySet schema 中定义的任意字段 |

**响应 `200`：** `WriteResult`

```json
{
  "accepted": 9,
  "failed": 0,
  "items": [
    { "id": "my-domain/my.entity_type/entity-001", "ok": true }
  ]
}
```

### `POST /api/v1/entitystore/{workspace}/entities:expire`

通过 stable key 过期实体。

**请求：**

```json
{
  "workspace": "my-workspace",
  "ids": [
    "my-domain/my.entity_type/entity-001",
    "my-domain/my.entity_type/entity-002"
  ],
  "reason": "可选的原因说明"
}
```

**响应 `200`：** `WriteResult`

### `POST /api/v1/entitystore/{workspace}/relations:write`

向 workspace 写入关系。

**请求：**

```json
{
  "workspace": "my-workspace",
  "idempotency_key": "可选的唯一键",
  "partial_success": false,
  "relations": [
    {
      "__src_domain__": "my-domain",
      "__src_entity_type__": "my.service",
      "__src_entity_id__": "svc-001",
      "__relation_type__": "depends_on",
      "__dest_domain__": "my-domain",
      "__dest_entity_type__": "my.service",
      "__dest_entity_id__": "svc-002",
      "__category__": "relation",
      "__method__": "Update",
      "__first_observed_time__": 1700000000,
      "__last_observed_time__": 1800000000,
      "__keep_alive_seconds__": 3600
    }
  ]
}
```

**Relation payload 字段：**

| 字段 | 必填 | 类型 | 说明 |
|---|---|---|---|
| `__src_domain__` | 是 | string | 源实体 domain |
| `__src_entity_type__` | 是 | string | 源实体类型 |
| `__src_entity_id__` | 是 | string | 源实体 ID |
| `__relation_type__` | 是 | string | 关系类型，匹配已导入的模型 |
| `__dest_domain__` | 是 | string | 目标实体 domain |
| `__dest_entity_type__` | 是 | string | 目标实体类型 |
| `__dest_entity_id__` | 是 | string | 目标实体 ID |
| `__method__` | 否 | string | `"Update"`、`"Upsert"`、`"Delete"` |
| *(时间戳等)* | 否 | | 与 entity 时间戳字段相同 |

**响应 `200`：** `WriteResult`

### `POST /api/v1/entitystore/{workspace}/relations:expire`

通过 stable key 过期关系。

**请求：**

```json
{
  "workspace": "my-workspace",
  "ids": [
    "src-domain/src-type/src-id/relation-type/dest-domain/dest-type/dest-id"
  ],
  "reason": "可选的原因说明"
}
```

**响应 `200`：** `WriteResult`

---

## UModel — 模型管理

### `POST /api/v1/umodel/{workspace}/validate`

验证模型元素，不导入。

**请求：**

```json
{
  "workspace": "my-workspace",
  "elements": [
    { "kind": "EntitySet", "domain": "my-domain", "name": "my.service", "spec": { ... } }
  ]
}
```

**响应 `200`：** `ValidationResult`

### `POST /api/v1/umodel/{workspace}/import`

从本地路径导入模型定义。

**请求：**

```json
{
  "path": "examples/quickstart-multidomain"
}
```

**响应 `200`：** `UModelImportResult`

### `POST /api/v1/umodel/{workspace}/elements`

直接写入模型元素。

**请求：**

```json
{
  "workspace": "my-workspace",
  "elements": [
    { "kind": "EntitySet", "domain": "my-domain", "name": "my.service", "spec": { ... } }
  ]
}
```

**响应 `200`：** `WriteResult`

### `DELETE /api/v1/umodel/{workspace}/elements`

删除模型元素。

**请求：**

```json
{
  "ids": ["my-domain/my.service/EntitySet"]
}
```

**响应 `200`：** `WriteResult`

---

## Samples

### `POST /api/v1/samples/{workspace}/multi-domain-quickstart:import`

导入内置的多域 quickstart 样例。

**响应 `200`：** `SampleImportResult`

---

## Query

### `POST /api/v1/query/{workspace}/execute`

执行查询。

**请求：**

```json
{
  "query": ".entity with(domain='my-domain', name='my.service') | limit 20",
  "limit": 20,
  "timeout_ms": 5000,
  "time_range": {
    "from": "2024-01-01T00:00:00Z",
    "to": "2024-12-31T23:59:59Z"
  }
}
```

| 字段 | 必填 | 类型 | 说明 |
|---|---|---|---|
| `query` | 是 | string | 以 `.umodel`、`.entity` 或 `.topo` 开头的 SPL |
| `limit` | 否 | integer | 1–1000 |
| `timeout_ms` | 否 | integer | |
| `time_range` | 否 | TimeRange | `from`/`to` 为 ISO 8601 格式 |
| `format` | 否 | string | |
| `parameters` | 否 | map[string]any | |

**响应 `200`：** `QueryResult`

### `POST /api/v1/query/{workspace}/explain`

解释查询计划，不执行。

**请求：** 与 `/execute` 相同。

**响应 `200`：** `QueryExplain`

---

## Agent

### `GET /api/v1/agent/{workspace}/discover`

获取 Agent 发现元数据。

**响应 `200`：** `AgentDiscovery`

### `POST /api/v1/agent/{workspace}/tools:execute`

执行 Agent 工具。

**请求：**

```json
{
  "name": "tool-name",
  "arguments": { "key": "value" }
}
```

**响应 `200`：** `AgentToolCallResult`

### `POST /api/v1/agent/{workspace}/resources:read`

读取 Agent 资源。

**请求：**

```json
{
  "uri": "umodel://my-workspace/resources/some-resource"
}
```

**响应 `200`：** `AgentResourceReadResult`

---

## 常见错误码

| HTTP 状态码 | Code | 含义 |
|---|---|---|
| 400 | `INVALID_ARGUMENT` | URL 格式错误、JSON 格式错误、未知 JSON 字段、HTTP 方法错误 |
| 400 | `VALIDATION_FAILED` | Schema 校验失败 |
| 400 | `QUERY_PARSE_ERROR` | SPL 语法错误 |
| 404 | `NOT_FOUND` | 资源或 action 不存在 |
| 409 | `CONFLICT` | Workspace 冲突或版本不匹配 |
| 409 | `VERSION_CONFLICT` | 乐观锁失败 |
| 500 | `INTERNAL` | 服务端错误 |

---

## 外部集成快速检查清单

外部服务向 UModel 写入实体时，请确认：

1. **URL 格式**：`POST /api/v1/entitystore/{workspace}/entities:write`
   - 只有两个段：`{workspace}` 和 `entities:write`
   - 不能有尾部斜杠
   - 示例：`/api/v1/entitystore/otel_demo/entities:write`

2. **请求体**：包含 `entities` 数组的对象，不能是裸数组
   - ✅ `{"workspace": "otel_demo", "entities": [...]}`
   - ❌ `[{...}, {...}]`

3. **Entity 字段**：使用双下划线保留字段
   - ✅ `__domain__`、`__entity_type__`、`__entity_id__`
   - ❌ `domain`、`entity_type`、`entity_id`

4. **顶层禁止未知字段**
   - 顶层只允许 `workspace`、`idempotency_key`、`partial_success`、`entities`

5. **Content-Type**：`application/json`
