# CSC credentials/info

## **Scope**

Returns a credential's certificate and key/certificate details (CSC API v2
`credentials/info`). The response is built from the `ECDSA_SIGN_CERT` environment value:
the certificate is returned in `cert.certificates` (base64-DER) and the descriptive fields
are read from it.

## **Authorization**

Bearer access token. The mock accepts any bearer.

```sh
header 'Authorization: Bearer <access_token>'
```

## **Request**

The client sends the following request using TLS:

```sh
POST /csc/v2/credentials/info
```

### **Body**

JSON

```json
{
    "credentialID": "ecc-sign-1",
    "certificates": "chain",
    "certInfo": true
}
```

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `credentialID` | *string* | The credential to read. |
| `certificates` | *string* | `chain` / `single` / `none`. Accepted, not enforced. |
| `certInfo` | *boolean* | Request parsed certificate details. Accepted, not enforced. |

## **Response**

JSON object:

```json
{
    "key": { "status": "enabled", "algo": ["1.2.840.10045.2.1"], "len": 384 },
    "cert": {
        "status": "valid",
        "certificates": ["<base64-DER certificate>"],
        "subjectDN": "CN=...,SERIALNUMBER=...,C=LV",
        "issuerDN": "CN=...,O=...,C=LV",
        "serialNumber": "0A1B2C...",
        "validFrom": "20260101000000Z",
        "validTo": "20270101000000Z",
        "qcStatements": ["0.4.0.1862.1.1", "0.4.0.1862.1.6.1"]
    },
    "signatureQualifier": "eu_eidas_qes",
    "auth": { "mode": "oauth2code" },
    "SCAL": "2",
    "multisign": 50
}
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `key.status` | *string* | Key status (`enabled`). |
| `key.algo` | *array* | Public-key algorithm OIDs (`1.2.840.10045.2.1` = id-ecPublicKey). |
| `key.len` | *integer* | Key length in bits (curve size; omitted if the certificate cannot be parsed). |
| `cert.status` | *string* | Certificate status (`valid`). |
| `cert.certificates` | *array* | The certificate chain, leaf first, base64-DER (no PEM headers). |
| `cert.subjectDN` / `cert.issuerDN` | *string* | Subject / issuer distinguished name, read from the certificate. |
| `cert.serialNumber` | *string* | Certificate serial number (hex). |
| `cert.validFrom` / `cert.validTo` | *string* | Validity window (`YYYYMMDDHHMMSSZ`). |
| `cert.qcStatements` | *array* | QC statement OIDs read from the certificate's qcStatements extension (omitted if absent). |
| `signatureQualifier` | *string* | Signature qualifier (`eu_eidas_qes`). |
| `auth.mode` | *string* | Authorization mode (`oauth2code`). |
| `SCAL` | *string* | Sole-control assurance level. |
| `multisign` | *integer* | Maximum signatures per authorization. |
