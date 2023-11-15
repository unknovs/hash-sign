# Add file to signed asic-e container

Method will work only if containers folder `/tmp` will be mapped to host. Method need an access to temp folder in order to manage files.

## **Scope**

Method for adding a file to asice-e container. Usually needed if you sign a file hash and then add that file to a asic-e container containing signature of that file.

## **Authorization**

If "API_KEY" variable is set in environment, `API-Key` header shall be used in header

```
header 'API-Key: Strong_example'
```

## **Request**

The Service provider's application sends the following request using TLS:

```
POST /asice/addFile
```
### Key type in request

Request can contain a key `type`, there is a two `type`'s available:
`?type=binary` - default - response will contain asic-e container with added files as binary
`?type=base64` - response will contain JSON where asic-e container with added files will be base64 encoded

### **Body**

You will ask, why base64 encoded files, why not a binaries? Main reason - postman, you cant deal with binaries in Pre-request scripts.

JSON
```json
{
  "emptyAsice": "string",
  "signedFiles": [
    {
      "fileName": "string",
      "encodedFile": "string"
    }
  ]
}
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `emptyAsice` | *string* | base64 encoded asic-e container to hold files |
| `signedFiles` | *array* | Array of files with filenames to add in asic-e container |
| `signedFiles.fileName` | *string* | filename with extension |
| `signedFiles.encodedFile` | *string* | base64 encoded file that is planned to be placed in the asic-e container |

### **Example**

```json
{
  "emptyAsice": "UEsDBAoAAAgAANpEV...AANzEAAAAA",
  "signedFiles": [
    {
      "fileName": "example.txt",
      "encodedFile": "dGVzdDE="
    }
  ]
}
```

## **Response**

### If type is binary or without a type key

Body will contain binary file

### If type is base64

JSON object:

```json
{
    "packedAsice":"string"
}
```

Description of properties

|**Property**|**Type**|**Description**|
| --- | --- | --- |
| `packedAsice`  | *string* | Base64 encoded asic-e container with added files |

### **Example** 

```json
{
    "packedAsice":"UEsDBAoAAAgAANasdpEV...AANzEAasdasdAAA"
}
```
