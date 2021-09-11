# Go Lang API structure

This project goal is to have a basic structure for Go Lang API's. There is an
effort to keep the code as simple as possible using Go standard library, but some
few libraries were also add to improve code quality and readability.

## Migration

[Migrate](https://github.com/golang-migrate/migrate) is the option for database migration management.

Two steps are required to run the migration:

1. set the PostreSQL schema path:

```bash
export POSTGRESQL_URL_GAAP=postgres://localenv:the-secret@localhost:5432/hotel?sslmode=disable
```

2. run the migration

```bash
migrate -database ${POSTGRESQL_URL_GAAP} -path store/migrations up
```

todo:

- validate current code with OpenAPI specification
- add SQL connection to PostgreSQL
- define Problem struct better
- add in code documentation

