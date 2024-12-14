# Calculate digest summary

## **Scope**

Generate and sign JWT token

## **Authorization**

If "API_KEY" variable is set in environment, `API-Key` header shall be used in header

```bash
header 'API-Key: Strong_example'
```

## **Request**

The Service provider's application sends the following request using TLS:

```bash
POST /jwt/generate
```

### Body

JSON

```json
{
    "iss":"string",
    "aud":"string",
    "sub":"string"
}
```

### Description of request properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `iss` | *string* | Issuer |
| `aud` | *string* | Audience |
| `sub` | *string* | Subject |

## **Response**

JSON object:

```json
{
    "token": "string"
}
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `token` | *string* | Signed JWT token |
