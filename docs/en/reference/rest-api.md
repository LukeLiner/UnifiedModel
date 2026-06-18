# REST API Reference

中文：[REST API 参考](../../zh/reference/rest-api.md)

Complete REST API contract for UModel. All endpoints accept and return `application/json`.

OpenAPI 3.1 specification: [api/openapi/openapi.yaml](../../../api/openapi/openapi.yaml)

---

## Base URL

```
http://localhost:8080
```

---

## Common Rules

### Request

- All `POST`/`PUT` requests require `Content-Type: application/json`.
- JSON bodies must match the schema exactly. **Unknown fields are rejected** with `400 INVALID_ARGUMENT "invalid json body"`.
- Path parameters must be **non-empty** and match the exact URL template. Missing or extra segments return `400` or `404`.

### Response

All error responses use a stable error envelope:

```json
{
  "error": {
    "code": "INVALID_ARGUMENT",
    "message": "human-readable description",
    "retryable": false
  }
}
```

### URL Pattern

For scoped endpoints (`/api/v1/umodel/`, `/api/v1/entitystore/`, `/api/v1/samples/`, `/api/v1/query/`, `/api/v1/agent/`), the URL format is strictly:

```
/api/v1/{prefix}/{workspace}/{action}
```

**Two segments only** after the prefix. No trailing slash. No extra segments.

---

## Service

### `GET /`

Service index.

**Response `200`:**

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

Health check.

**Response `200`:**

```json
{
  "status": "ok",
  "graphstore": { "provider": "memory", "status": "healthy" }
}
```

---

## Workspaces

### `POST /api/v1/workspaces`

Create a workspace.

**Request:**

```json
{
  "id": "my-workspace",
  "name": "My Workspace",
  "description": "Optional description",
  "labels": { "env": "production" }
}
```

| Field | Required | Type | Notes |
|---|---|---|---|
| `id` | yes | string | Pattern: `^[a-z0-9](?:[a-z0-9_-]{0,62}[a-z0-9])?$` |
| `name` | no | string | |
| `description` | no | string | |
| `labels` | no | map[string]string | |
| `config` | no | map[string]map[string]any | |

**Response `201`:** `WorkspaceMetadata`

### `GET /api/v1/workspaces`

List workspaces.

**Query parameters:**

| Parameter | Type | Default | Notes |
|---|---|---|---|
| `page_size` | integer | 100 | 1–100 |
| `page_token` | string | | Pagination token |
| `include_deleted` | boolean | false | |
| `include_conflicts` | boolean | false | |

**Response `200`:** `WorkspacePage`

### `GET /api/v1/workspaces/{workspace}`

Get a workspace.

**Response `200`:** `WorkspaceMetadata`

### `PUT /api/v1/workspaces/{workspace}`

Update a workspace.

**Request:**

```json
{
  "name": "New Name",
  "labels": { "env": "staging" },
  "replace_labels": true,
  "if_match_version": 3
}
```

**Response `200`:** `WorkspaceMetadata`

### `DELETE /api/v1/workspaces/{workspace}`

Soft-delete a workspace.

**Response `200`:** `WorkspaceMetadata` (status: `deleted`)

---

## EntityStore — Writing Runtime Data

All EntityStore endpoints use `POST` only.

### `POST /api/v1/entitystore/{workspace}/entities:write`

Write entities to a workspace.

**Request:**

```json
{
  "workspace": "my-workspace",
  "idempotency_key": "optional-unique-key",
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
      "...custom fields...": "..."
    }
  ]
}
```

**Important rules:**

| Rule | Detail |
|---|---|
| Top-level structure | Must be an object with `entities` array. A bare JSON array `[{...}]` is **rejected** with `400 "invalid json body"`. |
| No unknown top-level fields | Only `workspace`, `idempotency_key`, `partial_success`, `entities` are allowed at the top level. |
| `entities` | Required. Array of entity payload objects. |
| `workspace` | Optional in body (server overwrites with URL value). Include it to avoid unknown-field rejection. |

**Entity payload fields:**

| Field | Required | Type | Notes |
|---|---|---|---|
| `__domain__` | yes | string | Domain name matching an imported EntitySet |
| `__entity_type__` | yes | string | Entity type matching an imported EntitySet |
| `__entity_id__` | yes | string | Unique entity identifier |
| `__category__` | no | string | `"entity"` (standard) |
| `__method__` | no | string | `"Update"`, `"Upsert"`, `"Delete"`, `"Expire"` |
| `__first_observed_time__` | no | integer | Unix timestamp (seconds) |
| `__last_observed_time__` | no | integer | Unix timestamp (seconds) |
| `__keep_alive_seconds__` | no | integer | TTL hint |
| *(custom fields)* | no | any | Arbitrary fields defined in EntitySet schema |

**Response `200`:** `WriteResult`

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

Expire entities by stable key.

**Request:**

```json
{
  "workspace": "my-workspace",
  "ids": [
    "my-domain/my.entity_type/entity-001",
    "my-domain/my.entity_type/entity-002"
  ],
  "reason": "optional reason"
}
```

**Response `200`:** `WriteResult`

### `POST /api/v1/entitystore/{workspace}/relations:write`

Write relations to a workspace.

**Request:**

