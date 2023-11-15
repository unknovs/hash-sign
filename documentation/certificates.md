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

### **Example**

```sh
GET {host}/certificates
```

## **Response**

JSON object:

```json
{
    "authentication_certificate": "string",
    "signing_certificate": "string"
}
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `authentication_certificate` | *string* | Base64 encoded authentication certificate from environment (if set)|
| `signing_certificate` | *string* | Base64 encoded signing certificate from environment (if set) |

### **Example** 

```json
{
    "authentication_certificate": "MIIGRzCCBC...1ohzvdaO+LaKIqazQ=",
    "signing_certificate": "MIIG6jCCBNKgAwIB...PZyabTTbNo6tUAim8j+2aew=="
}
```
