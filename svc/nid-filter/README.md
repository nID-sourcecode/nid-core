# Nid-filter
The nid-filter service contains four different request/response filters that are meant to be used as external processing filters with Envoy.

# Auditlog

The audit log filter logs requests with JWT claims and response code. This filter should be placed on the ingress of a service.

# Authswap

The authswap filter swaps the passed JWT token for a JWT token that is valid for the requested datasource. This filter should be placed on the egress of a service.

# Autopseudo

The autopseudo filter is able to replace the redacted bsn or subject inside the JWT token with the actual value. This filter should be placed on the ingress of a service.

# Scopeverification

The scopeverification verifies the gql/rest request with the scopes inside the JWT token. This filter should be placed on the ingress of a service. 

