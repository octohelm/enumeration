package enumgen

import (
	"regexp"
	"testing"

	"github.com/octohelm/gengo/pkg/gengo"
	"github.com/octohelm/gengo/pkg/gengo/testingutil"
	"github.com/octohelm/x/cmp"
	. "github.com/octohelm/x/testing/v2"
)

func TestParseValueFrom(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		expect ValueFrom
	}{
		{name: "默认值", input: "", expect: ValueFromConstName},
		{name: "按后缀", input: "ConstNameSuffix", expect: ValueFromConstNameSuffix},
		{name: "按常量值", input: "ConstValue", expect: ValueFromConstValue},
		{name: "未知值回退默认", input: "Unknown", expect: ValueFromConstName},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			Then(t, "应解析为预期策略",
				Expect(ParseValueFrom(tt.input), Equal(tt.expect)),
			)
		})
	}
}

func TestEnumTypeGenerate(t *testing.T) {
	m := MustValue(t, func() (*testingutil.Module, error) {
		return testingutil.NewModule(t, map[string]string{
			"sample/types.go": `package sample

// +gengo:enum
type Gender int

const (
	GENDER_UNKNOWN Gender = iota
	GENDER__MALE
	GENDER__FEMALE
)

// +gengo:enum
type GenderExt string

const (
	GENDER_EXT__MALE   GenderExt = "MALE"   // 男
	GENDER_EXT__FEMALE GenderExt = "FEMALE" // 女
)

// +gengo:enum
// +gengo:enum:value-from=ConstValue
type Status string

const (
	STATUS_UNKNOWN Status = ""
	STATUS__OK     Status = "ok" // 好
)

const hidden Gender = 99
`,
		})
	})

	files := MustValue(t, func() (map[string]string, error) {
		return m.Generate(gengo.GeneratorArgs{
			Entrypoint:         []string{m.ImportPath("sample")},
			OutputFileBaseName: "zz_generated_test",
			Force:              true,
		}, &enumGen{})
	})

	Then(t, "应生成整数与字符串枚举辅助方法",
		Expect(files, Be(testingutil.File("sample/zz_generated_test.enum.go",
			testingutil.Contains(
				"var InvalidGender = errors.New(\"invalid Gender\")",
				"func ParseGenderFromString(s string) (Gender, error)",
				"func (v *Gender) Scan(src any) error",
				"func ParseGenderExtLabelString(label string) (GenderExt, error)",
				"case \"女\":",
				"func (v Status) Label() string",
				"return \"ok\"",
			),
			testingutil.NotContains(
				"hidden",
			),
			testingutil.Count("func (v *Gender) Scan(src any) error", cmp.Eq(1)),
		))),
	)
}

func TestEnumTypeGenerateIntConstValue(t *testing.T) {
	m := MustValue(t, func() (*testingutil.Module, error) {
		return testingutil.NewModule(t, map[string]string{
			"sample/types.go": `package sample

// +gengo:enum
// +gengo:enum:value-from=ConstValue
type Color int

const (
	COLOR_UNKNOWN Color = 0
	COLOR__RED    Color = 100
	COLOR__GREEN  Color = 200
	COLOR__BLUE   Color = 300
)
`,
		})
	})

	files := MustValue(t, func() (map[string]string, error) {
		return m.Generate(gengo.GeneratorArgs{
			Entrypoint:         []string{m.ImportPath("sample")},
			OutputFileBaseName: "zz_generated_test",
			Force:              true,
		}, &enumGen{})
	})

	Then(t, "ConstValue int 枚举应使用字面值作为 String 和 Parse 依据",
		Expect(files, Be(testingutil.File("sample/zz_generated_test.enum.go",
			testingutil.Contains(
				"return \"100\"",
				"return \"200\"",
				"return \"300\"",
				"case \"100\":",
				"case \"200\":",
				"case \"300\":",
				"func (v *Color) Scan(src any) error",
			),
		))),
	)
}

func TestEnumTypeGenerateCustomConstUnknownName(t *testing.T) {
	m := MustValue(t, func() (*testingutil.Module, error) {
		return testingutil.NewModule(t, map[string]string{
			"sample/types.go": `package sample

// +gengo:enum
// +gengo:enum:const-unknown-name=MY_ZERO
type Level int

const (
	MY_ZERO       Level = iota
	LEVEL__LOW
	LEVEL__MEDIUM
	LEVEL__HIGH
)
`,
		})
	})

	files := MustValue(t, func() (map[string]string, error) {
		return m.Generate(gengo.GeneratorArgs{
			Entrypoint:         []string{m.ImportPath("sample")},
			OutputFileBaseName: "zz_generated_test",
			Force:              true,
		}, &enumGen{})
	})

	Then(t, "自定义 ConstUnknownName 时应使用指定名称而非自动推导",
		Expect(files, Be(testingutil.File("sample/zz_generated_test.enum.go",
			testingutil.Contains(
				"MY_ZERO",
				"func (v Level) IsZero() bool",
			),
			testingutil.NotContains("LEVEL_UNKNOWN"),
		))),
	)
}

func TestEnumTypeGenerateError(t *testing.T) {
	m := MustValue(t, func() (*testingutil.Module, error) {
		return testingutil.NewModule(t, map[string]string{
			"sample/bad.go": `package sample

// +gengo:enum
type Bad int
`,
		})
	})

	Then(t, "缺少常量时不应生成枚举文件",
		ExpectDo(func() error {
			files, err := m.Generate(gengo.GeneratorArgs{
				Entrypoint:         []string{m.ImportPath("sample")},
				OutputFileBaseName: "zz_generated_test",
				Force:              true,
			}, &enumGen{})
			if err != nil {
				return err
			}
			return testingutil.File("sample/zz_generated_test.enum.go")(files)
		}, ErrorMatch(regexp.MustCompile("生成文件不存在"))),
	)
}
