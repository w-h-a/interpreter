package boolean

import (
	"fmt"

	"github.com/w-h-a/interpreter/internal/object"
)

type Boolean struct {
	Value bool
}

func (o *Boolean) Inspect() string {
	return fmt.Sprintf("%t", o.Value)
}

func (o *Boolean) Type() object.ObjectType {
	return object.BOOLEAN
}
