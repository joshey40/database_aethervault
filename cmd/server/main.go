package main

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5"
	"github.com/joshey40/database_aethervault/internal/api"
	"github.com/joshey40/database_aethervault/internal/config"
	dbconn "github.com/joshey40/database_aethervault/internal/db/gen"
	"github.com/joshey40/database_aethervault/internal/logger"
	"go.uber.org/zap"
)

func run() error {
	ctx := context.Background()

	config, err := config.LoadConf("../../config.yml")

	if err != nil {
		logger.L().Error("Error loading config", zap.Error(err))
		return err
	}
	connectionString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=REPLACEME", config.DB.Host, config.DB.Port, config.DB.User, config.DB.Name)
	logger.L().Info("", zap.String("Connection String", connectionString))
	conn, err := pgx.Connect(ctx, connectionString)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	queries := dbconn.New(conn)

	err = queries.CreateUserTable(ctx)
	if err != nil {
		return err
	}

	// list all users
	authors, err := queries.ListUsers(ctx)
	if err != nil {
		return err
	}
	logger.L().Sugar().Info(authors)

	pwhash, err := argon2id.CreateHash("VerySavePassword", argon2id.DefaultParams)
	if err != nil {
		return err
	}
	userinfo := dbconn.CreateUserParams{
		Username: "I am a cool user",
		Pwhash:   pwhash,
	}

	// create an user
	insertedUser, err := queries.CreateUser(ctx, userinfo)
	if err != nil {
		return err
	}
	logger.L().Sugar().Info(insertedUser)

	// get the user we just inserted by name
	fetchedUser, err := queries.GetUserByName(ctx, insertedUser.Username)
	if err != nil {
		return err
	}

	// get the user we just inserted by uuid
	fetchedUser2, err := queries.GetUserByUUID(ctx, insertedUser.Uuid)
	if err != nil {
		return err
	}

	// prints true
	logger.L().Sugar().Info(reflect.DeepEqual(insertedUser, fetchedUser))
	logger.L().Sugar().Info(reflect.DeepEqual(insertedUser, fetchedUser2))

	return nil
}

func main() {
	r := api.Router()
	http.ListenAndServe(":3333", r)

	// err := run()
	// if err != nil {
	// 	logger.L().Error("Some kind of error", zap.Error(err))
	// }
}
