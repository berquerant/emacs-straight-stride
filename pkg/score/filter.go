package score

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type Filter struct {
	program *vm.Program
}

func NewFilter(text string) (*Filter, error) {
	p, err := expr.Compile(text, expr.AsBool())
	if err != nil {
		return nil, err
	}
	return &Filter{
		program: p,
	}, nil
}

func (f *Filter) Select(x *Log) (bool, error) {
	v, err := expr.Run(f.program, x)
	if err != nil {
		return false, fmt.Errorf("%w: failed to filter log, %v", err, x)
	}
	return v.(bool), nil
}

const DefaultFilterFormula = "true"
