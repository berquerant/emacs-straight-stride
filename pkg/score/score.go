package score

import (
	"fmt"

	"github.com/berquerant/grinfo"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type Calculator struct {
	program *vm.Program
}

func NewCalculator(text string) (*Calculator, error) {
	p, err := expr.Compile(text, expr.AsFloat64())
	if err != nil {
		return nil, err
	}
	return &Calculator{
		program: p,
	}, nil
}

func (c *Calculator) Calculate(x *grinfo.Log) (*Log, error) {
	y := NewLog(x)
	v, err := expr.Run(c.program, y)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to calculate score, %v", err, x)
	}
	y.Score = v.(float64)
	return y, nil
}

const DefaultCalculatorFormula = "diff_commit + 3 * diff_tag + 3 * diff_day"
