/*
Package enumgen 提供基于 gengo 的枚举代码生成器实现。

该包负责从带注解的 Go 常量声明中提取枚举元数据，并生成运行期辅助方法。

支持的声明方式：

	// +gengo:enum
	type Gender int

	// +gengo:enum
	// +gengo:enum:value-from=ConstName
	// +gengo:enum:const-unknown-name=ILLEGAL
	type Token int

`+gengo:enum[:args]` 当前支持的参数：

- `value-from=ConstName`：使用常量名作为生成值。
- `value-from=ConstNameSuffix`：使用常量名去掉类型前缀后的后缀作为生成值。
- `value-from=ConstValue`：使用常量字面值作为生成值。
- `const-unknown-name=<NAME>`：指定 zero/unknown 常量名，用于生成 int stringer 风格辅助方法。
*/
package enumgen
