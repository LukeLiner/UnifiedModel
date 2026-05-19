# Java Test Summary

中文版本：[测试总结.md](测试总结.md)

This document is the English companion for the generated Java SDK test summary.

## Purpose

The Java test summary records the verification status of generated Java SDK assets, including type generation, compilation, serialization behavior, and compatibility checks.

## How To Regenerate

From the repository root:

```bash
make expand
make verify-java
```

## Maintenance Notes

- Treat generated Java files as derived artifacts.
- Fix schema or generator issues at the source.
- Keep English and Chinese summaries aligned when the verification scope changes.
