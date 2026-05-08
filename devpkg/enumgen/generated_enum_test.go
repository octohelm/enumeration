package enumgen_test

import (
	"fmt"
	"testing"

	. "github.com/octohelm/x/testing/v2"

	"github.com/octohelm/enumeration/devpkg/enumgen/testdata/model"
	"github.com/octohelm/enumeration/pkg/enumeration"
)

type offsetGender model.Gender

func (offsetGender) Offset() int {
	return -3
}

type offsetColor struct {
	model.Color
}

func (offsetColor) Offset() int {
	return 1000
}

func TestGeneratedEnumContracts(t *testing.T) {
	t.Run("运行期接口契约", func(t *testing.T) {
		var canValues enumeration.CanEnumValues = model.GENDER__MALE
		var withOffset enumeration.DriverValueOffset = offsetGender(model.GENDER__FEMALE)

		Then(
			t, "生成类型应实现枚举值接口",
			Expect(canValues.EnumValues(), Equal([]any{model.GENDER__MALE, model.GENDER__FEMALE})),
		)

		Then(
			t, "偏移类型应暴露数据库值偏移",
			Expect(withOffset.Offset(), Equal(-3)),
		)
	})

	t.Run("生成辅助方法行为", func(t *testing.T) {
		Then(
			t, "字符串与标签解析应符合样例定义",
			ExpectMustValue(func() (model.Gender, error) {
				return model.ParseGenderFromString("FEMALE")
			}, Equal(model.GENDER__FEMALE)),
			ExpectMustValue(func() (model.Gender, error) {
				return model.ParseGenderLabelString("男")
			}, Equal(model.GENDER__MALE)),
		)

		Then(
			t, "自定义值来源的 token 应按常量名生成字符串",
			Expect(model.EOF.String(), Equal("EOF")),
			Expect(model.COMMENT.Label(), Equal("COMMENT")),
		)
	})
}

func TestGeneratedGenderEnum(t *testing.T) {
	t.Run("IsZero 应基于 GENDER_UNKNOWN", func(t *testing.T) {
		Then(
			t, "零值返回 true",
			Expect(model.Gender(0).IsZero(), Equal(true)),
		)
		Then(
			t, "非零值返回 false",
			Expect(model.GENDER__MALE.IsZero(), Equal(false)),
		)
	})

	t.Run("String 应返回后缀名称", func(t *testing.T) {
		Then(
			t, "应返回后缀名称",
			Expect(model.GENDER__MALE.String(), Equal("MALE")),
			Expect(model.GENDER__FEMALE.String(), Equal("FEMALE")),
		)
	})

	t.Run("ParseFromString 应识别后缀名称", func(t *testing.T) {
		Then(
			t, "应识别后缀名称",
			ExpectMustValue(func() (model.Gender, error) {
				return model.ParseGenderFromString("MALE")
			}, Equal(model.GENDER__MALE)),
		)
	})

	t.Run("Value 应返回正确数值", func(t *testing.T) {
		Then(
			t, "应返回正确数值",
			ExpectMustValue(func() (any, error) {
				return model.GENDER__MALE.Value()
			}, Equal[any](int64(1))),
		)
	})

	t.Run("MarshalText / UnmarshalText 应往返一致", func(t *testing.T) {
		text, err := model.GENDER__MALE.MarshalText()
		Then(
			t, "MarshalText 应成功",
			Expect(err, Equal[error](nil)),
			Expect(string(text), Equal("MALE")),
		)

		var v model.Gender
		err = v.UnmarshalText(text)
		Then(
			t, "UnmarshalText 应还原原值",
			Expect(err, Equal[error](nil)),
			Expect(v, Equal(model.GENDER__MALE)),
		)
	})

	t.Run("Label 应返回中文注释", func(t *testing.T) {
		Then(
			t, "应返回中文注释",
			Expect(model.GENDER__MALE.Label(), Equal("男")),
			Expect(model.GENDER__FEMALE.Label(), Equal("女")),
		)
	})
}

