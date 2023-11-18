package firebase

import (
	"context"
	"fmt"
	"strings"
	"tsm/src/db/dbi"

	"golang.org/x/exp/slices"
)

func (database *Database) Insert(table string, key string, record map[string]string) error {
	err := database.reference.Child(table).Child(key).Set(context.Background(), record)
	if err != nil {
		return fmt.Errorf("[Database] [Error] insert failed: %s", err)
	}

	return nil
}

func (database *Database) InsertRequest(request *dbi.Request) *dbi.Response {
	response := &dbi.Response{
		Records: nil,
		Success: true,
		Error:   nil,
	}

	schemaTable := database.schema[request.Table]

	record := map[string]string{}
	keyValues := []string{}
	for _, field := range request.Fields {
		record[field.Name] = field.Value
		if slices.Contains(schemaTable.PrimaryKey, field.Name) {
			keyValues = append(keyValues, field.Value)
		}
	}

	database.Insert(request.Table, strings.Join(keyValues, ","), record)
	return response
}
