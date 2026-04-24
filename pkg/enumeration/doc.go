/*
Package enumeration 提供枚举运行期最小接口约定。

它只定义生成代码与调用方共享的能力边界，例如枚举值枚举和数据库值偏移，
不负责扫描、格式化或代码生成实现。

定义的接口：

  - CanEnumValues — 返回全部可枚举值的集合。
  - DriverValueOffset — 声明数据库持久化值相对 Go 常量值的固定偏移量。
*/
package enumeration
