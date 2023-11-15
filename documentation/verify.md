# Verify signed data

## **Scope**

Method for verifying signed data.

## **Authorization**

If "API_KEY" variable is set in environment, `API-Key` header shall be used in header

```
header 'API-Key: Strong_example'
```

## **Request**

The Service provider's application sends the following request using TLS:

```
POST /digest/verify
```

### **Body**

JSON
```json
{
    "digestValue": "string",
    "signatureValue": "string",
    "certificate": "string"
}
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `digestValue` | *string* | digest before signature in base64 format. If u are using `/sign` then `hash` value in request or response |
| `signatureValue` | *string* | signatureValue (signed digest) in base64 format. If u are using `/sign` then `signatureValue` received in response |
| `certificate` | *string* |  Public certificate in base64 format. If u are using `/sign` then Public certificate of the private key loaded in `PEM_FILE` variable.|


### **Example**

```json
{
    "digestValue": "zH/19ZUeiZrDlFbnTunPt3pOpkYeF/KS8OjmJWDoaTg=",
    "signatureValue": "H/rUJkDf3eLykp+GIv...l8gXn6eSbxll69rlYc6Fg==",
    "certificate": "MIIG6jCCBNKgAwIBAgIQ...your_public_certificate_base64_here...Diyj+2aew=="
}
```

## **Response**

status code and Message

`200` - `Signature is valid!`  - signature is valid

`400`:
* `Invalid digest value` - digest provided can't be decoded
* `Failed to verify signature` - signature verification failed, `digestValue` is not the same as signed in `signatureValue`
* `Failed to parse certificate: x509: malformed certificate` - provided certificate cant be parsed
* `Invalid signature value` - provided signature value cant be decoded
* `invalid public key algorithm` - Provided certificate do not contain RSA key