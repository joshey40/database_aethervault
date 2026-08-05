# Aethervault database backend
This project serves as a backend to the aethervault app [link](https://github.com/joshey40/flutter_app_aethervault). It needs to run in combination with a postgresql database.

## Generate SQL Code
This project uses **sqlc** for generating typesafe go functions from SQL. The SQL files needed for sqlc are located in `internal/db/sqlc` while the generated code is in `internal/db/gen`. 

To regenerate the code (i.e. after a change), use the command 

```bash
sqlc generate -f internal/db/sqlc/sqlc.yml
```
