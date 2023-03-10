# Simple API for working with RSA signatures

POST `/sign` For signing hash with RSA keys using SHA256 algorithm

POST `/verify` For verification of signed hash using public certificate

## Signing 

* Application decodes received base64 hash to binary format
* Signs with RSA SHA256 
* Encodes signed value to base64
* Returns base64 signed value.

## Verification 

* input shall contain 
-  `digestValue` - digest before signature
-  `signatureValue` - signatureValue (signed digest)
-  `certificate` - Public certificate in base64 format


## Environment

```
    environment:
      PEM_FILE: "/run/secrets/key.pem"
      API_KEY: "Put_your_api_key_here"
    secrets:
      - source: "private_key"
        target: "key.pem"
```

`PEM_FILE` unencrypted RSA signing key in PEM format. 

`API_KEY` Api key. Optional. If set, `API-Key` header shall be used in header.

### Secret creation from server terminal (SSH with root privileges)

Example for creating Docker swarm secrets from file.

Log into server with ssh and administrator privileges. Copy key file to server. Private key must be in PKCS8 unencrypted format - starts with `-----BEGIN PRIVATE KEY-----` and end with `-----END PRIVATE KEY-----`.

```sh
docker secret create private_key /path/to/file/key.pem
```

## Methods

`/sign` method description [here](./docs/sign.md)

`/verify` method description [here](./docs/verify.md)

## Useful commands

You can find some useful commands for preparing key [here](./docs/helper.md)