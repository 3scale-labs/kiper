package envoy.authz

import input.attributes.request.http as http_request

default allow = false

allow {
	authorize_with_3scale(input)
}
