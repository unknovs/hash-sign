# Simple API for working with signatures

POST `/digest/sign` For signing hash with RSA keys using SHA256 algorithm

POST `/digest/sign-ecc` For signing hash with ECC keys

POST `/digest/verify` For verification of signed hash using public certificate

POST `/digest/calculateSummary` For digests summary calculation for one signature for use in Entrust TrustedX eIDAS Platform

POST `/digest/verificationCode` Calculates verification code by principle `integer(SHA256(hash)[-2:-1]) mod 10000`

GET `/certificates` For receiving a signing and authentication certificates stored in environment variables

POST `/asice/addFile` For adding a file to a asic-e container

POST `/encrypt/publicKey` For data encryption (RSA PKCS1Padding) using a PKCS1 RSA public key in PEM format.

POST `/digest/verificationCode` 4 digit verification code generation from hash to be signed.

POST `/jwt/generate` JWT token generation and signing with specific data

## Image

Latest image available on [docker hub](https://hub.docker.com/r/unknovs/hash-sign)

## Signing

* Application decodes received base64 hash to binary format
* Signs with RSA and ECC
* Encodes signed value to base64
* Returns base64 signed value.

## Verification

* input shall contain
  * `digestValue` - digest before signature
  * `signatureValue` - signatureValue (signed digest)
  * `certificate` - Public certificate in base64 format


## Environment

```yaml
    environment:
      PEM_FILE: "/run/secrets/key.pem"
      EC_PEM_FILE: "/run/secrets/ecc_key.pem"
      API_KEY: "Put_your_api_key_here"
      RSA_AUTH_CERT: "base64 encoded RSA signing certificate"
      RSA_SIGN_CERT: "base64 encoded RSA authentication certificate"
      ECDSA_AUTH_CERT: "base64 encoded ECDSA signing certificate"
      ECDSA_SIGN_CERT: "base64 encoded ECDSA authentication certificate"
      JWT_SIGNING_KEY: "jwt signing private key"

    secrets:
      - source: "rsa_private_key"
        target: "key.pem"
      - source: "ecc_private_key"
        target: "ecc_key.pem"
      - source: "jwt_signing_key"
        target: "jwt_signing_key"
    volumes:
      - temp:/tmp
volumes:
  temp:
secrets:
  rsa_private_key:
    external: true
  ecc_private_key:
    external: true 
  jwt_signing_key:
    external: true
```

`PEM_FILE` unencrypted RSA signing key in PEM format. Description below.

`EC_PEM_FILE` unencrypted ECDSA signing key in PEM format. Description below.

`API_KEY` Api key. Optional. If set, `API-Key` header shall be used in header.

`RSA_AUTH_CERT` base64 encoded RSA authentication certificate. Value between the `-----BEGIN CERTIFICATE-----` and `-----END CERTIFICATE-----` shall be provided. 

`RSA_SIGN_CERT` base64 encoded RSA signing certificate. Value between the `-----BEGIN CERTIFICATE-----` and `-----END CERTIFICATE-----` shall be provided.

`ECDSA_AUTH_CERT` base64 encoded ECDSA authentication certificate. Value between the `-----BEGIN CERTIFICATE-----` and `-----END CERTIFICATE-----` shall be provided.

`ECDSA_SIGN_CERT` base64 encoded ECDSA signing certificate. Value between the `-----BEGIN CERTIFICATE-----` and `-----END CERTIFICATE-----` shall be provided.

`JWT_SIGNING_KEY` PKCS8 PRIVATE KEY `FILE` in PEM format. **Including** `-----BEGIN PRIVATE KEY-----` and `-----END PRIVATE KEY-----`

### Secret creation from server terminal (SSH with root privileges)

Example for creating Docker swarm secrets from file.

Log into server with ssh and administrator privileges. Copy key file to server. For example, for RSA, Private key must be in PKCS#1 unencrypted format - starts with `-----BEGIN RSA PRIVATE KEY-----` and end with `-----END RSA PRIVATE KEY-----`.

```sh
docker secret create private_key /path/to/file/key.pem
```

### Secret creation from Portainer

#### For RSA private key

When creating a secret, copy content of pem file - starts with `-----BEGIN RSA PRIVATE KEY-----` and end with `-----END RSA PRIVATE KEY-----` to a secret.

#### For ECDSA private key

When creating a secret, copy content of pem file - starts with `-----BEGIN PRIVATE KEY-----` and end with `-----END PRIVATE KEY-----` to a secret.

## Methods

`/digest/sign` and `/digest/sign-ecc` method description [here](./documentation/sign.md)

`/digest/verify` method description [here](./documentation/verify.md)

`/digest/calculateSummary` method description [here](./documentation/calculateSummary.md)

`/digest/verificationCode` method description [here](./documentation/verificationCode.md)

`/certificates` method description [here](./documentation/certificates.md)

`/asice/addFile` method description [here](./documentation/addFile.md)

`/encrypt/publicKey` method description [here](./documentation/encrypt_with_public_key.md)

`/digest/verificationCode` method description [here](./documentation/verificationCode.md)

`/jwt/generate` method description [here](./documentation/generateJwt.md)

## Useful commands

You can find some useful commands for preparing key [here](./documentation/helper.md)
