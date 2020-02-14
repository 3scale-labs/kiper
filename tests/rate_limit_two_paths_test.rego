package envoy.authz

import input.attributes.request.http as http_request

test_path_abc_has_limit_of_2_rpmin {
    allow with http_request as {"path": "/abc" }
    allow with http_request as {"path": "/abc" }
    not allow with http_request as {"path": "/abc" }
}

test_path_def_has_limit_of_1_rpmin {
    allow with http_request as {"path": "/def" }
    not allow with http_request as {"path": "/def" }
}

test_rest_of_paths_unlimited {
    # To test "unlimited", just make more reqs than the sum of the two other
    # limits to be sure.
    allow with http_request as {"path": "/" }
    allow with http_request as {"path": "/" }
    allow with http_request as {"path": "/" }
    allow with http_request as {"path": "/" }
}
