# Rate limiting examples

Different limits for two paths. The rest are unlimited:

```
default allow = false
default rate_limited = false

allow {
    not rate_limited
    update_limits_usage()
}

rate_limited {
    http_request.path == "/abc"
    rate_limit({"by": {"path": http_request.path}, "count": 5, "seconds": 60})
}

rate_limited {
    http_request.path == "/def"
    rate_limit({"by": {"path": http_request.path}, "count": 2, "seconds": 60})
}
```


Global rate limit based on client source IP with a set of restricted and unrestricted IPs.

```
# Let's imagine we want to rate limit based on client ip addresses the traffic to some servers
#   * Global Limit by IP of 100 r/min
#   * A set of IPs can do up to 200 r/min instead of 100r/min
#   * A set of IPs can only do 10r/min as those are more restricted. 

# Alias req to input.attributes, so we can get the "Source" doing req.source.address
#	
#  Why? see the golang structs: 
#
#  type Input struct {
#		ParsedPath  []string    `json:"parsed_path"`
#		ParsedQuery ParsedQuery `json:"parsed_query"`
#		ParsedBody  ParsedBody  `json:"parsed_body"`
#		Attributes  Attributes  `json:"attributes"`
#	}
#
#	type Attributes struct {
#		Source      Destination `json:"source"`
#		Destination Destination `json:"destination"`
#		Request     Request     `json:"request"`
#   }

import input.attributes as req

# Set defaults.
default allow = false
default rate_limited = false

# Define to arrays of IPs, the restricted and the unrestricted.
non_restricted_ips = {"127.0.0.1","2.2.2.2"}
restricted_ips = {"172.16.0.1"}

# We allow the request to go through only if is "not" rate_limited.
# All the conditions inside allow are treated as "AND".
allow {
    not rate_limited
    update_limits_usage()
}

# Let's define all the rate_limited conditions. Those are evaluated and "ORed" together.


# If the source address is in the non_restricted_ip list, we check for a 200r/60s limit.
rate_limited {
    non_restricted_ips[req.source.address]
    rate_limit({"by": {"client_ip": req.source.address}, "count": 200, "seconds": 60})
}

# If the source address is in the restricted list, we check for a 20r/60s limit.
rate_limited {
    restricted_ips[req.source.address]
    rate_limit({"by": {"client_ip": req.source.address}, "count": 20, "seconds": 60})
}

# If the source address is not in either list, we check for a 100r/60s limit.
rate_limited {
    not non_restricted_ips[req.source.address]
    not restricted_ips[req.source.address]
    rate_limit({"by": {"client_ip": req.source.address}, "count": 100, "seconds": 60})
}
```



Different limits for two methods. The rest are unlimited:

```
default allow = false
default rate_limited = false

allow {
    not rate_limited
    update_limits_usage()
}

rate_limited {
    http_request.method == "GET"
    rate_limit({"by": {"method": http_request.method}, "count": 5, "seconds": 60})
}

rate_limited {
    http_request.method == "HEAD"
    rate_limit({"by": {"path": http_request.method}, "count": 2, "seconds": 60})
}
```


Limit by header. All the values have the same limit and denies if the header is not set:

```
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
    rate_limit({"by": {"user_id": http_request.headers["user_id"]}, "count": 5, "seconds": 60})
}
```

Use a generic limit of 5 rps for a path but limit some users to 3:

```
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
```
