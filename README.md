# Auth Service

This is the Auth service

Generated with

```
micro new --namespace=DMS.SMS.v1 --type=service auth
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: DMS.SMS.v1.service.auth
- Type: service
- Alias: auth

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

```
# install etcd
brew install etcd

# run etcd
etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./auth-service
```

Build a docker image
```
make docker
```