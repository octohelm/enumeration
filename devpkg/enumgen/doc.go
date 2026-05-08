/*
Package enumgen 提供基于 gengo 的枚举代码生成器实现。

该包负责从带注解的 Go 常量声明中提取枚举元数据，并生成运行期辅助方法。

常量约定：

  - 常量名使用 TYPE__VALUE 形式，生成器按 __ 拆分为类型前缀与值后缀。
  - 常量尾部带中文注释时，该注释会作为 Label() 的来源。
  - 以下常量会被跳过：internal 常量、以下划线开头的常量、被识别为 unknown 零值的常量。
  - int enum 的第一项默认为 unknown/zero 值，可通过 const-unknown-name 自定义。

声明方式：

	// +gengo:enum
	type Gender int

	// +gengo:enum:value-from=ConstName
	// +gengo:enum:const-unknown-name=ILLEGAL
	type Token int

`+gengo:enum[:args]` 支持的参数：

  - value-from=ConstName：使用常量名作为生成值。
  - value-from=ConstNameSuffix：使用常量名去掉类型前缀后的后缀作为生成值。
  - value-from=ConstValue：使用常量字面值作为生成值。
  - const-unknown-name=<NAME>：指定 zero/unknown 常量名。
*/
package enumgen
