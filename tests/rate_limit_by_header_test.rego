package threescale.tests.rate_limit_by_header

import input.attributes.request.http as http_request

test_limit_is_per_user_id_header {
    allow with http_request as { "headers": { "user_id": "1" } }
    allow with http_request as { "headers": { "user_id": "1" } }
    not allow with http_request as { "headers": { "user_id": "1" } }

    allow with http_request as { "headers": { "user_id": "2" } }
    allow with http_request as { "headers": { "user_id": "2" } }
    not allow with http_request as { "headers": { "user_id": "2" } }
}

test_without_user_id_header {
    not allow with http_request as { "headers": {} }
}
