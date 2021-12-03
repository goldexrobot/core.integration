# Goldex Robot Integration

## Signing/verifying a request

In order to call Goldex API JWT is used to sign a request.

Goldex callbacks are also signed with per-project key.
It is not mandatory but is preferred to validate those callbacks.
Developer is fully responsible for security.

JWT token should be transferred in `Authorization` header with `Bearer` prefix:

```text
Authorization: Bearer <jwt.goes.here>
```

### JWT claims

Here are default fields of JWT used during signing a request to Goldex API:

| Field | Purpose | Format |
| --- | --- | --- |
| aud | Recipient of the request | string(3-32): alphanumeric, `-`, `_` |
| iss | Issuer of JWT or API login | string(3-32): alphanumeric, `-`, `_` |
| jti | Unique request ID | string(6-36): alphanumeric, `-` (UUID compatible) |
| sub | The request entity or domain | string(32): alphanumeric |
| exp, nbf, iat | Are expiration, not before and issuance time | According to [RFC 7519](https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.5) |

Additional JWT fields:

| Field | Purpose | Format |
| --- | --- | --- |
| bha | Request body hash algorithm | string(16): `SHA-256`, `SHA-384`, `SHA-512`, `SHA3-224`, `SHA3-256`, `SHA3-384`, `SHA3-512` |
| bhs | Request body hash | string(32-128): hexadecimal without leading `0x` |
| mtd | Request method | string(8): GET, POST etc. |
| url | Request URL | string(256): valid URL |

Body hash alg and hash fields have to be empty for bodiless request such as GET.

Goldex **callbacks** carries JWT with next content:

| Field | Content |
| --- | --- |
| aud | ["project-N"] where N is project ID |
| iss | "goldex" |
| sub | "notification" or "sync-request" depending on context |

Allowed JWT signing algorithms: `HS256` (HMAC SHA-256), `HS384` (HMAC SHA-384), `HS512` (HMAC SHA-512).
