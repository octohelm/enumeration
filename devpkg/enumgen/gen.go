package enumgen

import (
	"database/sql/driver"
	"fmt"
	"go/types"

	"github.com/octohelm/enumeration/pkg/enumeration"
	"github.com/octohelm/enumeration/pkg/scanner"
	"github.com/octohelm/gengo/pkg/gengo"
	"github.com/octohelm/gengo/pkg/gengo/snippet"
)

func init() {
	gengo.Register(&enumGen{})
}

type enumGen struct {
	enumTypes EnumTypes
}

func (*enumGen) Name() string {
	return "enum"
}

func (*enumGen) New(ctx gengo.Context) gengo.Generator {
	g := &enumGen{
		enumTypes: EnumTypes{},
	}
	g.enumTypes.Walk(ctx, ctx.Package("").Pkg().Path())
	return g
}

func (g *enumGen) GenerateType(c gengo.Context, named *types.Named) error {
	if enum, ok := g.enumTypes.ResolveEnumType(named); ok {
		if enum.IsIntStringer() {
			g.genIntStringerEnums(c, named, enum)
			return nil
		}
		g.genEnums(c, named, enum)
	}
	return nil
}

func (g *enumGen) genEnums(c gengo.Context, named *types.Named, enum *EnumType) {
	options := make([]Option, len(enum.Constants))
	tpeObj := named.Obj()

	for i := range enum.Constants {
		options[i].Name = enum.Constants[i].Name()
		options[i].Value = fmt.Sprintf("%v", enum.Value(enum.Constants[i]))
		options[i].Label = enum.Label(enum.Constants[i])
	}

	c.RenderT(`
var Invalid@Type = @errorsNew("invalid @Type")

func (@Type) EnumValues() []any {
	return []any{
		@constValues
	}
}
`, snippet.Args{
		"errorsNew": snippet.PkgExpose("errors", "New"),
		"Type":      snippet.ID(tpeObj),
		"constValues": snippet.Snippets(func(yield func(snippet.Snippet) bool) {
			for _, o := range options {
				if !yield(snippet.Sprintf("%T,", o.Name)) {
					return
				}
			}
		}),
	})

	g.genLabel(c, tpeObj, enum, options)
}

type Option struct {
	Name  string
	Label string
	Value any
}

func (g *enumGen) genLabel(c gengo.Context, typ *types.TypeName, enum *EnumType, options []Option) {
	c.RenderT(`
func Parse@Type'LabelString(label string) (@Type, error) {
	switch label {
		@labelToConstCases
		default:
			return @ConstUnknown, Invalid@Type
	}
}

func (v @Type) Label() string {
	switch v {
		@constToLabelCases
		default:
			return @fmtSprint(v)
	}
}

`,
		snippet.Args{
			"Type": snippet.ID(typ.Name()),
			"ConstUnknown": snippet.Snippets(func(yield func(snippet.Snippet) bool) {
				if enum.ConstUnknown != nil {
					if !yield(snippet.ID(enum.ConstUnknown)) {
						return
					}
					return
				}

				if !yield(snippet.Value("")) {
					return
				}
				return
			}),
			"fmtSprint": snippet.PkgExpose("fmt", "Sprint"),
			"labelToConstCases": snippet.Snippets(func(yield func(snippet.Snippet) bool) {
				for _, o := range options {
					if !yield(snippet.Sprintf("case %v:\n\treturn %T,nil\n", o.Label, o.Name)) {
					}
				}
			}),
			"constToLabelCases": snippet.Snippets(func(yield func(snippet.Snippet) bool) {
				for _, o := range options {
					if !yield(snippet.Sprintf("case %T:\n\treturn %v\n", o.Name, o.Label)) {
					}
				}
			}),
		})
}

