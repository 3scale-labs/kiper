package threescale

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/3scale/3scale-go-client/threescale"

	threescaleAPI "github.com/3scale/3scale-go-client/threescale/api"
	apisonator "github.com/3scale/3scale-go-client/threescale/http"
	porta "github.com/3scale/3scale-porta-go-client/client"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/types"
)

const (
	adminPortalEnv         = "THREESCALE_ADMIN_PORTAL"
	accessTokenEnv         = "THREESCALE_ACCESS_TOKEN"
	serviceIdEnv           = "THREESCALE_SERVICE_ID"
	proxyConfigEnvironment = "production"
	adminPortalScheme      = "https"
	adminPortalPort        = 443
)

var ThreescaleAuthrepBuiltin = &ast.Builtin{
	Name: "authorize_with_3scale",
	Decl: types.NewFunction(types.Args(types.A), types.B),
}

func AuthrepWithThreescaleImpl(httpRequest ast.Value) (ast.Value, error) {
	// TODO: avoid instantiating the clients on every request
	apisonatorClient, err := apisonator.NewDefaultClient()
	if err != nil {
		return nil, err
	}

	adminPortalHost := os.Getenv(adminPortalEnv)
	if adminPortalHost == "" {
		return nil, fmt.Errorf("admin portal not set")
	}
	adminPortal, err := porta.NewAdminPortal(
		adminPortalScheme, adminPortalHost, adminPortalPort,
	)
	if err != nil {
		return nil, err
	}

	accessToken := os.Getenv(accessTokenEnv)
	if accessToken == "" {
		return nil, fmt.Errorf("access token not set")
	}

	service := serviceFromEnv()
	if service == "" {
		return ast.Boolean(false), fmt.Errorf("service ID not found")
	}

	portaClient := porta.NewThreeScale(adminPortal, accessToken, nil)
	proxyConfig, err := portaClient.GetLatestProxyConfig(string(service), proxyConfigEnvironment)
	if err != nil {
		return nil, err
	}

	request := Input{}
	_ = json.Unmarshal([]byte(httpRequest.String()), &request)

	clientAuth := clientAuthFromProxyConfig(&proxyConfig.ProxyConfig)
	if clientAuth == nil {
		return ast.Boolean(false), fmt.Errorf("service credentials not found")
	}

	appCreds := appCredsFromRequest(&request)
	if appCreds == nil {
		return ast.Boolean(false), fmt.Errorf("app credentials not found")
	}

	usage, err := usageFromMatchedRules(
		request.Attributes.Request.HTTP.Path,
		proxyConfig.ProxyConfig.Content.Proxy.ProxyRules,
	)
	if err != nil {
		return nil, err
	}

	// If there are no matches, it means that request to the request path are
	// not authorized
	if len(usage) == 0 {
		return ast.Boolean(false), nil
	}

	threescaleRequest := threescale.Request{
		Auth:    *clientAuth,
		Service: service,
		Transactions: []threescaleAPI.Transaction{
			{
				Metrics: usage,
				Params:  *appCreds,
			},
		},
	}

	resp, err := apisonatorClient.AuthRep(threescaleRequest)

	if err != nil {
		return nil, err
	}

	return ast.Boolean(resp.Success()), nil
}

func clientAuthFromProxyConfig(proxyConfig *porta.ProxyConfig) *threescaleAPI.ClientAuth {
	authType := proxyConfig.Content.BackendAuthenticationType
	authVal := proxyConfig.Content.BackendAuthenticationValue

	return &threescaleAPI.ClientAuth{
		Type:  threescaleAPI.AuthType(authType),
		Value: authVal,
	}
}

func serviceFromEnv() threescaleAPI.Service {
	serviceId := os.Getenv(serviceIdEnv)
	if serviceId == "" {
		return ""
	}

	return threescaleAPI.Service(serviceId)
}

func usageFromMatchedRules(path string, rules []porta.ProxyRule) (threescaleAPI.Metrics, error) {
	res := threescaleAPI.Metrics(make(map[string]int))

	for _, rule := range rules {
		// TODO: There are probably some pattern for which this does not work.
		regex, err := regexp.Compile(rule.Pattern)
		if err != nil {
			return nil, err
		}

		if regex.Match([]byte(path)) {
			res[rule.MetricSystemName] += int(rule.Delta)
		}
	}

	return res, nil
}

func appCredsFromRequest(request *Input) *threescaleAPI.Params {
	userKey := userKeyFromRequest(request)

	if userKey != "" {
		return &threescaleAPI.Params{UserKey: userKey}
	}

	appId, appKey := appIdAndKeyFromRequest(request)

	if appId != "" && appKey != "" {
		return &threescaleAPI.Params{
			AppID:  appId,
			AppKey: appKey,
		}
	}

	return nil
}

// Note: the location of the user key is configurable in 3scale. It can be in
// any header or query argument. For now we'll assume that if specified, it is
// in the "user_key" query arg.
func userKeyFromRequest(request *Input) string {
	return request.queryArgs()["user_key"]
}

// Note: same assumption as in userKeyFromRequest() for now.
func appIdAndKeyFromRequest(request *Input) (appID string, appKey string) {
	args := request.queryArgs()
	return args["app_id"], args["app_key"]
}
