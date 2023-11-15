# Usefull commands

## OpenSSL

Converting pfx to pem:

```
openssl pkcs12 -in your.pfx -out your.pem -nodes
```

To extract the private key form a PFX to a PEM file use this command:

```
openssl pkcs12 -in your.pfx -nocerts -out encrypted.pem
```

To receive unencrypted private key in PEM format from encrypted, use this command:

```
openssl rsa -in encrypted.pem -out unencrypted.pem
```

## Powershell

To encode files to base64 you can use this powershell command:

```
 [convert]::ToBase64String((Get-Content -path "unencrypted.pem" -Encoding byte))
```

