package object

type ObjectType string

const (
	INTEGER ObjectType = "INTEGER"
	BOOLEAN ObjectType = "BOOLEAN"
	NULL    ObjectType = "NULL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}
