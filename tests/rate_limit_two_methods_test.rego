package threescale.tests.rate_limit_two_methods

import input.attributes.request.http as http_request

test_get_has_a_limit_of_two_rpmin {
    test_request := { "method": "GET" }
    allow with http_request as test_request
    allow with http_request as test_request
    not allow with http_request as test_request
}

test_post_has_a_limit_of_one_rpmin {
    test_request := { "method": "POST" }
    allow with http_request as test_request
    not allow with http_request as test_request
}

test_rest_of_methods_are_unlimited {
    # To test "unlimited", just make more reqs than the sum of the two other
    # limits to be sure.
    test_request := { "method": "HEAD" }
    allow with http_request as test_request
    allow with http_request as test_request
    allow with http_request as test_request
}
