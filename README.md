# Chunnel
A wrapper for docker clients that uses an SSH tunnel to access a remote
hosts docker API. Sets the `DOCKER_HOST` / `DOCKER_CERT_PATH` and unsets 
the `DOCKER_TLS_VERIFY` environment variables.

## Usage
```
go get github.com/evq/chunnel
```

## Examples
```
chunnel ubuntu@foo.bar docker ps
chunnel ubuntu@foo.bar docker images
chunnel ubuntu@foo.bar docker-compose up
```

## Limitations
Note that currently only a insecure docker API locally bound to port 2375 
(on the remote host) is supported.
