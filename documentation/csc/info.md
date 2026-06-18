# CSC info

## **Scope**

Service discovery — returns the supported specification version, methods, signature
algorithms and authorization types (CSC API v2 `info`, the one mandatory CSC method).

## **Authorization**

None. `info` is a public discovery endpoint.

## **Request**

The client sends the following request using TLS:

```sh
POST /csc/v2/info
```

### **Body**

JSON (optional; `{}` is valid)

```json
{
    "lang": "en-US"
}
```

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `lang` | *string* | Optional language tag for human-readable fields. |

## **Response**

JSON object:

```json
{
    "specs": "2.2.0.0",
    "name": "hash-sign CSC mock",
    "authType": ["oauth2code"],
    "methods": ["credentials/list", "credentials/info", "signatures/signHash"],
    "signAlgorithms": {
        "algos": ["1.2.840.10045.4.3.2", "1.2.840.10045.4.3.3", "1.2.840.10045.4.3.4"]
    },
    "asynchronousOperationMode": false
}
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `specs` | *string* | CSC API specification version (`2.2.0.0`). |
| `name` | *string* | Service name. |
| `authType` | *array* | Supported authorization types (`oauth2code`). |
| `methods` | *array* | Supported CSC methods. |
| `signAlgorithms.algos` | *array* | Supported signature-algorithm OIDs (ECDSA with SHA-256 / SHA-384 / SHA-512). |
| `asynchronousOperationMode` | *boolean* | Whether asynchronous signing (`signPolling`) is supported. Always `false`. |
