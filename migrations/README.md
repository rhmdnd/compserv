# Database Migrations

This folder contains database migrations for the compliance service using
[golang-migrate/migrate](https://github.com/golang-migrate/migrate).

This project requires using the `migrate` CLI to change the database schema and
to add new migrations. Please refer to the
[documentation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
for installation instructions or simply run `make tools/migrate` from the root of
this project.

## Examples

These examples assume you have a PostgreSQL server accessible with a database
named `compliance`, and you've properly
[escaped](https://github.com/golang-migrate/migrate#database-urls) the database
password. You can refer to the deployment documentation in the
[README.md](../README.md) for instructions on using terraform to deploy a
database using AWS RDS.

Apply all migrations:

```console
$ POSTGRESQL_URL="postgres://$USER:$PASSWORD@$DATABASE_ENDPOINT/compliance
$ ./migrate -database $POSTGRESQL_URL -path migrations up
```

Apply all down migrations:

```console
$ ./migrate -database $POSTGRESQL_URL -path migrations down
```

## Test

To run a quick and easy test of the migrations, use the `test.sh` script:
```console
bash migrations/test.sh
```

If you don't have the `migrate` CLI tool installed in your `$PATH`, run:
```console
make tools/migrate
PATH=$PATH:$(pwd)/tools bash migrations/test.sh
```
