---
name: enumeration-guideline
description: 说明本仓库中枚举类型的定义约定与 `+gengo:enum[:args]` 用法；当需要新增、修改或审查 enum 定义时使用。
metadata:
  primary_pattern: tool-wrapper
---

# Enumeration Guideline

用于在本仓库中定义、修改或审查枚举类型时，快速选择合适的 enum 形式，并按生成器约定落代码。

## 适用范围

- 新增或修改带 `+gengo:enum` 注解的类型。
- 判断某个 enum 应该定义为 `int`、`string`，还是改用不同 `value-from` 策略。
- 审查 enum 常量命名、unknown 零值、label 注释和生成行为是否符合仓库约定。

## 输入

至少需要：

- 目标 enum 的业务语义
- 预期底层类型，例如 `int` 或 `string`
- 是否需要 unknown/zero 值
- 序列化或持久化值应来自常量名、常量后缀还是常量字面值

## 先看什么

1. 先读 [`references/enum-types.md`](references/enum-types.md) 选择 enum 类型。
2. 需要确认在 gengo 体系里的接入方式时，重点看其中的“gengo 体系下的使用方式”一节。
3. 需要对照宿主仓库现有样例时，优先查找宿主仓库自己的 enum 样例或生成入口。
4. 需要确认注解参数支持范围时，查看宿主仓库中 `enumgen` 的包文档或生成器说明。

## 工作方式

1. 先按 `references/enum-types.md` 判断当前 enum 属于哪一类。
2. 再确认是否需要 `+gengo:enum:value-from=...` 或 `+gengo:enum:const-unknown-name=...`。
3. 最后检查常量是否导出、unknown 常量是否清晰、label 注释是否写在常量行尾。

## 完成标准

- enum 类型与 `value-from` 策略匹配。
- 常量命名能被生成器稳定识别。
- 需要 label 时，标签注释位置正确。
- 需要 unknown 零值时，常量名与注解参数一致。