func TestGeneratedColorEnum(t *testing.T) {
	t.Run("ConstValue int 枚举应使用字面值", func(t *testing.T) {
		Then(
			t, "String 应返回字面字符串",
			Expect(model.COLOR__RED.String(), Equal("100")),
			Expect(model.COLOR__GREEN.String(), Equal("200")),
			Expect(model.COLOR__BLUE.String(), Equal("300")),
		)

		Then(
			t, "ParseFromString 应识别字面值",
			ExpectMustValue(func() (model.Color, error) {
				return model.ParseColorFromString("200")
			}, Equal(model.COLOR__GREEN)),
		)

		Then(
			t, "未知值应回到 UNKNOWN 模式",
			ExpectMustValue(func() (model.Color, error) {
				return model.ParseColorFromString("UNKNOWN_999")
			}, Equal(model.Color(999))),
		)

		Then(
			t, "IsZero 应基于 COLOR_UNKNOWN",
			Expect(model.Color(0).IsZero(), Equal(true)),
			Expect(model.COLOR__RED.IsZero(), Equal(false)),
		)
	})

	t.Run("Value 应返回正确的数据库值", func(t *testing.T) {
		Then(
			t, "无偏移时返回原始 int64",
			ExpectMustValue(func() (any, error) {
				return model.COLOR__BLUE.Value()
			}, Equal[any](int64(300))),
		)

		Then(
			t, "嵌入类型的 Value 不感知外层偏移（Go 嵌入语义）",
			ExpectMustValue(func() (any, error) {
				return offsetColor{Color: model.COLOR__RED}.Value()
			}, Equal[any](int64(100))),
		)
	})

	t.Run("MarshalText / UnmarshalText 应往返一致", func(t *testing.T) {
		text, err := model.COLOR__GREEN.MarshalText()
		Then(
			t, "MarshalText 应成功",
			Expect(err, Equal[error](nil)),
			Expect(string(text), Equal("200")),
		)

		var v model.Color
		err = v.UnmarshalText(text)
		Then(
			t, "UnmarshalText 应还原原值",
			Expect(err, Equal[error](nil)),
			Expect(v, Equal(model.COLOR__GREEN)),
		)
	})
}

func TestGeneratedFlagEnum(t *testing.T) {
	t.Run("自定义 const-unknown-name 的枚举", func(t *testing.T) {
		Then(
			t, "应使用 FLAG_NONE 作为零值",
			Expect(model.Flag(0).IsZero(), Equal(true)),
			Expect(model.FLAG__READ.IsZero(), Equal(false)),
		)

		Then(
			t, "零值的 String 应在 UNKNOWN 模式给出可解析表示",
			Expect(model.Flag(0).String(), Equal("0")),
		)

		Then(
			t, "String 应使用后缀名称",
			Expect(model.FLAG__READ.String(), Equal("READ")),
			Expect(model.FLAG__WRITE.String(), Equal("WRITE")),
			Expect(model.FLAG__EXEC.String(), Equal("EXEC")),
		)

		Then(
			t, "ParseFromString 应识别后缀名称",
			ExpectMustValue(func() (model.Flag, error) {
				return model.ParseFlagFromString("READ")
			}, Equal(model.FLAG__READ)),
		)

		Then(
			t, "不存在常量名时回到 FLAG_NONE",
			Expect(model.Flag(99).IsZero(), Equal(false)),
		)
	})

	t.Run("MarshalText / UnmarshalText 应往返一致", func(t *testing.T) {
		text, err := model.FLAG__WRITE.MarshalText()
		Then(
			t, "MarshalText 应成功",
			Expect(err, Equal[error](nil)),
			Expect(string(text), Equal("WRITE")),
		)

		var v model.Flag
		err = v.UnmarshalText(text)
		Then(
			t, "UnmarshalText 应还原原值",
			Expect(err, Equal[error](nil)),
			Expect(v, Equal(model.FLAG__WRITE)),
		)
	})
}

