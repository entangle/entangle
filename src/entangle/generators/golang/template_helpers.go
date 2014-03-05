package golang

import (
	"entangle/declarations"
	"entangle/utils"
	"fmt"
	"strings"
)

var (
	simpleTypeClassMapping = map[declarations.TypeClass]string{
		declarations.BoolClass:    "bool",
		declarations.StringClass:  "string",
		declarations.Float32Class: "float32",
		declarations.Float64Class: "float64",
		declarations.Int8Class:    "int8",
		declarations.Int16Class:   "int16",
		declarations.Int32Class:   "int32",
		declarations.Int64Class:   "int64",
		declarations.Uint8Class:   "uint8",
		declarations.Uint16Class:  "uint16",
		declarations.Uint32Class:  "uint32",
		declarations.Uint64Class:  "uint64",
	}
	simpleTypeClassDeserializationMapping = map[declarations.TypeClass]string{
		declarations.BoolClass:    "DeserializeBool",
		declarations.StringClass:  "DeserializeString",
		declarations.BinaryClass:  "DeserializeBinary",
		declarations.Float32Class: "DeserializeFloat32",
		declarations.Float64Class: "DeserializeFloat64",
		declarations.Int8Class:    "DeserializeInt8",
		declarations.Int16Class:   "DeserializeInt16",
		declarations.Int32Class:   "DeserializeInt32",
		declarations.Int64Class:   "DeserializeInt64",
		declarations.Uint8Class:   "DeserializeUint8",
		declarations.Uint16Class:  "DeserializeUint16",
		declarations.Uint32Class:  "DeserializeUint32",
		declarations.Uint64Class:  "DeserializeUint64",
	}
)

func documentationHelper(documentation []string, indentation int) string {
	if documentation == nil || len(documentation) == 0 {
		return ""
	}

	prefix := fmt.Sprintf("%s//", strings.Repeat("\t", indentation))
	wrapper := utils.NewSimpleTextWrapper(79 - len(prefix) - 1)
	lines := make([]string, 0, len(documentation)*2)

	for _, paragraph := range documentation {
		if len(lines) > 0 {
			lines = append(lines, "")
		}

		for _, line := range wrapper.Wrap(paragraph) {
			if len(line) > 0 {
				lines = append(lines, fmt.Sprintf("%s %s", prefix, line))
			}
		}
	}

	if len(lines) == 0 {
		return ""
	}

	return fmt.Sprintf("%s\n", strings.Join(lines, "\n"))
}

func typeHelper(typeDecl declarations.Type) string {
	star := ""
	if typeDecl.Nilable() {
		star = "*"
	}

	if simpleName, ok := simpleTypeClassMapping[typeDecl.Class()]; ok {
		return fmt.Sprintf("%s%s", star, simpleName)
	}

	switch typeDecl.Class() {
	case declarations.BinaryClass:
		return "[]byte"

	case declarations.EnumClass:
		return fmt.Sprintf("%s%s", star, typeDecl.(*declarations.EnumType).Enum().Name)

	case declarations.StructClass:
		return fmt.Sprintf("%s%s", star, typeDecl.(*declarations.StructType).Struct().Name)

	case declarations.MapClass:
		mapTypeDecl := typeDecl.(*declarations.MapType)
		return fmt.Sprintf("map[%s]%s", typeHelper(mapTypeDecl.KeyType()), typeHelper(mapTypeDecl.ValueType()))

	case declarations.ListClass:
		listTypeDecl := typeDecl.(*declarations.ListType)
		return fmt.Sprintf("[]%s", typeHelper(listTypeDecl.ElementType()))

	default:
		panic("Unimplemented type")
	}
}

func nonNilableTypeHelper(typeDecl declarations.Type) string {
	name := typeHelper(typeDecl)
	if name[0] == '*' {
		return name[1:]
	}
	return name
}

func canSkipBeforeFieldHelper(fieldDecl *declarations.Field, minimumDeserializedLength int) bool {
	return fieldDecl.Index > uint(minimumDeserializedLength)
}

func indent(input string, indentation int) string {
	split := strings.Split(input, "\n")
	prefix := strings.Repeat("\t", indentation)
	indented := make([]string, len(split))
	for i, s := range split {
		indented[i] = fmt.Sprintf("%s%s", prefix, s)
	}
	return strings.Join(indented, "\n")
}

