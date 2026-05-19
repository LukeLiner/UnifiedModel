# Links And Field Mappings

中文：[Link 与字段映射](../../zh/concepts/links-and-field-mappings.md)

Links turn independent model definitions into an object graph. Field mappings explain how one side of a link matches another side.


## Link Kinds

| Kind | Source | Destination | Purpose |
|---|---|---|---|
| `data_link` | EntitySet | DataSet | Shows which telemetry belongs to which object type. |
| `entity_set_link` | EntitySet | EntitySet | Defines topology relation semantics between object types. |
| `storage_link` | DataSet | Storage | Routes datasets to physical storage. |
| `runbook_link` | EntitySet | RunbookSet | Attaches operational guidance to an object type. |

## DataLink

DataLink maps entity fields to dataset labels or fields.

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

Query-time meaning: from a `devops.service` entity, entity field `service_id` matches metric label `service_id`.

## EntitySetLink

EntitySetLink defines topology relation types such as:

- `contains`
- `calls`
- `instance_of`
- `parent_of`
- `same_as`

The link defines what a relation means. Runtime relation records provide the actual edges.

## StorageLink

StorageLink keeps physical storage separate from dataset semantics.

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

Use mappings when a value has different names in different systems:

| Source field | Destination field | Example |
|---|---|---|
| Entity field | Dataset label | `service_id` -> `acs_arms_service_id` |
| Source entity field | Relation metric label | `${{src.service_id}}` -> `acs_arms_p_service_id` |
| Destination entity field | Relation metric label | `${{dest.service_id}}` -> `acs_arms_service_id` |

## Design Rules

- Keep links small and explicit.
- Prefer stable IDs over display names.
- Put conditional availability in link-level filters when needed.
- Name links so source, relation type, and destination are obvious.
- Test links with a sample entity and a query plan before documenting them as supported.

## Related Concepts

- [Entity Sets](entity-sets.md)
- [Datasets](datasets.md)
- [Storage And GraphStore Providers](storage-and-graphstore.md)
