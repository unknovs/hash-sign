# Generate verification code

## **Scope**

Calculates verification code by principle `integer(SHA256(hash)[-2:-1]) mod 10000`

This calculation can be used to calculate SmartId verification codes

## **Authorization**

If "API_KEY" variable is set in environment, `API-Key` header shall be used in header

```bash
header 'API-Key: Strong_example'
```

## **Request**

The Service provider's application sends the following request using TLS for RSA keys:

```bash
POST /digest/sign
```

uzsing ECDSA key:

```bash
POST /digest/verificationCode
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
| `hash` | *string* | hash in base64 format |

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
    "verification_code": "number"
}
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `verification_code`  | *number* | 4 digit number |

### **Example** 

```json
{
    "verification_code": 1251
}
```
