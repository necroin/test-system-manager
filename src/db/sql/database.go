package sql

import (
	"database/sql"
	"fmt"
	"os"
	"tsm/src/config"

	_ "github.com/mattn/go-sqlite3"
	"github.com/necroin/golibs/sqlschema"
)

type Database struct {
	config *config.Config
	sql    *sql.DB
}

func New(config *config.Config) (*Database, error) {
	db, err := sql.Open("sqlite3", config.Database.Storage)
	if err != nil {
		return nil, fmt.Errorf("[Database] [Error] failed open database: %s", err)
	}

	schemaFile, err := os.Open(config.Database.Schema)
	if err != nil {
		return nil, fmt.Errorf("[Database] [Error] failed open schema file: %s", err)
	}
	defer schemaFile.Close()

	if err := sqlschema.SetSchema(db, schemaFile); err != nil {
		return nil, fmt.Errorf("[Database] [Error] failed set schema: %s", err)
	}

	return &Database{
		config: config,
		sql:    db,
	}, nil
}

func (database *Database) Close() {
	database.sql.Close()
}
