---
name: enumeration-guideline
description: 说明 `github.com/octohelm/enumeration` 中枚举类型的定义约定与 `+gengo:enum[:args]` 用法；当需要新增、修改或审查 enum 定义时使用。
---

# Enumeration Guideline

按 `github.com/octohelm/enumeration` 约定定义枚举类型。

## 定义一个 enum

```go
// +gengo:enum
type Gender int

const (
    GENDER_UNKNOWN Gender = iota
    GENDER__MALE          // 男
    GENDER__FEMALE        // 女
)
```

**关键约定**：
- 常量名用 `TYPE__VALUE` 形式，生成器按 `__` 拆分前缀与后缀
- 第一项为 unknown/zero 值
- 常量尾部中文注释自动作为 `Label()` 来源
- enum 挂在 gengo 链路上，通过 `import _ "github.com/octohelm/enumeration/devpkg/enumgen"` 注册

**选 value-from**：
- 默认 `ConstNameSuffix`：`GENDER__MALE` → `"MALE"`
- `ConstName`：保留完整常量名（词法 token 等场景）
- `ConstValue`：取常量字面值（外部协议字面值）
- 需自定义 unknown 名时加 `+gengo:enum:const-unknown-name=<NAME>`

具体参数和常量约定以 `go doc github.com/octohelm/enumeration/devpkg/enumgen` 为准。

## 更多类型和选型

见 [references/enum-types.md](references/enum-types.md)——string enum、ConstName、ConstValue 的详细写法和使用场景。
