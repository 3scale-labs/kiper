package threescale.tests.rate_limit_by_header

# Limit by the value of the "user_id" header. All the values have the same limit
# associated. Denies when the header is not set:

import input.attributes.request.http as http_request

default allow = false
default rate_limited = false

allow {
    has_user_id
    not rate_limited
    update_limits_usage()
}

has_user_id {
    http_request.headers["user_id"] != null
}

rate_limited {
    rate_limit({"by": {"user_id": http_request.headers["user_id"]}, "count": 2, "seconds": 60})
}
