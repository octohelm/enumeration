package enumeration

// CanEnumValues 表示可返回全部可枚举值集合的类型。
type CanEnumValues interface {
	EnumValues() []any
}

// DriverValueOffset 表示数据库持久化值相对 Go 常量值存在固定偏移量。
type DriverValueOffset interface {
	Offset() int
}