func TestGeneratedTokenEnum(t *testing.T) {
	t.Run("value-from=ConstName 的枚举", func(t *testing.T) {
		Then(
			t, "应使用 ILLEGAL 作为零值",
			Expect(model.Token(0).IsZero(), Equal(true)),
			Expect(model.EOF.IsZero(), Equal(false)),
		)

		Then(
			t, "零值的 String 应返回 ILLEGAL",
			Expect(model.Token(0).String(), Equal("ILLEGAL")),
		)

		Then(
			t, "String 应使用常量名",
			Expect(model.EOF.String(), Equal("EOF")),
			Expect(model.COMMENT.String(), Equal("COMMENT")),
		)

		Then(
			t, "ParseFromString 应识别常量名",
			ExpectMustValue(func() (model.Token, error) {
				return model.ParseTokenFromString("EOF")
			}, Equal(model.EOF)),
		)

		Then(
			t, "Value 应返回正确数值",
			ExpectMustValue(func() (any, error) {
				return model.EOF.Value()
			}, Equal[any](int64(1))),
		)

		Then(
			t, "Label 应回退到常量名（无注释）",
			Expect(model.EOF.Label(), Equal("EOF")),
		)
	})

	t.Run("MarshalText / UnmarshalText 应往返一致", func(t *testing.T) {
		text, err := model.COMMENT.MarshalText()
		Then(
			t, "MarshalText 应成功",
			Expect(err, Equal[error](nil)),
			Expect(string(text), Equal("COMMENT")),
		)

		var v model.Token
		err = v.UnmarshalText(text)
		Then(
			t, "UnmarshalText 应还原原值",
			Expect(err, Equal[error](nil)),
			Expect(v, Equal(model.COMMENT)),
		)
	})
}

func TestGeneratedGenderExtEnum(t *testing.T) {
	t.Run("string 类型枚举", func(t *testing.T) {
		Then(
			t, "EnumValues 应包含所有常量",
			Expect(model.GenderExt("").EnumValues(), Equal([]any{
				model.GENDER_EXT__FEMALE, model.GENDER_EXT__MALE,
			})),
		)

		Then(
			t, "Label 应返回中文注释",
			Expect(model.GENDER_EXT__MALE.Label(), Equal("男")),
			Expect(model.GENDER_EXT__FEMALE.Label(), Equal("女")),
		)

		Then(
			t, "ParseLabelString 应识别中文标签",
			ExpectMustValue(func() (model.GenderExt, error) {
				return model.ParseGenderExtLabelString("女")
			}, Equal(model.GENDER_EXT__FEMALE)),
		)

		Then(
			t, "未知标签应返回零值并报告错误",
			ExpectDo(func() error {
				v, err := model.ParseGenderExtLabelString("未知")
				if v != model.GenderExt("") {
					return fmt.Errorf("期望零值, 但得到 %v", v)
				}
				return err
			}, ErrorIs(model.InvalidGenderExt)),
		)
	})
}

func TestScanEnumValues(t *testing.T) {
	t.Run("Scan 应正确解析整数输入", func(t *testing.T) {
		var g model.Gender
		err := g.Scan(int64(1))
		Then(
			t, "int64(1) 应解析为 MALE",
			Expect(err, Equal[error](nil)),
			Expect(g, Equal(model.GENDER__MALE)),
		)

		var c model.Color
		err = c.Scan(int64(100))
		Then(
			t, "int64(100) 应解析为 RED",
			Expect(err, Equal[error](nil)),
			Expect(c, Equal(model.COLOR__RED)),
		)
	})

	t.Run("Scan 传入字符串也应正确解析", func(t *testing.T) {
		var g model.Gender
		err := g.Scan("1")
		Then(
			t, "字符串 '1' 应解析为 MALE",
			Expect(err, Equal[error](nil)),
			Expect(g, Equal(model.GENDER__MALE)),
		)

		var g2 model.Gender
		err = g2.Scan([]byte("0"))
		Then(
			t, "字节串 '0' 应解析为 UNKNOWN",
			Expect(err, Equal[error](nil)),
			Expect(g2, Equal(model.GENDER_UNKNOWN)),
		)
	})
}

func TestGeneratedEnumValuesCoverage(t *testing.T) {
	t.Run("所有 int 枚举类型的 EnumValues 应包含所有已定义常量", func(t *testing.T) {
		Then(
			t, "Gender",
			Expect(model.Gender(0).EnumValues(), Equal([]any{
				model.GENDER__MALE, model.GENDER__FEMALE,
			})),
		)

		Then(
			t, "Color",
			Expect(model.Color(0).EnumValues(), Equal([]any{
				model.COLOR__RED, model.COLOR__GREEN, model.COLOR__BLUE,
			})),
		)

		Then(
			t, "Flag",
			Expect(model.Flag(0).EnumValues(), Equal([]any{
				model.FLAG__READ, model.FLAG__WRITE, model.FLAG__EXEC,
			})),
		)

		Then(
			t, "Token",
			Expect(model.Token(0).EnumValues(), Equal([]any{
				model.EOF, model.COMMENT,
			})),
		)
	})
}
