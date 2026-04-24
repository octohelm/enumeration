package model

// +gengo:enum
type Gender int

const (
	GENDER_UNKNOWN Gender = iota
	GENDER__MALE          // 男
	GENDER__FEMALE        // 女
)

// +gengo:enum
type GenderExt string

const (
	GENDER_EXT__MALE   GenderExt = "MAIL"   // 男
	GENDER_EXT__FEMALE GenderExt = "FEMALE" // 女
)

// +gengo:enum
// +gengo:enum:value-from=ConstName
// +gengo:enum:const-unknown-name=ILLEGAL
type Token int

const (
	ILLEGAL Token = iota
	EOF
	COMMENT
)

// +gengo:enum
// +gengo:enum:value-from=ConstValue
type Color int

const (
	COLOR_UNKNOWN Color = 0
	COLOR__RED    Color = 100
	COLOR__GREEN  Color = 200
	COLOR__BLUE   Color = 300
)

// +gengo:enum
// +gengo:enum:const-unknown-name=FLAG_NONE
type Flag int

const (
	FLAG_NONE Flag = iota
	FLAG__READ
	FLAG__WRITE
	FLAG__EXEC
)
