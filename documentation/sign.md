# Sign hash

## **Scope**

Method for signing hash.

## **Authorization**

If "API_KEY" variable is set in environment, `API-Key` header shall be used in header

```sh
header 'API-Key: Strong_example'
```

## **Request**

The Service provider's application sends the following request using TLS for RSA keys:

```sh
POST /digest/sign
```

using ECDSA key:

```sh
POST /digest/sign-ecc?SignatureMethod=P1364
```

## sign-ecc keys

Since with ECDSA ypu can create diferent signatures, DER and P1363 signatures are implemented and fo choose signatureMethot, a key `SignatureMethod` with values `DER`  or `P1363` are implemented.

|**Key**|**Type**|**Description**|
| --- | --- | --- |
| `SignatureMethod` | *string* | Use `DER`  or `P1363` signing methods. if no key, `DER`  is default. |

### **Body** ECC

JSON

```json
{
    "hash": "string"
}
```

### **Body** RSA 

For RSA keys, for backwards compatability, there is a two ways how to make a request, single object (signature) or in array (batch)

#### RSA single object (signature) body

```json
    {
        "sessionId":"string",
        "hash": "string"
    }
```

#### RSA array (batch) body

```json
[
    {
        "sessionId":"string",
        "hash": "string"
    },
    {
        "sessionId":"string",
        "hash": "string"
    }
]
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `sessionId` | *string* | sessionId of the hash to be signed. For Single object (signature) optional. |
| `hash` | *string* | hash to be signed in base64 format |

### Example for single object (signature) request body (sessionId is optional)

```json
{
  "hash": "27aAZIjttlrjGyLMlcMcQh+nsltyVNLpxdog="
}
```

## Response

JSON object:

```json
{
    "sessionId": "string", // for rsa batch
    "signatureMethod": "string",
    "hash": "string",
    "signatureValue": "string"
}
```

### Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `sessionId` | *string* | sessionId of the hash to be signed if provided in request |
| `signatureMethod`  | *string* | Signature method used to sign|
| `hash`  | *string* | hash requested to be signed in base64 format|
| `signatureValue` | *string* | Signature value in base64 format |

### Note

If there was a batch request, response will be in array, if for single request, response will contain single object (without array)

### Example

```json
{
    "sessionId": "123kjn-131231",
    "signatureMethod": "P1363",
    "hash": "27aAZIjttlrjGyLMlcMcQh+nsltyVNLpxdog=",
    "signatureValue": "iyQGs/5hdq+....V/YsjOVA=="
}
```
