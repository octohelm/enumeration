/*
Package main 提供基于 gengo 的枚举代码生成入口。

通过匿名导入 devpkg/enumgen 注册生成器，对指定 entrypoint 中的
带 +gengo:enum 注解的类型执行代码生成，产物落为 zz_generated.*.go。
*/
package main
