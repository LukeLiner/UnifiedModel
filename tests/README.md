# tests/

测试套件，按 spec 11 分层。

| 目录 | 说明 |
|---|---|
| `unit/` | 模块内部逻辑单元测试 |
| `contract/` | 公共契约、Storage interface、MCP、SDK 契约测试 |
| `integration/` | 跨模块路径集成测试 |
| `e2e/` | Quickstart 端到端测试 |
| `golden/` | 查询语义、explain、docs 示例 golden tests |
| `architecture/` | 依赖方向、禁用 API architecture guard 测试 |
