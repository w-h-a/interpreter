package integer

import (
	"fmt"

	"github.com/w-h-a/interpreter/internal/evaluator/object"
)

type Integer struct {
	Value int64
}

func (o *Integer) Inspect() string {
	return fmt.Sprintf("%d", o.Value)
}

func (o *Integer) Type() object.ObjectType {
	return object.INTEGER
}
