package models

// DataType represents a data model
type DataType struct {
	Name   Name
	Fields []DataField
}

// DataField represents a data model's field
type DataField struct {
	Name                  Name
	Type                  string
	Pointer               bool
	ValidForCreationInput bool
	ValidForUpdateInput   bool
}

// Name is a handy struct for all the things we need from a data types name
type Name struct {
	Singular, // singular title cased
	RouteName, // usually snaked case and singular
	PluralRouteName, // usually snaked case and plural
	UnexportedVarName, // singular camelCased
	PluralUnexportedVarName, // plural camelCased
	Plural, // plural title cased
	_ string
}
