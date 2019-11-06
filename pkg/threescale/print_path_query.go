// Example custom OPA query that returns the path of a request

package threescale

import (
	"encoding/json"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/types"
)

var PrintPathBuiltin = &ast.Builtin{
	Name: "printpath",
	Decl: types.NewFunction(types.Args(types.A), types.S),
}

func PrintPathImpl(a ast.Value) (ast.Value, error) {
	input := a.String()
	request := Input{}
	_ = json.Unmarshal([]byte(input), &request)
	return ast.String(request.Attributes.Request.HTTP.Path), nil
}
