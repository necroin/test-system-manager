package sql

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"tsm/src/db/dbi"
	"tsm/src/logger"
)

func (database *Database) Insert(table string, columns []string, values []string) error {
	sqlCommand := fmt.Sprintf("INSERT OR REPLACE INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ", "), strings.Join(values, ", "))
	logger.Debug(sqlCommand)
	_, err := database.sql.Exec(sqlCommand)
	if err != nil {
		return fmt.Errorf("[Database] [Insert] [Error] failed database request: %s", err)
	}
	return nil
}

func (database *Database) InsertRequest(request *dbi.Request) *dbi.Response {
	response := &dbi.Response{
		Records: nil,
		Success: true,
		Error:   nil,
	}

	names := []string{}
	for _, field := range request.Fields {
		names = append(names, field.Name)
	}

	values := []string{}
	for _, field := range request.Fields {
		values = append(values, field.Value)
	}

	if err := database.Insert(request.Table, names, values); err != nil {
		response.Success = false
		response.Error = err
		return response
	}

	return response
}

func (database *Database) InsertHandler(responseWriter http.ResponseWriter, request *http.Request) {
	dbRequest := &dbi.Request{}
	if err := json.NewDecoder(request.Body).Decode(request); err != nil {
		logger.Error("[Database] [InsertHandler] [Error] failed decode json request: %s", err)
		return
	}

	response := database.InsertRequest(dbRequest)
	if err := json.NewEncoder(responseWriter).Encode(response); err != nil {
		logger.Error("[Database] [InsertHandler] [Error] failed encode json response: %s", err)
		return
	}

	json.NewEncoder(responseWriter).Encode(response)
}
