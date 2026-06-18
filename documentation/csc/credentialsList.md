# CSC credentials/list

## **Scope**

Lists the credential identifiers available to the token holder (CSC API v2
`credentials/list`).

## **Authorization**

Bearer access token (from `oauth2/token`). The mock accepts any bearer.

```sh
header 'Authorization: Bearer <access_token>'
```

## **Request**

The client sends the following request using TLS:

```sh
POST /csc/v2/credentials/list
```

### **Body**

JSON (`{}` is valid)

```json
{
    "onlyValid": true
}
```

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `onlyValid` | *boolean* | Optional; return only currently-valid credentials. Accepted, not enforced. |

## **Response**

JSON object:

```json
{
    "credentialIDs": ["ecc-sign-1"]
}
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `credentialIDs` | *array* | The available credential identifiers. The mock exposes a single ECC signing credential. |
