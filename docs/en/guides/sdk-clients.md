# SDK And Client Guide

中文：[SDK 与客户端指南](../../zh/guides/sdk-clients.md)

Public client surfaces: REST, CLI, MCP, generated model SDKs, and the minimal Go REST client.


## REST API

REST fits integrations that require explicit HTTP control.

Contract:

- [OpenAPI](../../../api/openapi/openapi.yaml)

Common endpoints:

```http
POST /api/v1/workspaces
POST /api/v1/umodel/{workspace}/import
POST /api/v1/entitystore/{workspace}/entities:write
POST /api/v1/query/{workspace}/execute
GET  /api/v1/agent/{workspace}/discover
```

## CLI

`umctl` covers local development, examples, and documentation workflows.

```bash
go run ./cmd/umctl --addr http://localhost:8080 workspace create demo '{"name":"Demo"}'
go run ./cmd/umctl --addr http://localhost:8080 query run demo ".umodel | limit 5"
```

Reference: [CLI Reference](../reference/cli.md).

## Go REST Client

The minimal Go REST client lives under [sdk/go/service](../../../sdk/go/service). It wraps public REST contracts used by the quick start.

Example shape:

```go
client := service.NewClient("http://localhost:8080")
ctx := context.Background()

_, err := client.CreateWorkspace(ctx, service.CreateWorkspaceRequest{
    ID:   "demo",
    Name: "Demo",
})
if err != nil {
    return err
}

result, err := client.Query(ctx, "demo", service.QueryRequest{
    Query: ".umodel | limit 5",
})
if err != nil {
    return err
}
_ = result
```

Run its tests:

```bash
cd sdk/go
go test ./service
```

## Generated Model SDKs

Generated model SDKs represent UModel schema types:

- Go: [sdk/go/umodel](../../../sdk/go/umodel/README.en.md)
- Python: [sdk/python/umodel](../../../sdk/python/umodel/README.md)
- Java: [generated/java](../../../generated/java)

Regenerate and verify:

```bash
make expand
make verify
```

Runnable examples:

- [SDK Examples](../../../examples/sdk/README.md)
- [SDK End-To-End Flow](../../../examples/sdk/full-flow.md)
- Go model inspector: [examples/sdk/go/model-inspector](../../../examples/sdk/go/model-inspector)
- Go REST quick start: [examples/sdk/go/service-quickstart](../../../examples/sdk/go/service-quickstart)
- Python model inspector: [examples/sdk/python/inspect_model_pack.py](../../../examples/sdk/python/inspect_model_pack.py)

## Generated SDK Integration Examples

Generated SDK workflow: parse a model pack, fail fast on unknown kinds or missing envelope fields, inspect metadata, and hand the verified pack to the REST API or CLI. Runtime entity writes and graph queries stay behind the public service contracts.

### Parse And Validate A Model File In Go

Repository-local examples under `sdk/go` use the `umodel_go_cli` module path. External projects replace that import path with the published or vendored SDK module path.

```go
package main

import (
    "fmt"
    "os"

    umodel "umodel_go_cli/umodel"
)

func main() {
    data, err := os.ReadFile("examples/quickstart-multidomain/devops/entity_set/devops.service.yaml")
    if err != nil {
        panic(err)
    }

    obj, err := umodel.ParseYamlUModel(data)
    if err != nil {
        panic(err)
    }
    if err := obj.Validate(); err != nil {
        panic(err)
    }

    metadata := obj.GetMetadata()
    fmt.Printf("%s %s/%s\n", obj.GetKind(), metadata.Domain, metadata.Name)
}
```

### Scan A Model Pack In Python

Run from the repository root with `PYTHONPATH=sdk/python`, or package the generated Python SDK in your application environment.

```python
from pathlib import Path

from umodel import (
    get_link_endpoints,
    get_object_metadata,
    is_link_object,
    parse_umodel_yaml,
)

for path in Path("examples/quickstart-multidomain").rglob("*.yaml"):
    obj = parse_umodel_yaml(path.read_bytes())
    metadata = get_object_metadata(obj)

    if is_link_object(obj):
        src, dest = get_link_endpoints(obj)
        if src and dest:
            print(f"{path}: {obj.get_kind()} {src.name} -> {dest.name}")
        else:
            print(f"{path}: {obj.get_kind()} {metadata['domain']}/{metadata['name']}")
    else:
        print(f"{path}: {obj.get_kind()} {metadata['domain']}/{metadata['name']}")
```

### Combine Generated Types With The REST Client

Generated model SDKs handle local validation and metadata inspection. The service client handles workspace import and runtime queries.

```go
client := service.NewClient("http://localhost:8080")
ctx := context.Background()

_, err := client.CreateWorkspace(ctx, service.CreateWorkspaceRequest{
    ID:   "demo",
    Name: "Demo",
})
if err != nil {
    return err
}

_, err = client.ImportUModel(ctx, "demo", service.UModelImportRequest{
    Path: "examples/quickstart-multidomain",
})
if err != nil {
    return err
}

result, err := client.Query(ctx, "demo", service.QueryRequest{
    Query: ".umodel with(kind='entity_set') | limit 5",
})
if err != nil {
    return err
}
_ = result
```

Integration patterns:

| Scenario | SDK role | Service role |
|---|---|---|
| CI for model packs | Parse, validate, and list model metadata | Optional import smoke test |
| Platform startup | Load bundled model definitions and fail fast on incompatible schema | Import into the target workspace |
| Developer tooling | Convert JSON/YAML, inspect kinds, and show source/destination links | Query live workspace state |
| Agent applications | Read static model metadata for prompts or tool descriptions | Use Query Service or MCP for runtime graph context |

## MCP

MCP serves agent clients that need UModel discovery, resources, query examples, or query tools.

```bash
go run ./cmd/umodel-mcp --data data --graphstore file.memory
```

Use Streamable HTTP when the client expects an HTTP MCP server:

```bash
go run ./cmd/umodel-mcp --transport http --addr 127.0.0.1:8090 --data data --graphstore file.memory
```

The MCP envelope stays JSON-RPC. Tool/resource text payloads use TOON (`text/toon`), and tool calls also return `structuredContent` for JSON-oriented clients.

Reference: [MCP Reference](../reference/mcp.md). Examples: [examples/mcp](../../../examples/mcp/README.md).

## Choosing A Surface

| Need | Use |
|---|---|
| Human local workflow | CLI or Web UI |
| HTTP integration | REST API or Go REST client |
| Generated model construction | Generated SDKs |
| Agent integration | MCP / AgentGateway |
| Contract inspection | OpenAPI and MCP schema |
