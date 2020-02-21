package threescale

import (
	"encoding/json"
	"time"

	"github.com/open-policy-agent/opa/topdown"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/types"
)

// There are two commands defined for rate limiting: "rate_limit" and
// "update_limits_usage".
// "rate_limit" defines a limit for some entity or combination of them (source
// IP, user ID, request path, etc.) and checks whether the request should be
// limited.
// "update_limits_usage" updates all the rate-limit counters associated with the
// request.
// The counters cannot be updated in the "rate_limit" command because there
// might be several limits that apply for a single request, and we want to
// update the usage only when we are within the delimits defined for all of
// them.
// Ideally, we'd like to run "update_limits_usage" transparently, but for now,
// we need to call it explicitly in the .rego after checking the limits.

type opaLimit struct {
	By      map[string]string
	Count   int
	Seconds int
}

// The string generated in this function should be unique per opaLimit.
func (limit *opaLimit) key() (string, error) {
	limitJSON, err := json.Marshal(&limit)

	if err != nil {
		return "", err
	}

	return string(limitJSON), nil
}

var storage = newLimitsStorage()

const limitsCtxKey = "limits_to_apply"

var rateLimitBuiltin = &ast.Builtin{
	Name: "rate_limit",
	Decl: types.NewFunction(types.Args(types.A), types.B),
}

// Returns true if limited, false otherwise. It does not update the counter used
// to rate-limit. Instead, when within limits, the increase value is stored in
// the context so it can be applied later. When rate-limited, it cleans the
// limits in the context, because in that case, we want to make sure that no
// updates are applied.
func rateLimitBuiltinImpl(limitValue ast.Value, bctx topdown.BuiltinContext) (ast.Value, error) {
	limit := opaLimit{}
	err := ast.As(limitValue, &limit)
	if err != nil {
		return nil, err
	}

	// Note: there's a hardcoded count of 1, we might want to change this later
	// and allow arbitrary updates on each request.
	withinLimits, err := withinLimits(1, &limit)
	if err != nil {
		return nil, err
	}

	if withinLimits {
		err := storeLimitUpdateInCtx(&limit, 1, bctx)
		if err != nil {
			return nil, err
		}
	} else {
		removeLimitUpdates(bctx)
	}

	return ast.Boolean(!withinLimits), nil
}

var updateLimitsUsageBuiltin = &ast.Builtin{
	Name: "update_limits_usage",
	Decl: types.NewFunction(types.Args(), types.B),
}

// Updates the counters used to rate-limit with the values stored in the
// context.
func updateLimitsUsageBuiltinImpl(bctx topdown.BuiltinContext) (ast.Value, error) {
	limits, exists := bctx.Cache.Get(limitsCtxKey)

	if !exists || limits == nil {
		return ast.Boolean(true), nil
	}

	for limitJSON, increaseBy := range limits.(map[string]int) {
		limit := opaLimit{}
		err := json.Unmarshal([]byte(limitJSON), &limit)
		if err != nil {
			return nil, err
		}

		err = updateUsage(increaseBy, &limit)
		if err != nil {
			return nil, err
		}
	}

	return ast.Boolean(true), nil
}

// Returns true if the limit if we do not go over limits after adding count.
// False otherwise.
func withinLimits(count int, limit *opaLimit) (bool, error) {
	key, err := limit.key()

	if err != nil {
		return false, err
	}

	currentVal, exists, err := storage.get(key)
	if err != nil {
		return false, err
	}

	if !exists {
		return limit.Count-count >= 0, nil
	}

	return currentVal-count >= 0, nil
}

func updateUsage(count int, limit *opaLimit) error {
	key, err := limit.key()
	if err != nil {
		return err
	}

	created, err := storage.create(key, limit.Count-1, time.Duration(limit.Seconds)*time.Second)
	if err != nil {
		return err
	}

	if created { // Already initialized with updated counter
		return nil
	}

	return storage.decrement(key, count)
}

func storeLimitUpdateInCtx(limit *opaLimit, increaseBy int, bctx topdown.BuiltinContext) error {
	val, exists := bctx.Cache.Get(limitsCtxKey)

	if !exists || val == nil {
		val = make(map[string]int)
		bctx.Cache.Put(limitsCtxKey, val)
	}

	key, err := limit.key()
	if err != nil {
		return err
	}

	val.(map[string]int)[key] = increaseBy

	return nil
}

func removeLimitUpdates(bctx topdown.BuiltinContext) {
	bctx.Cache.Put(limitsCtxKey, nil)
}
