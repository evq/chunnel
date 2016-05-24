# Chunnel
A wrapper for docker clients that uses an SSH tunnel to securely access the
docker API. Sets the `DOCKER_HOST` environment variable.

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
