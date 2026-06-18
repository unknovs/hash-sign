# CSC API mock

A minimal but **CSC API v2-faithful** test server (CSC API V2.2.0.0). It is **off by
default** and exposed only when `CSC_MOCK_ENABLED=true`. It implements the request and
response shapes of the CSC methods a remote-signing client uses, so any CSC client can
drive it: `signatures/signHash` is backed by the ECC signing key and `credentials/info`
by the `ECDSA_SIGN_CERT` environment value.

It is a test mock — it performs the raw `signHash` and serves the credential certificate,
but does no AdES packaging or timestamping (that is the calling application's job). The
`/csc/v2/*` routes are **not** behind `API-Key`; CSC uses its own OAuth bearer model and
the mock accepts any token.

## Methods

| Method | Path | Description |
| --- | --- | --- |
| `POST` | `/csc/v2/info` | [Service discovery](./info.md) (mandatory CSC method) |
| `GET`  | `/csc/v2/oauth2/authorize` | [Authorization](./authorize.md) (auto-approves) |
| `POST` | `/csc/v2/oauth2/token` | [Token](./token.md) |
| `POST` | `/csc/v2/credentials/list` | [Credential list](./credentialsList.md) |
| `POST` | `/csc/v2/credentials/info` | [Credential info](./credentialsInfo.md) |
| `POST` | `/csc/v2/signatures/signHash` | [Sign hash](./signHash.md) |

## Environment

| Variable | Description |
| --- | --- |
| `CSC_MOCK_ENABLED` | Set `true` to expose `/csc/v2/*` (default off). |
| `CSC_MOCK_SIG_ENCODING` | ECDSA signature encoding returned by `signHash`: `DER` (default) or `P1363`. |
| `EC_PEM_FILE` | ECC signing key (also used by `/digest/sign-ecc`). |
| `ECDSA_SIGN_CERT` | ECDSA signing certificate served by `credentials/info`. |

## Error responses

Errors use the CSC error body alongside a non-2xx HTTP status:

```json
{
    "error": "invalid_request",
    "error_description": "no hashes to sign"
}
```
