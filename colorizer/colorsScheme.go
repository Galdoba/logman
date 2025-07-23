package colorizer

type colorSchema struct {
	color256 map[ColorKey]uint8
}

func New(customKeys ...colorData) *colorSchema {
	c := colorSchema{}
	c.color256 = make(map[ColorKey]uint8)
	// c.color256 = defaultColorKey()
	for _, custom := range customKeys {
		c.color256[custom.key] = custom.val
	}
	return &c
}

func DefaultScheme() *colorSchema {
	c := colorSchema{}
	c.color256 = make(map[ColorKey]uint8)
	c.color256 = defaultColorKey()
	return &c
}

func (c *colorSchema) WithColors(colors ...colorData) *colorSchema {
	for _, custom := range colors {
		c.color256[custom.key] = custom.val
	}
	return c
}

type colorData struct {
	key ColorKey
	val uint8
}

func CustomColor(key ColorKey, color256Value uint8) colorData {
	return colorData{
		key: key,
		val: color256Value,
	}
}

func defaultColorKey() map[ColorKey]uint8 {
	colMap := make(map[ColorKey]uint8)
	colMap[fgKey("base")] = 7
	colMap[bgKey("base")] = 0

	colMap[fgKey("string")] = 208
	colMap[fgKey("byte")] = 95
	colMap[fgKey("rune")] = 95
	colMap[fgKey("int")] = 120
	colMap[fgKey("int8")] = 120
	colMap[fgKey("int16")] = 120
	colMap[fgKey("int32")] = 120
	colMap[fgKey("int64")] = 120
	colMap[fgKey("float32")] = 9
	colMap[fgKey("float64")] = 9
	colMap[fgKey("bool")] = 12

	colMap[fgKey("struct")] = 221
	colMap[fgKey("slice")] = 14 //248
	colMap[fgKey("interface")] = 2
	colMap[fgKey("nil")] = 12
	colMap[fgKey("map")] = 14  //207
	colMap[fgKey("ptr")] = 221 //207
	colMap[fgKey("func")] = 36 //207
	colMap[fgKey("chan")] = 2  //207

	colMap[fgKey("fatal")] = 88
	colMap[fgKey("error")] = 196
	colMap[fgKey("warn")] = 184
	colMap[fgKey("report")] = 40
	colMap[fgKey("info")] = 112
	colMap[fgKey("debug")] = 244
	colMap[fgKey("trace")] = 230

	colMap[fgKey("caller")] = 244
	return colMap
}

type ColorKey struct {
	keytype string //field/fg/bg
	value   string
}

const (
	FIELD_KEY = "field"
	FG_KEY    = "fg"
	BG_KEY    = "bg"
)

// fieldKey - дает цвет покраски для поля (например '[error]' - может быть полностью красным в независимости от типа переменной)
func fldKey(val string) ColorKey {
	return ColorKey{
		keytype: FIELD_KEY,
		value:   val,
	}
}

func fgKey(val string) ColorKey {
	return ColorKey{
		keytype: FG_KEY,
		value:   val,
	}
}

func bgKey(val string) ColorKey {
	return ColorKey{
		keytype: BG_KEY,
		value:   val,
	}
}

func NewKey(keyType, value string) ColorKey {
	return ColorKey{
		keytype: keyType,
		value:   value,
	}
}

func (c *colorSchema) getColor(key ColorKey) uint8 {
	if v, ok := c.color256[key]; ok {
		return v
	}
	switch key.keytype {
	case FG_KEY:
		return 7
	case BG_KEY:
		return 0
	}
	return 10
}

func colorToField(c *colorSchema, kind, keyType string) uint8 {
	// switch kind {
	// case "string", "bool",
	// 	"int", "int8", "int16", "int32", "int64",
	// 	"Int",
	// 	"uint", "uint8", "uint16", "uint32", "uint64",
	// 	"float32", "float64":
	// 	return c.getColor(NewKey(keyType, kind))
	// case "struct", "map", "slice", "interface", "ptr", "func", "chan", "nil":
	// 	return c.getColor(NewKey(keyType, kind))
	// }

	switch keyType {
	case FG_KEY:
		return c.getColor(NewKey(FG_KEY, kind))
	}
	return c.getColor(NewKey(BG_KEY, kind))
}
