# Encrypt data with public key

## **Scope**

Method for data encryption with `RSA PKCS1Padding` using a PKCS1 RSA public key in PEM format.

## **Authorization**

If "API_KEY" variable is set in environment, `API-Key` header shall be used in header

```
header 'API-Key: Strong_example'
```

## **Request**

The Service provider's application sends the following request using TLS:

```
POST /encrypt/publicKey
```

### **Body**

JSON
```json
{
    "dataToEncrypt": "string",
    "public_key": "string"
}
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `dataToEncrypt` | *string* | data to encrypt, shall be in format you need to encrypt. |
| `public_key` | *string* | PKCS1 public key. Data between header `-----BEGIN PUBLIC KEY-----` and footer `-----END PUBLIC KEY-----` shall be added|

### **Example**

```json
{
    "dataToEncrypt": "Encrypt this data",
    "public_key": "MIIBIjANBgk...iselNQIDAQAB"
}
```

## **Response**
JSON object:

```json
{
    "encryptedData": "string"
}
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `encryptedData`  | *string* | Encrypted data in base64 format|

### **Example** 

```json
{
    "encryptedData": "jAHzmXGQS...mgB3wWzA=="
}
```
