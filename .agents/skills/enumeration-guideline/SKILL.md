---
name: enumeration-guideline
description: 说明 `github.com/octohelm/enumeration` 中枚举类型的定义约定与 `+gengo:enum[:args]` 用法；当需要新增、修改或审查 enum 定义时使用。
---

# Enumeration Guideline

用于稳定使用 `github.com/octohelm/enumeration` 定义、修改或审查枚举类型，按生成器约定落代码。

## 使用范围

- 新增或修改带 `+gengo:enum` 注解的类型。
- 判断某个 enum 应定义为 `int`、`string`，还是改用不同 `value-from` 策略。
- 审查 enum 常量命名、unknown 零值、label 注释和生成行为是否符合约定。

## 必要输入

- 目标 enum 的业务语义。
- 预期底层类型（`int` 或 `string`）。
- 是否需要 `value-from` 参数或自定义 `const-unknown-name`。

## 关键约定

- 常量名优先使用 `TYPE__VALUE` 形式，让生成器稳定识别类型前缀与值后缀。
- 第一项通常定义为 unknown/zero 值。
- 若常量尾部带中文注释，该注释会作为 `Label()` 的来源。
- enum 挂在 `gengo` 生成链路上，通过匿名导入注册 `enumgen`。
- 具体 API 签名和参数以 `go doc` 为准，不在 skill 中复制手册。

## 资源导航

- 要选择 enum 类型和写法，读 [`references/enum-types.md`](references/enum-types.md)。
- 要查阅生成器支持的完整参数，运行 `go doc github.com/octohelm/enumeration/devpkg/enumgen`。
- 要查阅运行期接口，运行 `go doc github.com/octohelm/enumeration/pkg/enumeration` 和 `go doc github.com/octohelm/enumeration/pkg/scanner`。

## 工作方式

1. 先按 `references/enum-types.md` 判断当前 enum 属于哪一类。
2. 再确认是否需要 `+gengo:enum:value-from=...` 或 `+gengo:enum:const-unknown-name=...`。
3. 最后检查常量是否导出、unknown 常量是否清晰、label 注释是否写在常量行尾。

## 完成标准

- enum 类型与 `value-from` 策略匹配。
- 常量命名能被生成器稳定识别。
- 需要 label 时标签注释位置正确。
- 需要 unknown 零值时常量名与注解参数一致。