func typeDeserializationMethodHelper(typeDecl declarations.Type) string {
	if simpleMethod, ok := simpleTypeClassDeserializationMapping[typeDecl.Class()]; ok {
		return fmt.Sprintf("goentangle.%s", simpleMethod)
	}

	switch typeDecl.Class() {
	case declarations.EnumClass:
		return fmt.Sprintf("Deserialize%s", typeDecl.(*declarations.EnumType).Enum().Name)

	case declarations.StructClass:
		return fmt.Sprintf("Deserialize%s", typeDecl.(*declarations.StructType).Struct().Name)

	case declarations.MapClass, declarations.ListClass:
		return nameOfDeserializer(typeDecl)

	default:
		panic("Unimplemented field type")
	}
}

func structFieldDeserializationCodeHelper(structDecl *declarations.Struct, field *declarations.Field, indentation int) string {
	// Determine the method used for deserialization.
	method := typeDeserializationMethodHelper(field.Type)

	// Return the deserialization code.
	if field.Type.Nilable() {
		return indent(fmt.Sprintf(`if des.%s, err = %s(ser[%d]); err != nil {
	if err == goentangle.ErrDeserializationError {
		err = fmt.Errorf("invalid value for field %s in struct %s")
	}
	return
}`, field.Name, method, field.Index-1, field.Name, structDecl.Name), indentation)
	} else {
		return ""
	}
}

func typeSerializationCodeHelper(typeDecl declarations.Type, source, target string, indentation int) string {
	if typeDecl.Nilable() {
		switch typeDecl.Class() {
		case declarations.BoolClass, declarations.StringClass, declarations.BinaryClass, declarations.Float32Class, declarations.Float64Class, declarations.Int8Class, declarations.Int16Class, declarations.Int32Class, declarations.Int64Class, declarations.Uint8Class, declarations.Uint16Class, declarations.Uint32Class, declarations.Uint64Class:
			return indent(fmt.Sprintf(`if %s != nil {
	%s = *%s
}`, source, target, source), indentation)

		case declarations.StructClass, declarations.EnumClass:
			return indent(fmt.Sprintf(`if %s != nil {
	if %s, err = %s.Serialize(); err != nil {
		return
	}
}`, source, target, source), indentation)

		case declarations.ListClass, declarations.MapClass:
			return indent(fmt.Sprintf(`if %s != nil {
	if %s, err = %s(%s); err != nil {
		return
	}
}`, source, target, nameOfSerializer(typeDecl), source), indentation)

		default:
			panic("Unsupport type")
		}
	} else {
		switch typeDecl.Class() {
		case declarations.BoolClass, declarations.StringClass, declarations.BinaryClass, declarations.Float32Class, declarations.Float64Class, declarations.Int8Class, declarations.Int16Class, declarations.Int32Class, declarations.Int64Class, declarations.Uint8Class, declarations.Uint16Class, declarations.Uint32Class, declarations.Uint64Class:
			return indent(fmt.Sprintf(`%s = %s`, target, source), indentation)

		case declarations.StructClass, declarations.EnumClass:
			return indent(fmt.Sprintf(`if %s, err = %s.Serialize(); err != nil {
	return
}`, target, source), indentation)

		case declarations.ListClass, declarations.MapClass:
			return indent(fmt.Sprintf(`if %s == nil {
	err = errors.New("non-nilable type cannot be nil")
	return
}

if %s, err = %s(%s); err != nil {
	return
}`, source, target, nameOfSerializer(typeDecl), source), indentation)

		default:
			panic("Unsupport type")
		}
	}
}

