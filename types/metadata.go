package types

type Type int

const (
	Unknown Type = iota
	String
	Number
	List
)

type Metadata struct {
	Blocks *[]Block `@@*`
}

type Block struct {
	Key    string   `@Ident ("=")?`
	Value  *Value   `| @@`
	Blocks *[]Block `| "{" @@* "}"`
}

type Value struct {
	Str  *string  `@String`
	Num  *int     `| @Number`
	List []*Value `| "[" ( @@ ( "," @@ )* )? "]"`
}

func (v Value) Type() Type {
	switch {
	case v.Str != nil:
		return String
	case v.Num != nil:
		return Number
	case v.List != nil:
		return List
	default:
		return Unknown
	}
}
