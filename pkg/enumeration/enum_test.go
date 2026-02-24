package enumeration_test

import (
	"testing"

	"github.com/octohelm/enumeration/pkg/scanner"
	. "github.com/octohelm/x/testing/v2"
)

func TestScanEnum(t *testing.T) {
	t.Run("ScanIntEnumStringer", func(t *testing.T) {
		t.Run("Offset 为 -3 的偏移量扫描", func(t *testing.T) {
			offset := -3

			// 定义输入输出映射以增强可读性
			testCases := []struct {
				input  any
				expect int
			}{
				{nil, 0},
				{[]byte("-3"), 0},
				{"-2", 1},
				{int(-1), 2},
				{int8(0), 3},
				{int16(1), 4},
				{int32(2), 5},
				{int64(3), 6},
				{uint(4), 7},
				{uint8(5), 8},
				{uint16(6), 9},
				{uint32(7), 10},
				{uint64(8), 11},
			}

			for _, tc := range testCases {
				Then(t, "输入不同类型的数值或字符串应正确计算枚举索引",
					ExpectMustValue(
						func() (int, error) {
							return scanner.ScanIntEnumStringer(tc.input, offset)
						},
						Equal(tc.expect),
					),
				)
			}
		})
	})
}
