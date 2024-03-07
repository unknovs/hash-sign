# Sign hash

## **Scope**

Method for signing hash.

## **Authorization**

If "API_KEY" variable is set in environment, `API-Key` header shall be used in header

```
header 'API-Key: Strong_example'
```

## **Request**

The Service provider's application sends the following request using TLS for RSA keys:

```
POST /digest/sign
```

uzsing ECDSA key:

```
POST /digest/sign-ecc
```

### **Body**

JSON
```json
{
    "hash": "string"
}
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `hash` | *string* | hash to be signed in base64 format |

### **Example**

```json
{
  "hash": "27aAZIjttlrjGyLMlcMcQh+nsltyVNLpxdog="
}
```

## **Response**
JSON object:

```json
{
    "hash": "string",
    "signatureValue": "string"
}
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `hash`  | *string* | hash requested to be signed in base64 format|
| `signatureValue` | *string* | Signature value in base64 format |

### **Example** 

```json
{
    "hash": "27aAZIjttlrjGyLMlcMcQh+nsltyVNLpxdog=",
    "signatureValue": "iyQGs/5hdq+....V/YsjOVA=="
}
```
