# Simple API for working with signatures

POST `/digest/sign` For signing hash with RSA keys using SHA256 algorithm

POST `/digest/sign-ecc` For signing hash with ECC keys

POST `/digest/verify` For verification of signed hash using public certificate

GET `/digest/calculateSummary/` For digests summary calculation for one signature for use in Entrust TrustedX eIDAS Platform

GET `/certificates` For receiving a signing and authentication certificates stored in environment variables

POST `/asice/addFile` For adding a file to a asic-e container

POST `/encrypt/publicKey` For data encryption (RSA PKCS1Padding) using a PKCS1 RSA public key in PEM format.

## Image

Latest image available on [docker hub](https://hub.docker.com/r/unknovs/hash-sign)

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
      SIGN_CERT: "base64 encoded signing certificate"
      AUTH_CERT: "base64 encoded authentication certificate"
    secrets:
      - source: "private_key"
        target: "key.pem"
    volumes:
      - temp:/tmp
volumes:
  temp:
secrets:
  private_key_prod:
    external: true    
```

`PEM_FILE` unencrypted RSA signing key in PEM format. 

`API_KEY` Api key. Optional. If set, `API-Key` header shall be used in header.

`SIGN_CERT` base64 encoded signing certificate

`AUTH_CERT` base64 encoded authentication certificate

### Secret creation from server terminal (SSH with root privileges)

Example for creating Docker swarm secrets from file.

Log into server with ssh and administrator privileges. Copy key file to server. Private key must be in PKCS8 unencrypted format - starts with `-----BEGIN PRIVATE KEY-----` and end with `-----END PRIVATE KEY-----`.

```sh
docker secret create private_key /path/to/file/key.pem
```
### Secret creation from Portainer

When creating a secret, copy content of pem file - starts with `-----BEGIN PRIVATE KEY-----` and end with `-----END PRIVATE KEY-----` to a secret.


## Methods

`/digest/sign` and `/digest/sign-ecc` method description [here](./documentation/sign.md)

`/digest/verify` method description [here](./documentation/verify.md)

`/digest/calculateSummary/` method description [here](./documentation/calculateSummary.md)

`/certificates` method description [here](./documentation/certificates.md)

`/asice/addFile` method description [here](./documentation/addFile.md)

`/encrypt/publicKey` method description [here](./documentation/encrypt_with_public_key.md)

## Useful commands

You can find some useful commands for preparing key [here](./documentation/helper.md)