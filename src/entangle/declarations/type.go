package declarations

// Type class.
type TypeClass int

// List of type classes.
const (
	BoolClass TypeClass = iota
	StringClass
	BinaryClass
	Float32Class
	Float64Class
	Int8Class
	Int16Class
	Int32Class
	Int64Class
	Uint8Class
	Uint16Class
	Uint32Class
	Uint64Class
	EnumClass
	StructClass
	MapClass
	ListClass
)

// Type declaration.
type Type interface {
	// Base type.
	Class() TypeClass

	// Nilable.
	Nilable() bool
}

// Simple type.
//
// Represents a non-complex type.
type simpleType struct {
	class   TypeClass
	nilable bool
}

func (s *simpleType) Class() TypeClass {
	return s.class
}

func (s *simpleType) Nilable() bool {
	return s.nilable
}

// Struct type.
type StructType struct {
	decl    *Struct
	nilable bool
}

func (s *StructType) Class() TypeClass {
	return StructClass
}

func (s *StructType) Nilable() bool {
	return s.nilable
}

func (s *StructType) Struct() *Struct {
	return s.decl
}

// New struct type.
func NewStructType(decl *Struct, nilable bool) Type {
	return &StructType{
		decl:    decl,
		nilable: nilable,
	}
}

// Enum type.
type EnumType struct {
	decl    *Enum
	nilable bool
}

func (s *EnumType) Class() TypeClass {
	return EnumClass
}

func (s *EnumType) Nilable() bool {
	return s.nilable
}

func (s *EnumType) Enum() *Enum {
	return s.decl
}

// New enum type.
func NewEnumType(decl *Enum, nilable bool) Type {
	return &EnumType{
		decl:    decl,
		nilable: nilable,
	}
}

// List type.
type ListType struct {
	elementType Type
	nilable     bool
}

func (s *ListType) Class() TypeClass {
	return ListClass
}

func (s *ListType) ElementType() Type {
	return s.elementType
}

func (s *ListType) Nilable() bool {
	return s.nilable
}

// New list type.
func NewListType(elementType Type, nilable bool) Type {
	return &ListType{
		elementType: elementType,
		nilable:     nilable,
	}
}

// Map type.
type MapType struct {
	keyType   Type
	valueType Type
	nilable   bool
}

func (s *MapType) Class() TypeClass {
	return MapClass
}

func (s *MapType) KeyType() Type {
	return s.keyType
}

func (s *MapType) ValueType() Type {
	return s.valueType
}

func (s *MapType) Nilable() bool {
	return s.nilable
}

// New map type.
func NewMapType(keyType, valueType Type, nilable bool) Type {
	return &MapType{
		keyType:   keyType,
		valueType: valueType,
		nilable:   nilable,
	}
}

var (
	BoolType           = &simpleType{BoolClass, false}
	NilableBoolType    = &simpleType{BoolClass, true}
	StringType         = &simpleType{StringClass, false}
	NilableStringType  = &simpleType{StringClass, true}
	BinaryType         = &simpleType{BinaryClass, false}
	NilableBinaryType  = &simpleType{BinaryClass, true}
	Float32Type        = &simpleType{Float32Class, false}
	NilableFloat32Type = &simpleType{Float32Class, true}
	Float64Type        = &simpleType{Float64Class, false}
	NilableFloat64Type = &simpleType{Float64Class, true}
	Int8Type           = &simpleType{Int8Class, false}
	NilableInt8Type    = &simpleType{Int8Class, true}
	Int16Type          = &simpleType{Int16Class, false}
	NilableInt16Type   = &simpleType{Int16Class, true}
	Int32Type          = &simpleType{Int32Class, false}
	NilableInt32Type   = &simpleType{Int32Class, true}
	Int64Type          = &simpleType{Int64Class, false}
	NilableInt64Type   = &simpleType{Int64Class, true}
	Uint8Type          = &simpleType{Uint8Class, false}
	NilableUint8Type   = &simpleType{Uint8Class, true}
	Uint16Type         = &simpleType{Uint16Class, false}
	NilableUint16Type  = &simpleType{Uint16Class, true}
	Uint32Type         = &simpleType{Uint32Class, false}
	NilableUint32Type  = &simpleType{Uint32Class, true}
	Uint64Type         = &simpleType{Uint64Class, false}
	NilableUint64Type  = &simpleType{Uint64Class, true}
)
