package envoy.authz

# Different limits for two paths. The rest are unlimited:

import input.attributes.request.http as http_request

default allow = false
default rate_limited = false

allow {
    not rate_limited
    update_limits_usage()
}

rate_limited {
    http_request.path == "/abc"
    rate_limit({"by": {"path": http_request.path}, "count": 2, "seconds": 60})
}

rate_limited {
    http_request.path == "/def"
    rate_limit({"by": {"path": http_request.path}, "count": 1, "seconds": 60})
}
