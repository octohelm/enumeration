# enumeration

`enumeration` 是一个面向 Go 的枚举支持库，围绕枚举值描述、扫描与代码生成提供最小实现面。

仓库当前同时包含运行期库代码、生成器实现和最小测试样例。它负责沉淀枚举相关的通用能力，不负责承载独立业务项目或多模块 workspace 管理。

## 职责与边界

- 提供枚举运行期接口与辅助实现，代码位于 `pkg/`。
- 提供枚举生成器及开发期支撑代码，代码位于 `devpkg/` 与 `tool/`。
- 提供用于验证生成结果的最小测试数据，位于 `testdata/`。
- 仓库级协作约束与统一执行入口分别收敛到 `AGENTS.md` 与 `justfile`，不在 README 展开命令细节。

## 从哪里开始

- `pkg/enumeration`：运行期枚举接口与基础类型约束。
- `devpkg/enumgen`：基于 `gengo` 的枚举代码生成实现。
- `devpkg/enumgen/testdata/model`：枚举生成器输入与输出样例，适合快速理解生成约定。
- `tool/go/justfile`：Go 工具链入口，包含依赖整理、测试、格式化和生成命令。

## 协作与执行入口

- `AGENTS.md`：仓库级协作规则与停机门禁。
- `justfile`：统一执行入口；可通过 `just --list --list-submodules` 查看所有已暴露命令。
