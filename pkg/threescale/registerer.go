package threescale

import (
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/topdown"
)

func RegisterThreeScaleQueries() {
	ast.RegisterBuiltin(PrintPathBuiltin)
	topdown.RegisterFunctionalBuiltin1(PrintPathBuiltin.Name, PrintPathImpl)
}
