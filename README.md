# compserv

[![Go](https://github.com/rhmdnd/compserv/actions/workflows/go.yml/badge.svg)](https://github.com/rhmdnd/compserv/actions/workflows/go.yml)

The name "compserv" is short for "compliance service" and its goal is to manage
and aggregate compliance information.

The service should be deployed with containers, ideally for Kubernetes-based
platforms. Also, the service will initially store compliance data for
Kubernetes deployments.

## Design Goals

1. Object models specific to compliance concepts must remain generic, so the
   service can be flexible for other implementations
2. Any tools required to run the service must be automated or scripted, and
   included in this repository with documentation

## Architecture

The service will use a relational database for persisting data. Data will be
exposed using a gRPC API.

**Why gRPC?**

We ultimately want to expose compliance data using a NIST-approved RESTful API.
Since that work is still underway, we will implement a gRPC API to enable
service-to-service communication.

If or when we have a NIST-approved API specification, we can decide how to wrap
and re-use the gRPC API.

**Why a relational database?**

Compliance data can be unbounded, making it a good candidate for object storage
or no-SQL backends. But, the service will provide an API for querying
compliance results across targets, or infrastructure. We want to enable fast
queries across all compliance targets. For example, querying a fleet of
Kubernetes clusters for which nodes have checks related to NIST AC-2.

Sanitizing the compliance data and using a relational database will better
enable this at scale, as opposed to loading compliance reports for each cluster
from object storage.

## Building

You can use the `Makefile` target to build go binaries, which are output to
`builds/` by default:

```console
$ make build
```

## Deploying

### AWS

You can deploy the database using [terraform](https://www.terraform.io/).

```console
$ cd terraform; terraform init
$ terraform apply
```

Cleanup the resources using `terraform destroy`.

### Kubernetes

Alternatively, you can deploy the database to a Kubernetes cluster with
[Kustomize](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization).
This functionality is primarily for development purposes and requires
additional work to intergate with the compliance service. In the future, we may
consider breaking the current Kustomize structure into overlays for different
environments.

```console
$ kubectl apply -k kustomize
```

Alternatively, you can use the `deploy` Makefile target to deploy the
compliance service and a PostgreSQL database into a Kubernetes cluster.

```console
$ make deploy
```

## Migrations

Please refer to the [migrations documentation](./migrations/README.md) for
instructions on creating and managing database migrations.

## gRPC API

This API is marked as **EXPERIMENTAL** and may change in backwards incompatible
ways.

Placeholder for documentation of gRPC endpoints implemented by the service.

## Releases

We will tag and branch each release from the `main` branch, following [semantic
versioning](https://semver.org/). Releases will be determined as needed since
the project doesn't follow a release cadence or schedule.

When in doubt, release early and often.

## Contributing

Please use [issues](https://github.com/rhmdnd/compserv/issues) and [pull
requests](https://github.com/rhmdnd/compserv/pulls) for contributing to the
project.
