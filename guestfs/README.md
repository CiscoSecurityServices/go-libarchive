# Guestfs

This is just a git hosted version of the Guestfs go bindings.

## Docker
```
docker build -t guestfs .
docker run --rm -it -v $PWD:/app guestfs:latest go test ./...
```

## updating
To update the GuestFS bindings use the `update.sh` script.
