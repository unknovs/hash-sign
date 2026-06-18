# CSC oauth2/token

## **Scope**

OAuth 2.0 token endpoint (CSC API v2 `oauth2/token`). Exchanges the authorization code for
an access token used as the `Bearer` on the credential and signing calls.

## **Authorization**

Client authentication (HTTP Basic `client_id:client_secret`). The mock accepts any client
credentials.

## **Request**

The client sends the following request using TLS:

```sh
POST /csc/v2/oauth2/token
```

### **Body**

`application/x-www-form-urlencoded`

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `grant_type` | *string* | `authorization_code`. |
| `code` | *string* | The authorization code from `oauth2/authorize`. |
| `redirect_uri` | *string* | Must match the authorize request. |
| `code_verifier` | *string* | PKCE verifier. |
| `client_id` | *string* | Client id. |

## **Response**

JSON object:

```json
{
    "access_token": "string",
    "token_type": "Bearer",
    "expires_in": 3600
}
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `access_token` | *string* | Bearer token for `credentials/*` and `signatures/signHash`. |
| `token_type` | *string* | `Bearer`. |
| `expires_in` | *integer* | Token lifetime in seconds. |
