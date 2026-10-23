# Aethervault database backend
This project serves as a backend to the aethervault app [link](https://github.com/joshey40/flutter_app_aethervault). It needs to run in combination with a postgresql database.

# ENV file
The server needs a `.env` file in the root directory. It should include

```bash
DBPASSWD=[the db password]
JWT_SECRET=[the secret to generate JWTs]
```

## Generate SQL Code
This project uses **sqlc** for generating typesafe go functions from SQL. The SQL files needed for sqlc are located in `internal/db/sqlc` while the generated code is in `internal/db/gen`. 

To regenerate the code (i.e. after a change), use the command 

```bash
sqlc generate -f internal/db/sqlc/sqlc.yml
```

## Add database migration
Do not change the existing migrations. Every change is done with another migration step.

To add a new migration, execute this command in the migration folder
```bash
goose -s create [name] sql
```