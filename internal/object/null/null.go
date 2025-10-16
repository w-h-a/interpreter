package null

import "github.com/w-h-a/interpreter/internal/object"

type Null struct{}

func (o *Null) Inspect() string {
	return "null"
}

func (o *Null) Type() object.ObjectType {
	return object.NULL
}
