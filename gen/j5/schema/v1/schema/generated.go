package schema

// Code generated by jsonapi. DO NOT EDIT.
// Source: github.com/pentops/o5-aws-tool/gen/j5/schema/v1/schema

import (
	list "github.com/pentops/o5-aws-tool/gen/j5/list/v1/list"
)

// EnumField_Ext Proto: EnumField_Ext
type EnumField_Ext struct {
}

// BoolField_Rules Proto: BoolField_Rules
type BoolField_Rules struct {
	Const *bool `json:"const,omitempty"`
}

// Enum_OptionInfoField Proto: Enum_OptionInfoField
type Enum_OptionInfoField struct {
	Name        string `json:"name,omitempty"`
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`
}

// EnumField Proto: EnumField
type EnumField struct {
	Ref       *Ref             `json:"ref,omitempty"`
	Enum      *Enum            `json:"enum,omitempty"`
	Rules     *EnumField_Rules `json:"rules,omitempty"`
	ListRules *list.EnumRules  `json:"listRules,omitempty"`
	Ext       *EnumField_Ext   `json:"ext,omitempty"`
}

// EntityRef Proto: EntityRef
type EntityRef struct {
	Package string `json:"package,omitempty"`
	Entity  string `json:"entity,omitempty"`
}

// EnumField_Rules Proto: EnumField_Rules
type EnumField_Rules struct {
	In    []string `json:"in,omitempty"`
	NotIn []string `json:"notIn,omitempty"`
}

// Field Proto Oneof: j5.schema.v1.Field
type Field struct {
	J5TypeKey string          `json:"!type,omitempty"`
	Any       *AnyField       `json:"any,omitempty"`
	Oneof     *OneofField     `json:"oneof,omitempty"`
	Object    *ObjectField    `json:"object,omitempty"`
	Enum      *EnumField      `json:"enum,omitempty"`
	Array     *ArrayField     `json:"array,omitempty"`
	Map       *MapField       `json:"map,omitempty"`
	String    *StringField    `json:"string,omitempty"`
	Integer   *IntegerField   `json:"integer,omitempty"`
	Float     *FloatField     `json:"float,omitempty"`
	Bool      *BoolField      `json:"bool,omitempty"`
	Bytes     *BytesField     `json:"bytes,omitempty"`
	Decimal   *DecimalField   `json:"decimal,omitempty"`
	Date      *DateField      `json:"date,omitempty"`
	Timestamp *TimestampField `json:"timestamp,omitempty"`
	Key       *KeyField       `json:"key,omitempty"`
}

func (s Field) OneofKey() string {
	if s.Any != nil {
		return "any"
	}
	if s.Oneof != nil {
		return "oneof"
	}
	if s.Object != nil {
		return "object"
	}
	if s.Enum != nil {
		return "enum"
	}
	if s.Array != nil {
		return "array"
	}
	if s.Map != nil {
		return "map"
	}
	if s.String != nil {
		return "string"
	}
	if s.Integer != nil {
		return "integer"
	}
	if s.Float != nil {
		return "float"
	}
	if s.Bool != nil {
		return "bool"
	}
	if s.Bytes != nil {
		return "bytes"
	}
	if s.Decimal != nil {
		return "decimal"
	}
	if s.Date != nil {
		return "date"
	}
	if s.Timestamp != nil {
		return "timestamp"
	}
	if s.Key != nil {
		return "key"
	}
	return ""
}

func (s Field) Type() interface{} {
	if s.Any != nil {
		return s.Any
	}
	if s.Oneof != nil {
		return s.Oneof
	}
	if s.Object != nil {
		return s.Object
	}
	if s.Enum != nil {
		return s.Enum
	}
	if s.Array != nil {
		return s.Array
	}
	if s.Map != nil {
		return s.Map
	}
	if s.String != nil {
		return s.String
	}
	if s.Integer != nil {
		return s.Integer
	}
	if s.Float != nil {
		return s.Float
	}
	if s.Bool != nil {
		return s.Bool
	}
	if s.Bytes != nil {
		return s.Bytes
	}
	if s.Decimal != nil {
		return s.Decimal
	}
	if s.Date != nil {
		return s.Date
	}
	if s.Timestamp != nil {
		return s.Timestamp
	}
	if s.Key != nil {
		return s.Key
	}
	return nil
}

// BytesField_Ext Proto: BytesField_Ext
type BytesField_Ext struct {
}

// EntityKey Proto: EntityKey
type EntityKey struct {
	PrimaryKey bool       `json:"primaryKey,omitempty"`
	ForeignKey *EntityRef `json:"foreignKey,omitempty"`
}

// DateField_Ext Proto: DateField_Ext
type DateField_Ext struct {
}

// Enum_Option Proto: Enum_Option
type Enum_Option struct {
	Name        string            `json:"name,omitempty"`
	Number      int32             `json:"number,omitempty"`
	Description string            `json:"description,omitempty"`
	Info        map[string]string `json:"info,omitempty"`
}

// FloatField_Format Proto Enum: j5.schema.v1.FloatField_Format
type FloatField_Format string

const (
	FloatField_Format_UNSPECIFIED FloatField_Format = "UNSPECIFIED"
	FloatField_Format_FLOAT32     FloatField_Format = "FLOAT32"
	FloatField_Format_FLOAT64     FloatField_Format = "FLOAT64"
)

// DecimalField_Ext Proto: DecimalField_Ext
type DecimalField_Ext struct {
}

// TimestampField_Ext Proto: TimestampField_Ext
type TimestampField_Ext struct {
}

// ObjectField_EntityJoin Proto: ObjectField_EntityJoin
type ObjectField_EntityJoin struct {
	Entity     *EntityRef `json:"entity,omitempty"`
	EntityPart string     `json:"entityPart,omitempty"`
}

// ArrayField_Ext Proto: ArrayField_Ext
type ArrayField_Ext struct {
	SingleForm *string `json:"singleForm,omitempty"`
}

// IntegerField Proto: IntegerField
type IntegerField struct {
	Format    string              `json:"format"`
	Rules     *IntegerField_Rules `json:"rules,omitempty"`
	ListRules *list.IntegerRules  `json:"listRules,omitempty"`
	Ext       *IntegerField_Ext   `json:"ext,omitempty"`
}

// FloatField_Ext Proto: FloatField_Ext
type FloatField_Ext struct {
}

// BoolField_Ext Proto: BoolField_Ext
type BoolField_Ext struct {
}

// KeyField Proto: KeyField
type KeyField struct {
	Rules     *KeyField_Rules `json:"rules,omitempty"`
	Format    *KeyFormat      `json:"format,omitempty"`
	ListRules *list.KeyRules  `json:"listRules,omitempty"`
	Ext       *KeyField_Ext   `json:"ext,omitempty"`
	Entity    *EntityKey      `json:"entity,omitempty"`
}

// BytesField_Rules Proto: BytesField_Rules
type BytesField_Rules struct {
	MinLength *uint64 `json:"minLength,omitempty"`
	MaxLength *uint64 `json:"maxLength,omitempty"`
}

// AnyField Proto: AnyField
type AnyField struct {
	OnlyDefined bool     `json:"onlyDefined,omitempty"`
	Types       []string `json:"types,omitempty"`
}

// ObjectField Proto: ObjectField
type ObjectField struct {
	Ref     *Ref                    `json:"ref,omitempty"`
	Object  *Object                 `json:"object,omitempty"`
	Rules   *ObjectField_Rules      `json:"rules,omitempty"`
	Ext     *ObjectField_Ext        `json:"ext,omitempty"`
	Flatten bool                    `json:"flatten,omitempty"`
	Entity  *ObjectField_EntityJoin `json:"entity,omitempty"`
}

// TimestampField Proto: TimestampField
type TimestampField struct {
	Rules     *TimestampField_Rules `json:"rules,omitempty"`
	ListRules *list.TimestampRules  `json:"listRules,omitempty"`
	Ext       *TimestampField_Ext   `json:"ext,omitempty"`
}

// ObjectField_Rules Proto: ObjectField_Rules
type ObjectField_Rules struct {
	MinProperties *uint64 `json:"minProperties,omitempty"`
	MaxProperties *uint64 `json:"maxProperties,omitempty"`
}

// EntityPart Proto Enum: j5.schema.v1.EntityPart
type EntityPart string

const (
	EntityPart_UNSPECIFIED EntityPart = "UNSPECIFIED"
	EntityPart_KEYS        EntityPart = "KEYS"
	EntityPart_STATE       EntityPart = "STATE"
	EntityPart_EVENT       EntityPart = "EVENT"
	EntityPart_DATA        EntityPart = "DATA"
	EntityPart_REFERENCES  EntityPart = "REFERENCES"
	EntityPart_DERIVED     EntityPart = "DERIVED"
)

// KeyFormat_UUID Proto: KeyFormat_UUID
type KeyFormat_UUID struct {
}

// OneofField Proto: OneofField
type OneofField struct {
	Ref       *Ref              `json:"ref,omitempty"`
	Oneof     *Oneof            `json:"oneof,omitempty"`
	Rules     *OneofField_Rules `json:"rules,omitempty"`
	ListRules *list.OneofRules  `json:"listRules,omitempty"`
	Ext       *OneofField_Ext   `json:"ext,omitempty"`
}

// DateField Proto: DateField
type DateField struct {
	Rules *DateField_Rules `json:"rules,omitempty"`
	Ext   *DateField_Ext   `json:"ext,omitempty"`
}

// DecimalField_Rules Proto: DecimalField_Rules
type DecimalField_Rules struct {
}

// FloatField_Rules Proto: FloatField_Rules
type FloatField_Rules struct {
	ExclusiveMaximum *bool    `json:"exclusiveMaximum,omitempty"`
	ExclusiveMinimum *bool    `json:"exclusiveMinimum,omitempty"`
	Minimum          *float64 `json:"minimum,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty"`
	MultipleOf       *float64 `json:"multipleOf,omitempty"`
}

// ObjectProperty Proto: ObjectProperty
type ObjectProperty struct {
	Schema             *Field  `json:"schema,omitempty"`
	Name               string  `json:"name,omitempty"`
	Required           bool    `json:"required,omitempty"`
	ExplicitlyOptional bool    `json:"explicitlyOptional,omitempty"`
	Description        string  `json:"description,omitempty"`
	ProtoField         []int32 `json:"protoField,omitempty"`
}

// MapField_Rules Proto: MapField_Rules
type MapField_Rules struct {
	MinPairs *uint64 `json:"minPairs,omitempty"`
	MaxPairs *uint64 `json:"maxPairs,omitempty"`
}

// Enum Proto: Enum
type Enum struct {
	Name        string                  `json:"name,omitempty"`
	Description string                  `json:"description,omitempty"`
	Prefix      string                  `json:"prefix,omitempty"`
	Options     []*Enum_Option          `json:"options,omitempty"`
	Info        []*Enum_OptionInfoField `json:"info,omitempty"`
}

// FloatField Proto: FloatField
type FloatField struct {
	Format    string            `json:"format,omitempty"`
	Rules     *FloatField_Rules `json:"rules,omitempty"`
	ListRules *list.FloatRules  `json:"listRules,omitempty"`
	Ext       *FloatField_Ext   `json:"ext,omitempty"`
}

// KeyFormat_Custom Proto: KeyFormat_Custom
type KeyFormat_Custom struct {
	Pattern string `json:"pattern"`
}

// OneofField_Ext Proto: OneofField_Ext
type OneofField_Ext struct {
}

// EntityObject Proto: EntityObject
type EntityObject struct {
	Entity string `json:"entity,omitempty"`
	Part   string `json:"part,omitempty"`
}

// Object Proto: Object
type Object struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Entity      *EntityObject     `json:"entity,omitempty"`
	Properties  []*ObjectProperty `json:"properties,omitempty"`
	AnyMember   []string          `json:"anyMember,omitempty"`
}

// TimestampField_Rules Proto: TimestampField_Rules
type TimestampField_Rules struct {
}

// DateField_Rules Proto: DateField_Rules
type DateField_Rules struct {
}

// KeyFormat_ID62 Proto: KeyFormat_ID62
type KeyFormat_ID62 struct {
}

// KeyField_Rules Proto: KeyField_Rules
type KeyField_Rules struct {
}

// RootSchema Proto Oneof: j5.schema.v1.RootSchema
type RootSchema struct {
	J5TypeKey string  `json:"!type,omitempty"`
	Oneof     *Oneof  `json:"oneof,omitempty"`
	Object    *Object `json:"object,omitempty"`
	Enum      *Enum   `json:"enum,omitempty"`
}

func (s RootSchema) OneofKey() string {
	if s.Oneof != nil {
		return "oneof"
	}
	if s.Object != nil {
		return "object"
	}
	if s.Enum != nil {
		return "enum"
	}
	return ""
}

func (s RootSchema) Type() interface{} {
	if s.Oneof != nil {
		return s.Oneof
	}
	if s.Object != nil {
		return s.Object
	}
	if s.Enum != nil {
		return s.Enum
	}
	return nil
}

// StringField_Ext Proto: StringField_Ext
type StringField_Ext struct {
}

// IntegerField_Rules Proto: IntegerField_Rules
type IntegerField_Rules struct {
	ExclusiveMaximum *bool  `json:"exclusiveMaximum,omitempty"`
	ExclusiveMinimum *bool  `json:"exclusiveMinimum,omitempty"`
	Minimum          *int64 `json:"minimum,omitempty"`
	Maximum          *int64 `json:"maximum,omitempty"`
	MultipleOf       *int64 `json:"multipleOf,omitempty"`
}

// OneofField_Rules Proto: OneofField_Rules
type OneofField_Rules struct {
}

// StringField Proto: StringField
type StringField struct {
	Format    *string             `json:"format,omitempty"`
	Rules     *StringField_Rules  `json:"rules,omitempty"`
	ListRules *list.OpenTextRules `json:"listRules,omitempty"`
	Ext       *StringField_Ext    `json:"ext,omitempty"`
}

// BytesField Proto: BytesField
type BytesField struct {
	Rules *BytesField_Rules `json:"rules,omitempty"`
	Ext   *BytesField_Ext   `json:"ext,omitempty"`
}

// ObjectField_Ext Proto: ObjectField_Ext
type ObjectField_Ext struct {
}

// Ref Proto: Ref
type Ref struct {
	Package string `json:"package,omitempty"`
	Schema  string `json:"schema,omitempty"`
}

// KeyField_Ext Proto: KeyField_Ext
type KeyField_Ext struct {
}

// ArrayField Proto: ArrayField
type ArrayField struct {
	Rules *ArrayField_Rules `json:"rules,omitempty"`
	Items *Field            `json:"items,omitempty"`
	Ext   *ArrayField_Ext   `json:"ext,omitempty"`
}

// DecimalField Proto: DecimalField
type DecimalField struct {
	Rules *DecimalField_Rules `json:"rules,omitempty"`
	Ext   *DecimalField_Ext   `json:"ext,omitempty"`
}

// IntegerField_Format Proto Enum: j5.schema.v1.IntegerField_Format
type IntegerField_Format string

const (
	IntegerField_Format_UNSPECIFIED IntegerField_Format = "UNSPECIFIED"
	IntegerField_Format_INT32       IntegerField_Format = "INT32"
	IntegerField_Format_INT64       IntegerField_Format = "INT64"
	IntegerField_Format_UINT32      IntegerField_Format = "UINT32"
	IntegerField_Format_UINT64      IntegerField_Format = "UINT64"
)

// BoolField Proto: BoolField
type BoolField struct {
	Rules     *BoolField_Rules `json:"rules,omitempty"`
	ListRules *list.BoolRules  `json:"listRules,omitempty"`
	Ext       *BoolField_Ext   `json:"ext,omitempty"`
}

// Oneof Proto: Oneof
type Oneof struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Properties  []*ObjectProperty `json:"properties,omitempty"`
}

// KeyFormat Proto Oneof: j5.schema.v1.KeyFormat
type KeyFormat struct {
	J5TypeKey string              `json:"!type,omitempty"`
	Informal  *KeyFormat_Informal `json:"informal,omitempty"`
	Custom    *KeyFormat_Custom   `json:"custom,omitempty"`
	Uuid      *KeyFormat_UUID     `json:"uuid,omitempty"`
	Id62      *KeyFormat_ID62     `json:"id62,omitempty"`
}

func (s KeyFormat) OneofKey() string {
	if s.Informal != nil {
		return "informal"
	}
	if s.Custom != nil {
		return "custom"
	}
	if s.Uuid != nil {
		return "uuid"
	}
	if s.Id62 != nil {
		return "id62"
	}
	return ""
}

func (s KeyFormat) Type() interface{} {
	if s.Informal != nil {
		return s.Informal
	}
	if s.Custom != nil {
		return s.Custom
	}
	if s.Uuid != nil {
		return s.Uuid
	}
	if s.Id62 != nil {
		return s.Id62
	}
	return nil
}

// MapField_Ext Proto: MapField_Ext
type MapField_Ext struct {
	SingleForm *string `json:"singleForm,omitempty"`
}

// MapField Proto: MapField
type MapField struct {
	ItemSchema *Field          `json:"itemSchema,omitempty"`
	KeySchema  *Field          `json:"keySchema,omitempty"`
	Rules      *MapField_Rules `json:"rules,omitempty"`
	Ext        *MapField_Ext   `json:"ext,omitempty"`
}

// IntegerField_Ext Proto: IntegerField_Ext
type IntegerField_Ext struct {
}

// ArrayField_Rules Proto: ArrayField_Rules
type ArrayField_Rules struct {
	MinItems    *uint64 `json:"minItems,omitempty"`
	MaxItems    *uint64 `json:"maxItems,omitempty"`
	UniqueItems *bool   `json:"uniqueItems,omitempty"`
}

// KeyFormat_Informal Proto: KeyFormat_Informal
type KeyFormat_Informal struct {
}

// StringField_Rules Proto: StringField_Rules
type StringField_Rules struct {
	Pattern   *string `json:"pattern,omitempty"`
	MinLength *uint64 `json:"minLength,omitempty"`
	MaxLength *uint64 `json:"maxLength,omitempty"`
}