```json
{
  "workspace": "my-workspace",
  "idempotency_key": "optional-unique-key",
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

**Relation payload fields:**

| Field | Required | Type | Notes |
|---|---|---|---|
| `__src_domain__` | yes | string | Source entity domain |
| `__src_entity_type__` | yes | string | Source entity type |
| `__src_entity_id__` | yes | string | Source entity ID |
| `__relation_type__` | yes | string | Relation type matching imported model |
| `__dest_domain__` | yes | string | Destination entity domain |
| `__dest_entity_type__` | yes | string | Destination entity type |
| `__dest_entity_id__` | yes | string | Destination entity ID |
| `__method__` | no | string | `"Update"`, `"Upsert"`, `"Delete"` |
| *(timestamps etc.)* | no | | Same as entity timestamps |

**Response `200`:** `WriteResult`

### `POST /api/v1/entitystore/{workspace}/relations:expire`

Expire relations by stable key.

**Request:**

```json
{
  "workspace": "my-workspace",
  "ids": [
    "src-domain/src-type/src-id/relation-type/dest-domain/dest-type/dest-id"
  ],
  "reason": "optional reason"
}
```

**Response `200`:** `WriteResult`

---

## UModel — Model Management

### `POST /api/v1/umodel/{workspace}/validate`

Validate model elements without importing.

**Request:**

```json
{
  "workspace": "my-workspace",
  "elements": [
    { "kind": "EntitySet", "domain": "my-domain", "name": "my.service", "spec": { ... } }
  ]
}
```

**Response `200`:** `ValidationResult`

### `POST /api/v1/umodel/{workspace}/import`

Import model definitions from a local path.

**Request:**

```json
{
  "path": "examples/quickstart-multidomain"
}
```

**Response `200`:** `UModelImportResult`

### `POST /api/v1/umodel/{workspace}/elements`

Write model elements directly.

**Request:**

```json
{
  "workspace": "my-workspace",
  "elements": [
    { "kind": "EntitySet", "domain": "my-domain", "name": "my.service", "spec": { ... } }
  ]
}
```

**Response `200`:** `WriteResult`

### `DELETE /api/v1/umodel/{workspace}/elements`

Delete model elements.

**Request:**

```json
{
  "ids": ["my-domain/my.service/EntitySet"]
}
```

**Response `200`:** `WriteResult`

---

## Samples

### `POST /api/v1/samples/{workspace}/multi-domain-quickstart:import`

Import the bundled multi-domain quickstart sample.

**Response `200`:** `SampleImportResult`

---

## Query

### `POST /api/v1/query/{workspace}/execute`

Execute a query.

**Request:**

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

| Field | Required | Type | Notes |
|---|---|---|---|
| `query` | yes | string | SPL starting with `.umodel`, `.entity`, or `.topo` |
| `limit` | no | integer | 1–1000 |
| `timeout_ms` | no | integer | |
| `time_range` | no | TimeRange | `from`/`to` as ISO 8601 |
| `format` | no | string | |
| `parameters` | no | map[string]any | |

**Response `200`:** `QueryResult`

### `POST /api/v1/query/{workspace}/explain`

Explain a query plan without executing.

**Request:** Same as `/execute`.

**Response `200`:** `QueryExplain`

---

## Agent

### `GET /api/v1/agent/{workspace}/discover`

Get agent discovery metadata.

**Response `200`:** `AgentDiscovery`

### `POST /api/v1/agent/{workspace}/tools:execute`

Execute an agent tool.

**Request:**

```json
{
  "name": "tool-name",
  "arguments": { "key": "value" }
}
```

**Response `200`:** `AgentToolCallResult`

### `POST /api/v1/agent/{workspace}/resources:read`

Read an agent resource.

**Request:**

```json
{
  "uri": "umodel://my-workspace/resources/some-resource"
}
```

**Response `200`:** `AgentResourceReadResult`

---

## Common Error Codes

| HTTP Status | Code | Meaning |
|---|---|---|
| 400 | `INVALID_ARGUMENT` | Bad URL, bad JSON, unknown JSON fields, wrong HTTP method |
| 400 | `VALIDATION_FAILED` | Schema validation failed |
| 400 | `QUERY_PARSE_ERROR` | Invalid SPL syntax |
| 404 | `NOT_FOUND` | Resource or action not found |
| 409 | `CONFLICT` | Workspace conflict or version mismatch |
| 409 | `VERSION_CONFLICT` | Optimistic lock failure |
| 500 | `INTERNAL` | Server error |

---

## Quick Integration Checklist

For external services writing entities to UModel:

1. **URL format**: `POST /api/v1/entitystore/{workspace}/entities:write`
   - Two segments only: `{workspace}` and `entities:write`
   - No trailing slash
   - Example: `/api/v1/entitystore/otel_demo/entities:write`

2. **Request body**: Object with `entities` array, not a bare array
   - ✅ `{"workspace": "otel_demo", "entities": [...]}`
   - ❌ `[{...}, {...}]`

3. **Entity fields**: Use double-underscore reserved fields
   - ✅ `__domain__`, `__entity_type__`, `__entity_id__`
   - ❌ `domain`, `entity_type`, `entity_id`

4. **No unknown top-level fields** in the JSON body
   - Only `workspace`, `idempotency_key`, `partial_success`, `entities` at the root

5. **Content-Type**: `application/json`
