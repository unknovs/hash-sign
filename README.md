# Simple API for signing hashes

For signing hash with RSA keys using SHA256 algorithm

## Process

* Application decodes received base64 hash to binary format
* Signs with RSA SHA256 
* Encodes signed value to base64
* Returns base64 signed value.

## Environment

```
    environment:
      PFX_FILE: "/run/secrets/key.pem"
    secrets:
      - source: "private_key"
        target: "key.pem"
```

`PFX_FILE` unencrypted RSA signing key in PEM format. 

### Secret creation from server terminal (SSH with root privileges)

Example for creating Docker swarm secrets from file.

Log into server with ssh and administrator privileges. Copy key file to server. Private key must be in PKCS8 unencrypted format - starts with `-----BEGIN PRIVATE KEY-----` and end with `-----END PRIVATE KEY-----`.

```sh
docker secret create private_key /path/to/file/key.pem
```

## Methods

You can find method description [here](./docs/sign.md)

## Useful commands

You can find some usefull commands for preparing key [here](./docs/helper.md)
