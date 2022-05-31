# Database Migrations

This folder contains database migrations for the compliance service using
[golang-migrate/migrate](https://github.com/golang-migrate/migrate).

This project requires using the `migrate` CLI to change the database schema and
to add new migrations. Please refer to the
[documentation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
for installation instructions.

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