func (g *enumGen) genIntStringerEnums(c gengo.Context, tpe types.Type, enum *EnumType) {
	options := make([]Option, len(enum.Constants))
	tpeObj := tpe.(*types.Named).Obj()

	for i := range enum.Constants {
		options[i].Name = enum.Constants[i].Name()
		options[i].Value = fmt.Sprintf("%v", enum.Value(enum.Constants[i]))
		options[i].Label = enum.Label(enum.Constants[i])
	}

	c.RenderT(`
var Invalid@Type = @errorsNew("invalid @Type")

func (@Type) EnumValues() []any {
	return []any{
		@constValues
	}
}
`, snippet.Args{
		"Type":      snippet.ID(tpeObj.Name()),
		"errorsNew": snippet.PkgExpose("errors", "New"),
		"constValues": snippet.Snippets(func(yield func(snippet.Snippet) bool) {
			for _, o := range options {
				if !yield(snippet.Sprintf("%T,", o.Name)) {
					return
				}
			}
		}),
	})

	c.RenderT(`
func (v @Type) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v *@Type) UnmarshalText(data []byte) (error) {
	vv, err := Parse@Type'FromString(string(@bytesToUpper(data)))
	if err != nil {
		return err
	}
	*v = vv
	return nil
}

func Parse@Type'FromString(s string) (@Type, error) {
	switch s {
		@strValueToConstCases
		default:
			var i @Type
			_, err := @fmtSscanf(s, "UNKNOWN_%d", &i)
			if err == nil {
				return i, nil
			}
			return @ConstUnknown, Invalid@Type
	}
}

func (v @Type) IsZero() bool {
	return v == @ConstUnknown
}

func (v @Type) String() string {
	switch v {
		@constToStrValueCases
		case @ConstUnknown:
            return "UNKNOWN"
		default:
			return @fmtSprintf("UNKNOWN_%d", v)
	}
}

`, snippet.Args{
		"Type": snippet.ID(tpeObj.Name()),
		"ConstUnknown": snippet.Snippets(func(yield func(snippet.Snippet) bool) {
			if enum.ConstUnknown != nil {
				if !yield(snippet.ID(enum.ConstUnknown)) {
					return
				}
				return
			}

			if !yield(snippet.Value("")) {
				return
			}

			return
		}),
		"stringsHasPrefix": snippet.PkgExpose("strings", "HasPrefix"),
		"fmtSscanf":        snippet.PkgExpose("fmt", "Sscanf"),
		"fmtSprintf":       snippet.PkgExpose("fmt", "Sprintf"),
		"bytesToUpper":     snippet.PkgExpose("bytes", "ToUpper"),
		"strValueToConstCases": snippet.Snippets(func(yield func(snippet.Snippet) bool) {
			for _, o := range options {
				if !yield(snippet.Sprintf("case %v:\n\treturn %T, nil\n", o.Value, o.Name)) {
				}
			}

		}),
		"constToStrValueCases": snippet.Snippets(func(yield func(snippet.Snippet) bool) {
			for _, o := range options {
				if !yield(snippet.Sprintf("case %T:\n\treturn %v\n", o.Name, o.Value)) {
				}
			}
		}),
	})

	g.genLabel(c, tpeObj, enum, options)

	c.RenderT(`
func (v @Type) Value() (@driverValue, error) {
	offset := 0
	if o, ok := any(v).(@enumerationDriverValueOffset); ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}

func (v *@Type) Scan(src any) error {
	offset := 0
	if o, ok := any(v).(@enumerationDriverValueOffset); ok {
		offset = o.Offset()
	}

	i, err := @scannerScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*v = @Type(i)
	return nil
}

`, snippet.Args{
		"Type":                         snippet.ID(tpeObj),
		"driverValue":                  snippet.PkgExposeFor[driver.Value](),
		"enumerationDriverValueOffset": snippet.PkgExposeFor[enumeration.DriverValueOffset](),
		"scannerScanIntEnumStringer":   snippet.PkgExposeFor[scanner.P]("ScanIntEnumStringer"),
	})
}
