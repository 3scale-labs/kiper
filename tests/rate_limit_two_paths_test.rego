package threescale.tests.rate_limit_two_paths

import input.attributes.request.http as http_request

test_path_abc_has_limit_of_2_rpmin {
    test_request := { "path": "/abc"}
    allow with http_request as test_request
    allow with http_request as test_request
    not allow with http_request as test_request
}

test_path_def_has_limit_of_1_rpmin {
    test_request := { "path": "/def" }
    allow with http_request as test_request
    not allow with http_request as test_request
}

test_rest_of_paths_unlimited {
    # To test "unlimited", just make more reqs than the sum of the two other
    # limits to be sure.
    test_request := { "path": "/" }
    allow with http_request as test_request
    allow with http_request as test_request
    allow with http_request as test_request
    allow with http_request as test_request
}
