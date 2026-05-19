# Link 与字段映射

English: [Links And Field Mappings](../../en/concepts/links-and-field-mappings.md)

Link 把独立模型定义连接成对象图。`fields_mapping` 定义链接两端的匹配方式。


## Link Kinds

| Kind | Source | Destination | 作用 |
|---|---|---|---|
| `data_link` | EntitySet | DataSet | 定义哪类遥测属于哪类对象。 |
| `entity_set_link` | EntitySet | EntitySet | 定义对象类型之间的拓扑关系语义。 |
| `storage_link` | DataSet | Storage | 将数据集路由到物理存储。 |
| `runbook_link` | EntitySet | RunbookSet | 将处置知识挂到对象类型上。 |

## DataLink

DataLink 将实体字段映射到数据集 label 或字段。

```yaml
kind: data_link
metadata:
  name: "devops.service_related_to_devops.metric.devops.service"
  domain: devops
spec:
  src:
    domain: devops
    kind: entity_set
    name: devops.service
  dest:
    domain: devops
    kind: metric_set
    name: devops.metric.devops.service
  data_link_type: related_to
  fields_mapping:
    "service_id": "service_id"
```

查询时语义：从 `devops.service` 实体出发，实体字段 `service_id` 匹配指标 label `service_id`。

## EntitySetLink

EntitySetLink 定义关系类型，例如：

- `contains`
- `calls`
- `instance_of`
- `parent_of`
- `same_as`

它定义“关系是什么意思”，运行时 Relation record 提供具体边。

## StorageLink

StorageLink 将物理存储与数据集语义分离。

```yaml
kind: storage_link
metadata:
  name: "devops.storage_link.devops.metric.devops.service"
  domain: devops
spec:
  src:
    domain: devops
    kind: metric_set
    name: devops.metric.devops.service
  dest:
    domain: devops
    kind: sls_metricstore
    name: devops.metric_set.core.storage
```

## Field Mapping Rules

同一个值在不同系统中使用不同名称时，字段映射负责对齐：

| Source field | Destination field | Example |
|---|---|---|
| Entity field | Dataset label | `service_id` -> `service_id` |
| Source entity field | Relation metric label | `${{src.service_id}}` -> `acs_arms_p_service_id` |
| Destination entity field | Relation metric label | `${{dest.service_id}}` -> `acs_arms_service_id` |

## 设计规则

- 优先映射稳定 ID，不要映射易变展示名。
- Link 保持小而明确。
- 条件启用放在 link-level filter。
- 命名中尽量体现 source、关系类型、destination。
- 将 link 文档化为 supported 之前，先用样例 entity 和 query plan 测试它。

## 相关概念

- [EntitySet](entity-sets.md)
- [DataSet](datasets.md)
- [Storage 与 GraphStore](storage-and-graphstore.md)
