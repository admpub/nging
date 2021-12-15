package confl

// represents any Go type that corresponds to a internal type.
type confType interface {
	typeString() string
}

// typeEqual accepts any two types and returns true if they are equal.
func typeEqual(t1, t2 confType) bool {
	if t1 == nil || t2 == nil {
		return false
	}
	return t1.typeString() == t2.typeString()
}

func typeIsHash(t confType) bool {
	return typeEqual(t, confHash) || typeEqual(t, confArrayHash)
}

type confBaseType string

func (btype confBaseType) typeString() string {
	return string(btype)
}

func (btype confBaseType) String() string {
	return btype.typeString()
}

var (
	confInteger   confBaseType = "Integer"
	confFloat     confBaseType = "Float"
	confDatetime  confBaseType = "Datetime"
	confString    confBaseType = "String"
	confBool      confBaseType = "Bool"
	confArray     confBaseType = "Array"
	confHash      confBaseType = "Hash"
	confArrayHash confBaseType = "ArrayHash"
)

// typeOfPrimitive returns a confType of any primitive value in conf.
// Primitive values are: Integer, Float, Datetime, String and Bool.
//
// Passing a lexer item other than the following will cause a BUG message
// to occur: itemString, itemBool, itemInteger, itemFloat, itemDatetime.
func (p *parser) typeOfPrimitive(lexItem item) confType {
	switch lexItem.typ {
	case itemInteger:
		return confInteger
	case itemFloat:
		return confFloat
	case itemDatetime:
		return confDatetime
	case itemString:
		return confString
	case itemBool:
		return confBool
	}
	p.bug("Cannot infer primitive type of lex item '%s'.", lexItem)
	panic("unreachable")
}

// typeOfArray returns a confType for an array given a list of types of its
// values.
//
// In the current spec, if an array is homogeneous, then its type is always
// "Array". If the array is not homogeneous, an error is generated.
func (p *parser) typeOfArray(types []confType) confType {
	// Empty arrays are cool.
	if len(types) == 0 {
		return confArray
	}

	theType := types[0]
	for _, t := range types[1:] {
		if !typeEqual(theType, t) {
			p.panicf("Array contains values of type '%s' and '%s', but arrays "+
				"must be homogeneous.", theType, t)
		}
	}
	return confArray
}
