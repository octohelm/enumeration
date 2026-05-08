# Enum 类型定义指南

用于按 `github.com/octohelm/enumeration` 约定选择和接入 enum 定义。

具体常量约定、声明语法和参数始终以 `go doc` 为准：
- `go doc github.com/octohelm/enumeration/devpkg/enumgen`

## gengo 接入方式

1. 匿名导入 `github.com/octohelm/enumeration/devpkg/enumgen` 触发 `init()` 注册。
2. 用 `gengo.NewExecutor` 创建执行器，将目标包加入 `Entrypoint`。
3. 通过 `OutputFileBaseName` 让输出文件落成 `zz_generated.*.go`。

测试或隔离验证时，复用 `github.com/octohelm/gengo/pkg/gengo/testingutil`。

## 选型建议

- 需要数据库整数扫描和 unknown 零值：`int` enum。
- 需要字面值直接出现在协议里：`string` enum 或 `value-from=ConstValue`。
- 需要完整常量名作为协议值：`value-from=ConstName`。
- 只想保留语义后缀：默认行为（`value-from=ConstNameSuffix`）。
