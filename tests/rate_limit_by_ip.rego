package threescale.tests.rate_limit_by_ip

# Limit based on client ip addresses the traffic to some servers:
#   * Global Limit by IP of 2 r/min.
#   * A set of IPs can do up to 3 r/min instead of 2 r/min.
#   * A set of IPs can only do 1 r/min as those are more restricted.

import input.attributes as req

default allow = false
default rate_limited = false

less_restricted_ips = {"1.1.1.1", "2.2.2.2"}
restricted_ips = {"3.3.3.3", "4.4.4.4"}

allow {
    not rate_limited
    update_limits_usage()
}

rate_limited {
    less_restricted_ips[req.source.address]
    rate_limit({"by": {"client_ip": req.source.address}, "count": 3, "seconds": 60})
}

rate_limited {
    restricted_ips[req.source.address]
    rate_limit({"by": {"client_ip": req.source.address}, "count": 1, "seconds": 60})
}

rate_limited {
    not less_restricted_ips[req.source.address]
    not restricted_ips[req.source.address]
    rate_limit({"by": {"client_ip": req.source.address}, "count": 2, "seconds": 60})
}
