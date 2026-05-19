# Java TestRunner

English version: [TestRunner.md](TestRunner.md)

本文档说明 Java SDK 生成结果的测试入口。TestRunner 用于验证生成类型、序列化/反序列化以及基础行为。

## 运行

从仓库根目录执行：

```bash
make verify-java
```

或进入 Java SDK 目录运行 Maven 测试，具体命令以 `Makefile` 中的 `verify-java` 为准。

## 维护规则

- schema 或 Java generator 变化后运行 Java 验证。
- 测试失败时优先修复 schema/generator，而不是手工改生成文件。
- 行为变化应同步更新英文和中文文档。
