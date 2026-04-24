package enumgen

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/types"
	"sort"
	"strconv"
	"strings"

	"github.com/octohelm/gengo/pkg/gengo"
)

type ValueFrom int

const (
	// ValueFromConstName 使用常量名作为生成值。
	ValueFromConstName ValueFrom = iota
	// ValueFromConstNameSuffix 使用常量名去掉类型前缀后的后缀部分作为生成值。
	ValueFromConstNameSuffix
	// ValueFromConstValue 使用常量字面值作为生成值。
	ValueFromConstValue
)

// ParseValueFrom 将配置字符串解析为取值策略。
func ParseValueFrom(s string) ValueFrom {
	switch s {
	case "ConstNameSuffix":
		return ValueFromConstNameSuffix
	case "ConstValue":
		return ValueFromConstValue
	}
	return ValueFromConstName
}

// EnumTypes 按包路径和目标类型维护枚举元数据索引。
type EnumTypes map[string]map[types.Type]*EnumType

// ResolveEnumType 根据类型查找已收集的枚举元数据。
func (e EnumTypes) ResolveEnumType(t types.Type) (*EnumType, bool) {
	if n, ok := t.(*types.Named); ok {
		if enumTypes, ok := e[n.Obj().Pkg().Path()]; ok && enumTypes != nil {
			if enumType, ok := enumTypes[t]; ok {
				return enumType, ok
			}
		}
	}
	return nil, false
}

// Walk 遍历指定包的常量并收集枚举元数据。
func (e EnumTypes) Walk(gc gengo.Context, inPkgPath string) {
	p := gc.Package(inPkgPath)

	constants := p.Constants()

	for k := range p.Constants() {
		if e[inPkgPath] == nil {
			e[inPkgPath] = map[types.Type]*EnumType{}
		}

		constv := constants[k]

		if e[inPkgPath][constv.Type()] == nil {
			et := &EnumType{
				ValueFrom: ValueFromConstNameSuffix,
			}

			switch named := constv.Type().(type) {
			case *types.Named:
				et.ConstUnknownName = fmt.Sprintf("%s_UNKNOWN", gengo.UpperSnakeCase(named.Obj().Name()))

				opts := gc.OptsOf(named.Obj(), "enum")

				if valueFrom, ok := opts.Get("valueFrom"); ok {
					et.ValueFrom = ParseValueFrom(valueFrom)
				}

				if constUnknownName, ok := opts.Get("constUnknownName"); ok {
					et.ConstUnknownName = constUnknownName
				}
			}

			e[inPkgPath][constv.Type()] = et
		}

		e[inPkgPath][constv.Type()].Add(constv, p.Comment(constv.Pos())...)
	}
}

// EnumType 描述单个枚举类型的生成期元数据。
type EnumType struct {
	ValueFrom        ValueFrom
	ConstUnknownName string

	ConstUnknown *types.Const
	Constants    []*types.Const
	Comments     map[*types.Const][]string
}

// IsIntStringer 返回当前枚举是否可生成 int stringer 风格辅助方法。
func (e *EnumType) IsIntStringer() bool {
	return e.ConstUnknown != nil && len(e.Constants) > 0
}

// Label 返回常量对应的人类可读标签。
func (e *EnumType) Label(cv *types.Const) string {
	if comments, ok := e.Comments[cv]; ok {
		label := strings.Join(comments, "")

		if label != "" {
			return label
		}
	}

	return fmt.Sprintf("%v", e.Value(cv))
}

// Value 返回常量对应的生成值。
func (e *EnumType) Value(cv *types.Const) any {
	if named, ok := cv.Type().(*types.Named); ok {
		switch e.ValueFrom {
		case ValueFromConstName:
			return cv.Name()
		case ValueFromConstNameSuffix:
			parts := strings.SplitN(cv.Name(), "__", 2)
			if len(parts) == 2 && parts[0] == gengo.UpperSnakeCase(named.Obj().Name()) {
				return parts[1]
			}
			if strings.HasSuffix(cv.Name(), "_UNKNOWN") {
				return "UNKNOWN"
			}
		case ValueFromConstValue:
			//
		}
	}

	val := cv.Val()

	if val.Kind() == constant.Int {
		i, _ := strconv.ParseInt(val.String(), 10, 64)
		return i
	}

	s, _ := strconv.Unquote(val.String())
	return s
}

// Add 向当前枚举元数据中注册一个常量。
func (e *EnumType) Add(cv *types.Const, comments ...string) {
	if e.Comments == nil {
		e.Comments = map[*types.Const][]string{}
	}

	n := cv.Name()

	// skip internal
	if !ast.IsExported(n) {
		return
	}

	if n[0] == '_' {
		return
	}

	if n == e.ConstUnknownName {
		e.ConstUnknown = cv
		return
	}

	e.Comments[cv] = comments

	parts := strings.SplitN(n, "__", 2)

	if len(parts) == 2 {
		names := strings.Split(cv.Type().String(), ".")
		name := names[len(names)-1]

		if gengo.UpperSnakeCase(name) == parts[0] {
			e.Constants = append(e.Constants, cv)
		}
	} else {
		e.Constants = append(e.Constants, cv)
	}

	sort.Slice(e.Constants, func(i, j int) bool {
		return e.Constants[i].Val().String() < e.Constants[j].Val().String()
	})
}
