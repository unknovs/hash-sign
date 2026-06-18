# CSC signatures/signHash

## **Scope**

Signs one or more digests (CSC API v2 `signatures/signHash`). Backed by the ECC signing
key; one call signs the whole batch and returns the signatures in request order.

## **Authorization**

Bearer credential-scoped access token. The mock accepts any bearer.

```sh
header 'Authorization: Bearer <access_token>'
```

## **Request**

The client sends the following request using TLS:

```sh
POST /csc/v2/signatures/signHash
```

### **Body**

JSON

```json
{
    "credentialID": "ecc-sign-1",
    "hashes": ["<base64-digest>", "<base64-digest>"],
    "signAlgo": "1.2.840.10045.4.3.2",
    "operationMode": "S"
}
```

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `credentialID` | *string* | The credential to sign with. |
| `hashes` | *array* | Base64 digests to sign, order preserved. `hash` is accepted as an alias. |
| `signAlgo` | *string* | Signature-algorithm OID. |
| `operationMode` | *string* | `S` (synchronous). Optional; async is not supported. |
| `SAD` | *string* | Optional signature activation data; unused (the oauth2code credential token is the authorization). |

### Signature encoding

The ECDSA signature encoding is controlled by the `CSC_MOCK_SIG_ENCODING` environment
value: `DER` (default, ASN.1 `SEQUENCE{r,s}`) or `P1363` (raw `r‖s`).

### Example

```json
{
    "credentialID": "ecc-sign-1",
    "hashes": ["27aAZIjttlrjGyLMlcMcQh+nsltyVNLpxdog="],
    "signAlgo": "1.2.840.10045.4.3.2"
}
```

## **Response**

JSON object:

```json
{
    "signatures": ["<base64-signature>"]
}
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `signatures` | *array* | Base64 signature values, in the same order as the request digests. |

### Errors

`400` with the CSC error body, e.g.:

```json
{
    "error": "invalid_request",
    "error_description": "no hashes to sign"
}
```
