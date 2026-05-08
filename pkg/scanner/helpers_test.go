package scanner

import (
	"regexp"
	"testing"

	. "github.com/octohelm/x/testing/v2"
)

func TestScanIntEnumStringer(t *testing.T) {
	t.Run("成功解析多种输入类型", func(t *testing.T) {
		offset := -3

		testCases := []struct {
			name   string
			input  any
			expect int
		}{
			{name: "nil 使用默认值", input: nil, expect: 0},
			{name: "空字节串使用默认值", input: []byte(""), expect: 0},
			{name: "字节串数字", input: []byte("-2"), expect: 1},
			{name: "字符串数字", input: "-1", expect: 2},
			{name: "int", input: int(0), expect: 3},
			{name: "int8", input: int8(1), expect: 4},
			{name: "int16", input: int16(2), expect: 5},
			{name: "int32", input: int32(3), expect: 6},
			{name: "int64", input: int64(4), expect: 7},
			{name: "uint", input: uint(5), expect: 8},
			{name: "uint8", input: uint8(6), expect: 9},
			{name: "uint16", input: uint16(7), expect: 10},
			{name: "uint32", input: uint32(8), expect: 11},
			{name: "uint64", input: uint64(9), expect: 12},
			{name: "未知类型回退默认值", input: struct{}{}, expect: 0},
		}

		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				Then(
					t, "应返回去偏移后的索引",
					ExpectMustValue(func() (int, error) {
						return ScanIntEnumStringer(tt.input, offset)
					}, Equal(tt.expect)),
				)
			})
		}
	})

	t.Run("非法字符串返回错误并保留偏移值", func(t *testing.T) {
		Then(
			t, "应返回解析错误",
			ExpectDo(func() error {
				got, err := ScanIntEnumStringer("invalid", -2)
				if got != -2 {
					t.Fatalf("unexpected result: %d", got)
				}
				return err
			}, ErrorMatch(regexp.MustCompile("invalid syntax"))),
		)
	})
}

func TestScanEnum(t *testing.T) {
	Then(
		t, "应复用整数枚举扫描逻辑",
		ExpectMustValue(func() (int, error) {
			return ScanEnum("3", 1)
		}, Equal(2)),
	)
}
