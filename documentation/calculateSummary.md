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

### **Example**

```sh
GET {host}/digest/calculateSummary/27aAZIjttlrBYu0hDHjGyLMlcMcQh+nsltyVNLpxdog=
```

## **Response**

JSON object:

```json
{
    "digestSummary": "string"
}
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `digestSummary` | *string* | calculated digest summary value in base64 format for use in Entrust TrustedX eIDAS Platform |

### **Example** 

```json
{
    "digestSummary": "6N5NR4zjElA...kwn9QcV2Q="
}
```
