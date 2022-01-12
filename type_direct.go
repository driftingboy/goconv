package goconv

type Typ int

const (
	UnKnow Typ = iota
	Int
	Int32
	Int64
	Float32
	Float64
	String
	Boolean
	Date
	DateTime

	Ints
	Float32s
	Float64s
	Strings
	Booleans
	Dates
	DateTimes
)

var DefaultTypeDirect = NewTypeDirectWithConfig(
	true,
	map[Typ][]string{
		Int:       {"int"},
		Int32:     {"int32"},
		Int64:     {"int64"},
		Boolean:   {"bool", "boolean"},
		Float32:   {"float32"},
		Float64:   {"float64", "double"},
		String:    {"string"},
		Date:      {"date"},
		DateTime:  {"datetime"},
		Ints:      {"ints"},
		Float32s:  {"float32s"},
		Float64s:  {"float64s"},
		Strings:   {"strings"},
		Booleans:  {"bools"},
		Dates:     {"dates"},
		DateTimes: {"datetimes"},
	},
)

type TypeDirect struct {
	isMatchUPDown bool
	typeAlias     map[string]Typ
}

func NewTypeDirect() *TypeDirect {
	return &TypeDirect{}
}

func NewTypeDirectWithConfig(isMatchUPDown bool, typAlias map[Typ][]string) *TypeDirect {
	direct := new(TypeDirect)
	direct.isMatchUPDown = isMatchUPDown

	var length int
	for _, as := range typAlias {
		length += len(as)
	}
	direct.typeAlias = make(map[string]Typ, length)

	for t, as := range typAlias {
		for _, a := range as {
			direct.typeAlias[a] = t
		}
	}

	return direct
}

func (td *TypeDirect) SetIsMatchUPDown(b bool) {
	td.isMatchUPDown = b
}

func (td *TypeDirect) Grow(size int) {
	td.typeAlias = make(map[string]Typ, size)
}

func (td *TypeDirect) AddTypeAliases(t Typ, aliases ...string) {
	for _, a := range aliases {
		td.typeAlias[a] = t
	}
}

func (td TypeDirect) DirectTyp(typ string) Typ {
	t, ok := td.typeAlias[typ]
	if !ok {
		return UnKnow
	}
	return t
}
