package threescale.tests.rate_limit_two_methods

# Different limits for two HTTP methods. The rest are unlimited:

import input.attributes.request.http as http_request

default allow = false
default rate_limited = false

allow {
    not rate_limited
    update_limits_usage()
}

rate_limited {
    http_request.method == "GET"
    rate_limit({"by": {"method": http_request.method}, "count": 2, "seconds": 60})
}

rate_limited {
    http_request.method == "POST"
    rate_limit({"by": {"path": http_request.method}, "count": 1, "seconds": 60})
}
