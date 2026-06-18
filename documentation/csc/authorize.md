# CSC oauth2/authorize

## **Scope**

OAuth 2.0 Authorization Code endpoint (CSC API v2 `oauth2/authorize`). In this mock it
**auto-approves** (no user interaction) and immediately redirects back with an
authorization code.

## **Authorization**

None (browser-facing OAuth endpoint).

## **Request**

A browser GET redirect:

```sh
GET /csc/v2/oauth2/authorize?response_type=code&client_id=..&redirect_uri=..&state=..&code_challenge=..&code_challenge_method=S256&scope=..
```

### Query

|**Query key**|**Type**|**Description**|
| --- | --- | --- |
| `redirect_uri` | *string* | Where the response redirects back with the code. **Required.** |
| `state` | *string* | Opaque value echoed back on the redirect. |
| `response_type` | *string* | `code`. Accepted, not enforced. |
| `client_id` | *string* | Client id. Accepted, not enforced. |
| `code_challenge`, `code_challenge_method` | *string* | PKCE challenge. Accepted, not enforced. |
| `scope` | *string* | Requested scope (`service` / `credential`). Accepted, not enforced. |

## **Response**

`302 Found` — a redirect to:

```sh
{redirect_uri}?code=<authorization_code>&state=<state>
```

`400` with the CSC error body when `redirect_uri` is missing or invalid:

```json
{
    "error": "invalid_request",
    "error_description": "missing redirect_uri"
}
```
