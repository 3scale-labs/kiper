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
