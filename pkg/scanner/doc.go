/*
Package scanner 提供枚举相关的最小扫描辅助函数。

它负责把数据库或外部输入（int、string、[]byte）转换为去偏移后的整数枚举索引，
不承载具体枚举类型定义。

主要函数：

  - ScanEnum(src, offset) — 通用扫描入口，委托给 ScanIntEnumStringer。
  - ScanIntEnumStringer(src, offset) — 将输入解析为枚举索引，减去偏移量后返回。
*/
package scanner

// P 是供代码生成阶段引用包路径的公开标记类型。
type P struct{}
