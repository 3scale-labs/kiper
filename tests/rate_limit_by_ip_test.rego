package threescale.tests.rate_limit_by_ip

import input.attributes as req

test_restricted_ips_have_a_limit_of_1_rpm {
    allow with req as { "source": { "address": "3.3.3.3" } }
    not allow with req as { "source": { "address": "3.3.3.3" } }

    allow with req as { "source": { "address": "4.4.4.4" } }
    not allow with req as { "source": { "address": "4.4.4.4" } }
}

test_less_restricted_ips_have_a_limit_of_3_rpm {
    allow with req as { "source": { "address": "1.1.1.1" } }
    allow with req as { "source": { "address": "1.1.1.1" } }
    allow with req as { "source": { "address": "1.1.1.1" } }
    not allow with req as { "source": { "address": "1.1.1.1" } }

    allow with req as { "source": { "address": "2.2.2.2" } }
    allow with req as { "source": { "address": "2.2.2.2" } }
    allow with req as { "source": { "address": "2.2.2.2" } }
    not allow with req as { "source": { "address": "2.2.2.2" } }
}

test_rest_of_ips_have_a_limit_of_2_rpm {
    allow with req as { "source": { "address": "5.5.5.5" } }
    allow with req as { "source": { "address": "5.5.5.5" } }
    not allow with req as { "source": { "address": "5.5.5.5" } }
}
