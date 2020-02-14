package envoy.authz

import input.attributes.request.http as http_request

default allow = false
default rate_limited = false

limited_user_ids := [ "a", "b" ]

allow {
    not rate_limited
    update_limits_usage()
}

rate_limited {
    http_request.path == "/abc"
    rate_limit({"by": {"path": http_request.path}, "count": 5, "seconds": 60})
}

rate_limited {
    http_request.path == "/abc"
    http_request.headers["user_id"] == limited_user_ids[_]
    rate_limit({"by": {"path": http_request.path, "user_id": http_request.headers["user_id"]}, "count": 3, "seconds": 60})
}