func deserializationCodeHelper(typeDecl declarations.Type) string {
	parts := make([]string, 0, 8)

	switch typeDecl.Class() {
	case declarations.ListClass:
		listDecl := typeDecl.(*declarations.ListType)
		elemType := listDecl.ElementType()

		parts = append(parts, fmt.Sprintf(`	var ser []interface{}
	var serOk bool
	if ser, serOk = input.([]interface{}); !serOk {
		err = goentangle.ErrDeserializationError
		return
	}

	des = make(%s, len(ser))

	for i, serElem := range ser {`, nonNilableTypeHelper(typeDecl)))

		if elemType.Nilable() {
			parts = append(parts, fmt.Sprintf(`		if serElem == nil {
			continue
		}

		var nonNilDes %s
		if nonNilDes, err = %s(serElem); err != nil {
			return
		}
		des[i] = &nonNilDes`, nonNilableTypeHelper(elemType), typeDeserializationMethodHelper(elemType)))
		} else {
			parts = append(parts, fmt.Sprintf(`		if des[i], err = %s(serElem); err != nil {
			return
		}`, typeDeserializationMethodHelper(elemType)))
		}

		parts = append(parts, `	}

	return`)

		return strings.Join(parts, "\n")

	case declarations.MapClass:
		mapDecl := typeDecl.(*declarations.MapType)
		keyType := mapDecl.KeyType()
		valueType := mapDecl.ValueType()

		parts = append(parts, fmt.Sprintf(`	var ser map[interface{}]interface{}
	var serOk bool
	if ser, serOk = input.(map[interface{}]interface{}); !serOk {
		err = goentangle.ErrDeserializationError
		return
	}

	des = make(%s, len(ser))

	for serKey, serValue := range ser {
		var desKey %s
		var desValue %s
`, nonNilableTypeHelper(typeDecl), typeHelper(keyType), typeHelper(valueType)))

		if keyType.Nilable() {
			parts = append(parts, fmt.Sprintf(`		if serKey != nil {
			var nonNilDesKey %s
			if nonNilDesKey, err = %s(serKey); err != nil {
				return
			}
			desKey = &nonNilDesKey
		}`, nonNilableTypeHelper(keyType), typeDeserializationMethodHelper(keyType)))
		} else {
			parts = append(parts, fmt.Sprintf(`		if desKey, err = %s(serKey); err != nil {
			return
		}`, typeDeserializationMethodHelper(keyType)))
		}

		if valueType.Nilable() {
			parts = append(parts, fmt.Sprintf(`		if serValue != nil {
			var nonNilDesValue %s
			if nonNilDesValue, err = %s(serValue); err != nil {
				return
			}
			desValue = &nonNilDesValue
		}`, nonNilableTypeHelper(valueType), typeDeserializationMethodHelper(valueType)))
		} else {
			parts = append(parts, fmt.Sprintf(`		if desValue, err = %s(serValue); err != nil {
			return
		}`, typeDeserializationMethodHelper(valueType)))
		}

		parts = append(parts, `
		des[desKey] = desValue
	}

	return`)

		return strings.Join(parts, "\n")

	default:
		panic("Cannot generate deserialization code for type")
	}
}

func serializationCodeHelper(typeDecl declarations.Type) string {
	switch typeDecl.Class() {
	case declarations.ListClass:
		listDecl := typeDecl.(*declarations.ListType)
		elemType := listDecl.ElementType()

		return fmt.Sprintf(`	serArr := make([]interface{}, len(input))

	for i, des := range input {
%s
	}

	ser = serArr
	return`, typeSerializationCodeHelper(elemType, "des", "serArr[i]", 2))

	case declarations.MapClass:
		mapDecl := typeDecl.(*declarations.MapType)
		keyType := mapDecl.KeyType()
		valueType := mapDecl.ValueType()

		return fmt.Sprintf(`	serMap := make(map[interface{}]interface{}, len(input))

	for desKey, desValue := range input {
		var serKey interface{}
		var serValue interface{}

%s

%s

		serMap[serKey] = serValue
	}

	ser = serMap
	return`, typeSerializationCodeHelper(keyType, "desKey", "serKey", 2), typeSerializationCodeHelper(valueType, "desValue", "serValue", 2))

	default:
		panic("Cannot generate deserialization code for type")
	}
}

func structSerializationCodeHelper(structDecl *declarations.Struct) string {
	parts := make([]string, 0, 2+len(structDecl.Fields))

	parts = append(parts, fmt.Sprintf(`	serArr := make([]interface{}, %d)`, structDecl.SerializedLength()))

	for _, field := range structDecl.FieldsSortedByIndex() {
		parts = append(parts, fmt.Sprintf(`	// Serialize %s.
%s`, field.Name, typeSerializationCodeHelper(field.Type, fmt.Sprintf("s.%s", field.Name), fmt.Sprintf("serArr[%d]", field.Index-1), 1)))
	}

	parts = append(parts, `	ser = serArr
	return`)

	return strings.Join(parts, "\n\n")
}
