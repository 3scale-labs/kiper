package threescale

import (
	"fmt"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/topdown"
	"github.com/open-policy-agent/opa/topdown/builtins"
)

func RegisterThreeScaleQueries() {
	ast.RegisterBuiltin(PrintPathBuiltin)
	topdown.RegisterFunctionalBuiltin1(PrintPathBuiltin.Name, PrintPathImpl)
}

func RegisterRateLimitQueries() {
	registerRateLimitBuiltIn()
	registerUpdateLimitsUsageBuiltin()
}

func registerRateLimitBuiltIn() {
	ast.RegisterBuiltin(rateLimitBuiltin)

	name := rateLimitBuiltin.Name
	funcImpl := rateLimitBuiltinImpl

	// We can't use the helpers provided by OPA directly because we need to pass
	// the builtinContext as a param. We use that to store the counters that
	// need to be updated at the end of the query.
	builtinFunc := func(bctx topdown.BuiltinContext, args []*ast.Term, iter func(*ast.Term) error) error {
		result, err := funcImpl(args[0].Value, bctx)
		if err == nil {
			return iter(ast.NewTerm(result))
		}
		if _, empty := err.(topdown.BuiltinEmpty); empty {
			return nil
		}
		return handleBuiltinErr(name, bctx.Location, err)
	}

	topdown.RegisterBuiltinFunc(name, builtinFunc)
}

func registerUpdateLimitsUsageBuiltin() {
	ast.RegisterBuiltin(updateLimitsUsageBuiltin)

	name := updateLimitsUsageBuiltin.Name
	funcImpl := updateLimitsUsageBuiltinImpl

	builtinFunc := func(bctx topdown.BuiltinContext, args []*ast.Term, iter func(*ast.Term) error) error {
		result, err := funcImpl(bctx)
		if err == nil {
			return iter(ast.NewTerm(result))
		}
		if _, empty := err.(topdown.BuiltinEmpty); empty {
			return nil
		}
		return handleBuiltinErr(name, bctx.Location, err)
	}

	topdown.RegisterBuiltinFunc(name, builtinFunc)
}

// This is copy-pasted from the OPA library. It's a private func that we need to
// define the rate-limit commands.
func handleBuiltinErr(name string, loc *ast.Location, err error) error {
	switch err := err.(type) {
	case topdown.BuiltinEmpty:
		return nil
	case builtins.ErrOperand:
		return &topdown.Error{
			Code:     topdown.TypeErr,
			Message:  fmt.Sprintf("%v: %v", string(name), err.Error()),
			Location: loc,
		}
	default:
		return &topdown.Error{
			Code:     topdown.BuiltinErr,
			Message:  fmt.Sprintf("%v: %v", string(name), err.Error()),
			Location: loc,
		}
	}
}
