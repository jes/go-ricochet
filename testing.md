# Ricochet Testing Specification

This documents outlines each scenario this library must correctly handle and
links it to the automated test that exercises that functionality.


# Version Negotiation

## Open

File: [iricochet_test.go](./ricochet_test.go)

This stub test exercises the Open() function. `TestRicochetOpen`, and in cases where the server returns a bad version, `TestRicochetOpenWithError`, and where there is no server at all `TestRicochetOpenWithNoServer` - in both cases Open should return an error.

## Inbound

File: [inbound_version_negotiation_test.go](./inbound_version_negotiation_test.go)

### Invalid Protocol

If the inbound listener receives:

* Less than 4 bytes (`TestBadProtcolLength`)
* The first 2 bytes are not equal ot 0x49 and 0x4D
* A number of supported Versions &lt; 1 (`TestNoSupportedVersions`, `TestInvalidVersionList`)

Then it must close the connection.

### No Compatible Version Found

If the inbound listener does not receive a compatible version in the list of
supported versions. Then is must close the connection. `TestNoCompatibleVersions`

### Successful Version Negotiation 

Assuming the inbound listener receives a valid protocol message, and that message
contains a known supported version. Then the connection should remain open. `TestNegotiateInboundVersions`

## Outbound

File: [outbound_version_negotiation_test.go](./outbound_version_negotiation_test.go)

### No Compatible Version Found

If the outbound connection receives a response that does not match one of the versions
they sent out in their supporting list. Then then must close the connection `TestInvalidResponse` , `TestInvalidServer`

### Successful Version Negotiation 

Assuming the outbound connection receives a valid protocol message, and that message
contains a known supported version. Then the connection should remain open. `TestOutboundVersionNegotiation`

