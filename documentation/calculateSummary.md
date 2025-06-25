# Calculate digest summary

## Scope

Method for digest summary calculation for one signature for use in Entrust TrustedX eIDAS Platform to get digest signed

## Authorization

If "API_KEY" variable is set in environment, `API-Key` header shall be used in header

```sh
header 'API-Key: Strong_example'
```

## Request

The Service provider's application sends the following request using TLS:

```sh
POST /digest/calculateSummary
```

### Querry

Possible keys (Optional), if no keys defined, all registrated certifikates will be returned

|**Querry key**|**Possible values**|**Description**|
| --- | --- | --- |
| `hash` | *sha256* or *sha384* or *sha512* | Hashing algorithm in witch you will get diggest summary in response |

*NOTE: If no key is used, by default its set to `sha256`*

### Body

JSON

```json
{
    "digest": "string"
}
```

### Example

without using a key:

```sh
GET {host}/digest/calculateSummary
```

using a `?hash=sha384` key:

```sh
GET {host}/digest/calculateSummary?hash=sha384
```

body

```json
{
    "digest":"nofaW3trG2Q2OlnJZKLd620pCMjVOKezDoBveEAojA5W7yKjMyJv6fBA634hT2ns"
}
```

## Response

JSON object:

```json
{
    "digestSummary": "string",
    "URLSafeDigestSummary": "string",
    "algorithm": "string"
}
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `digestSummary` | *string* | calculated digest summary value in base64 format |
| `URLSafeDigestSummary` | *string* | calculated digest summary value in URL safe base64 format for use in Entrust TrustedX eIDAS Platform |
| `algorithm` | *string* | algorithm used to calculate digest summary value in base64 format |

### Response example

```json
{
    "digestSummary": "JWAUVjnwS3+xV0HS.....abRUiP3eWvr4S2",
    "URLSafeDigestSummary": "JWAUVjnwS3-xV0HS.....abRUiP3eWvr4S2",
    "algorithmUsed": "sha384"
}
```
