# Calculate digest summary

## **Scope**

Method for digest summary calculation for one signature for use in Entrust TrustedX eIDAS Platform to get digest signed

## **Authorization**

If "API_KEY" variable is set in environment, `API-Key` header shall be used in header

```
header 'API-Key: Strong_example'
```

## **Request**

The Service provider's application sends the following request using TLS:

```
GET /digest/calculateSummary/{digest-in-base64}
```

### Querry

Possible keys (Optional), if no keys defined, all registrated certifikates will be returned

|**Querry key**|**Possible values**|**Description**|
| --- | --- | --- |
| `hash` | *sha256* or *sha384* or *sha512* | Hashing algorithm in witch you will get diggest summary in response |

*NOTE: If no key is used, by default its set to `sha256`*

### **Example**

without using a key:

```sh
GET {host}/digest/calculateSummary/27aAZIjttlrBYu0hDHjGyLMlcMcQh+nsltyVNLpxdog=
```

using a `?hash=sha384` key:

```sh
GET {host}/digest/calculateSummary/27aAZIjttlrBYu0hDHjGyLMlcMcQh+nsltyVNLpxdog=?hash=sha384
```


## **Response**

JSON object:

```json
{
    "digestSummary": "string",
    "algorithm": "string"
}
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `digestSummary` | *string* | calculated digest summary value in base64 format for use in Entrust TrustedX eIDAS Platform |
| `algorithm` | *string* | algorithm used to calculate digest summary value in base64 format for use in Entrust TrustedX eIDAS Platform |


### **Example** 

```json
{
    "digestSummary": "6N5NR4zjElA...kwn9QcV2Q=",
    "algorithm": "sha512"
}
```
