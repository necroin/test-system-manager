package firebase

import (
	"context"
	"fmt"
	"io/ioutil"
	"tsm/src/config"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/necroin/golibs/sqlschema"
	"google.golang.org/api/option"
)

const (
	project_id   = "tsm"
	database_url = "https://tsm-local-default-rtdb.firebaseio.com/"
)

type firebaseApplication struct {
	app    *firebase.App
	config *firebase.Config
	client *db.Client
}

type DatabaseOptions struct {
	Credentials string
}

type Database struct {
	firebase  *firebaseApplication
	reference *db.Ref
	schema    map[string]sqlschema.Table
}

func New(config *config.Config) (*Database, error) {
	ctx := context.Background()

	opts := option.WithCredentialsFile(config.Credentials)

	firebaseConfig := &firebase.Config{
		ProjectID:   project_id,
		DatabaseURL: database_url,
	}

	app, err := firebase.NewApp(ctx, firebaseConfig, opts)
	if err != nil {
		return nil, fmt.Errorf("[Database] [Error] initializing app: %s", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		return nil, fmt.Errorf("[Database] [Error] initializing database client: %s", err)
	}

	reference := client.NewRef("")
	if err != nil {
		return nil, fmt.Errorf("[Database] [Error] create reference: %s", err)
	}

	firebaseAppInstance := &firebaseApplication{
		app:    app,
		config: firebaseConfig,
		client: client,
	}

	schemaFile, err := ioutil.ReadFile(config.Database.Schema)
	if err != nil {
		return nil, fmt.Errorf("[Database] [Error] failed open schema file: %s", err)
	}

	schema, err := sqlschema.Parse(schemaFile)
	if err != nil {
		return nil, fmt.Errorf("[Database] [Error] failed set schema: %s", err)
	}

	return &Database{
		firebase:  firebaseAppInstance,
		reference: reference,
		schema:    sqlschema.MapSchema(schema),
	}, nil
}

func (database *Database) Section(value string) *Database {
	return &Database{
		firebase:  database.firebase,
		reference: database.reference.Child(value),
	}
}

func (database *Database) Delete(path string, key string) error {
	err := database.reference.Child(path).Child(key).Delete(context.Background())
	if err != nil {
		return fmt.Errorf("[Database] [Error] delete data: %s", err)
	}
	return nil
}
