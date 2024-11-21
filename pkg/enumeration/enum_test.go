package enumeration_test

import (
	"github.com/octohelm/enumeration/pkg/scanner"
	"testing"

	testingx "github.com/octohelm/x/testing"
)

func TestScanEnum(t *testing.T) {
	cases := []struct {
		offset int
		values []interface{}
		expect []int
	}{
		{
			-3,
			[]interface{}{
				nil,
				[]byte("-3"),
				"-2",
				int(-1),
				int8(0),
				int16(1),
				int32(2),
				int64(3),
				uint(4),
				uint8(5),
				uint16(6),
				uint32(7),
				uint64(8),
			},
			[]int{
				0,
				0,
				1,
				2,
				3,
				4,
				5,
				6,
				7,
				8,
				9,
				10,
				11,
			},
		},
	}

	for _, c := range cases {
		for i, v := range c.values {
			n, err := scanner.ScanIntEnumStringer(v, c.offset)
			testingx.Expect(t, err, testingx.Be[error](nil))
			testingx.Expect(t, n, testingx.Equal(c.expect[i]))
		}
	}
}
