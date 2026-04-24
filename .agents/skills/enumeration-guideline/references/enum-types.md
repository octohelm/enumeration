# Enum 类型定义指南

用于在本仓库中选择和书写不同类型的 enum 定义。

## 共同约定

所有 enum 都遵循以下基线：

1. 在目标类型声明前添加 `// +gengo:enum`。
2. 参与生成的常量必须是导出的。
3. 常量名优先使用 `TYPE__VALUE` 形式，让生成器能稳定识别类型前缀与值后缀。
4. 若常量尾部带中文注释，该注释会作为 `Label()` 和 `Parse...LabelString()` 的来源。
5. 生成器会跳过 internal 常量、以下划线开头的常量，以及被识别为 unknown 零值的常量。

参考样例应由宿主仓库自行提供；若宿主仓库没有样例，按本文给出的最小写法落地。

## gengo 体系下的使用方式

本仓库的 enum 不是独立运行的宏系统，而是挂在 `gengo` 生成链路上的一个 generator。

接入时至少要满足：

1. 在类型声明前写 `// +gengo:enum` 及其参数。
2. 在生成器进程里通过匿名导入注册 `enumgen`。
3. 用 `gengo` 的 `Context` 或 `Executor` 把目标包加入 `Entrypoint`。
4. 通过 `OutputFileBaseName` 让输出文件落成 `zz_generated.*.go`。

最小生成器入口形态如下：

```go
import (
	_ "github.com/octohelm/enumeration/devpkg/enumgen"
)

c, err := gengo.NewContext(&gengo.GeneratorArgs{
	Entrypoint: []string{
		"github.com/octohelm/enumeration/pkg/testdata/model",
	},
	OutputFileBaseName: "zz_generated",
	Globals: map[string][]string{
		"gengo:runtimedoc": {},
	},
})

if err := c.Execute(ctx, gengo.GetRegisteredGenerators()...); err != nil {
	panic(err)
}
```

使用时应注意：

1. `enumgen` 依赖 `gengo.Register` 在 `init()` 中完成注册，所以不能省略匿名导入。
2. `Entrypoint` 应指向包含 enum 类型定义的 Go package，而不是单个文件路径。
3. 生成结果默认会引用本仓库的运行期包，例如 `pkg/enumeration` 与 `pkg/scanner`。
4. 若同时启用其他 `gengo` generator，应确认 `Globals` 与输出文件名策略不会互相冲突。

测试或隔离验证时，优先复用 `github.com/octohelm/gengo/pkg/gengo/testingutil` 创建临时 module 并执行生成。

## 1. `int` enum

适用场景：

- 需要稳定顺序值。
- 需要 `String()`、`MarshalText()`、`UnmarshalText()`、`Scan()`、`Value()` 等 int stringer 风格辅助方法。
- 需要明确的 unknown/zero 值。

写法：

```go
// +gengo:enum
type Gender int

const (
	GENDER_UNKNOWN Gender = iota
	GENDER__MALE          // 男
	GENDER__FEMALE        // 女
)
```

约定：

1. 第一项通常定义为 unknown/zero 值。
2. 默认 unknown 名称按类型推导，例如 `Gender` 对应 `GENDER_UNKNOWN`。
3. 这类 enum 会生成字符串解析、label 解析、数据库扫描和值转换方法。

## 2. `string` enum

适用场景：

- 业务持久化值本身就是字符串。
- 不需要按偏移值扫描数据库整数。
- 仍然需要 `EnumValues()`、`Label()`、`Parse...LabelString()` 等辅助方法。

写法：

```go
// +gengo:enum
type GenderExt string

const (
	GENDER_EXT__MALE   GenderExt = "MALE"   // 男
	GENDER_EXT__FEMALE GenderExt = "FEMALE" // 女
)
```

约定：

1. `string` enum 默认不依赖 unknown 零值也能工作。
2. 若未写标签注释，`Label()` 会回退到生成值本身。

## 3. 按常量后缀取值的 enum

适用场景：

- 代码里希望保留带类型前缀的常量名。
- 运行期字符串值只想暴露后缀，例如 `GENDER__MALE` 对应 `MALE`。

这也是默认策略，不写参数时即为该行为：

```go
// +gengo:enum
type Gender int
```

规则：

1. 生成器优先拆 `TYPE__VALUE` 中 `__` 后面的后缀。
2. 若常量命中 unknown 规则，例如 `GENDER_UNKNOWN`，生成值会按 unknown 处理。

## 4. 按常量名取值的 enum

适用场景：

- 运行期字符串值必须保留完整常量名。
- 常量名本身就是对外协议的一部分，例如词法 token。

写法：

```go
// +gengo:enum
// +gengo:enum:value-from=ConstName
// +gengo:enum:const-unknown-name=ILLEGAL
type Token int

const (
	ILLEGAL Token = iota
	EOF
	COMMENT
)
```

约定：

1. `String()` 和解析方法直接使用 `EOF`、`COMMENT` 这类常量名。
2. 若 unknown 名称不符合默认推导，必须显式写 `const-unknown-name`。

## 5. 按常量字面值取值的 enum

适用场景：

- 运行期值必须来自常量右侧字面值，而不是常量名。
- 常见于字符串状态码、外部协议字面值或手工指定数字。

写法：

```go
// +gengo:enum
// +gengo:enum:value-from=ConstValue
type Status string

const (
	STATUS_UNKNOWN Status = ""
	STATUS__OK     Status = "ok" // 正常
)
```

约定：

1. `String()` 优先返回字面值，例如 `"ok"`。
2. 标签仍然来自常量注释；若无注释，则回退到生成值。

## 参数速查

`+gengo:enum[:args]` 当前支持：

- `+gengo:enum:value-from=ConstName`
- `+gengo:enum:value-from=ConstNameSuffix`
- `+gengo:enum:value-from=ConstValue`
- `+gengo:enum:const-unknown-name=<NAME>`

## 选型建议

- 需要数据库整数扫描和 unknown 零值：优先 `int` enum。
- 需要字面值直接出现在协议里：优先 `string` enum 或 `ConstValue`。
- 需要完整常量名作为协议值：用 `ConstName`。
- 只想保留语义后缀：用默认的 `ConstNameSuffix`。
