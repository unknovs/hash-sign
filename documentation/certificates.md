# Get certificates

## **Scope**

Return certificates from environment values

## **Authorization**

If "API_KEY" variable is set in environment, `API-Key` header shall be used in header

```
header 'API-Key: Strong_example'
```

## **Request**

The Service provider's application sends the following request using TLS:

```
GET /certificates
```

### Querry

Possible keys (Optional), if no keys defined, all registrated certifikates will be returned

|**Querry key**|**Possible values**|**Description**|
| --- | --- | --- |
| `key` | *rsa* or *ecdsa* | If u need to specify witch key type you need in response (certs shall be registrated in environment)|
| `type` | *auth* or *sign* | If u need to specify witch type of certificate you need in response (certs shall be registrated in environment) |


### **Example**

```sh
GET {host}/certificates
```

or

```sh
GET {host}/certificates?key=rsa&type=auth
```

## **Response**

JSON object if no keys used in querry:

```json
{
    "rsa_authentication_certificate": "string",
    "rsa_signing_certificate": "string",
    "ecdsa_authentication_certificate": "string",
    "ecdsa_signing_certificate": "string"
}
```

JSON object if `?key=rsa&type=auth` used in querry:

```json
{
    "rsa_authentication_certificate": "string",
}
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `rsa_authentication_certificate` | *string* | Base64 encoded RSA authentication certificate from environment (if set)|
| `rsa_signing_certificate` | *string* | Base64 encoded RSA signing certificate from environment (if set) |
| `ecdsa_authentication_certificate` | *string* | Base64 encoded ECDSA authentication certificate from environment (if set) |
| `ecdsa_signing_certificate` | *string* | Base64 encoded ECDSA signing certificate from environment (if set) |

### **Example** 

if no keys used in querry:

```json
{
    "rsa_authentication_certificate": "MIIGRzCCBC...1ohzvdaO+LaKIqazQ=",
    "rsa_signing_certificate": "MIIG6jCCBNKgAwIB...PZyabTTbNo6tUAim8j+2aew==",
    "ecdsa_authentication_certificate": "MIIGRzCCBC...1ohzvdaO+LaKIqazQ=",
    "ecdsa_signing_certificate": "MIIG6jCCBNKgAwIB...PZyabTTbNo6tUAim8j+2aew=="
}
```

if `?key=rsa&type=auth` used in querry:

```json
{
    "rsa_authentication_certificate": "MIIGRzCCBC...1ohzvdaO+LaKIqazQ="
}
```