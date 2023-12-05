package sql

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"tsm/src/db/dbi"
	"tsm/src/logger"
)

func (database *Database) UpdateRequest(request *dbi.Request) *dbi.Response {
	response := &dbi.Response{
		Records: nil,
		Success: true,
		Error:   nil,
	}

	sqlCommand := fmt.Sprintf("UPDATE %s SET ", request.Table)

	fields := []string{}
	for _, field := range request.Fields {
		fields = append(fields, fmt.Sprintf("%s = %s", field.Name, field.Value))
	}
	sqlCommand = sqlCommand + strings.Join(fields, ", ")

	if len(request.Filters) != 0 {
		filters := []string{}
		for _, filter := range request.Filters {
			filters = append(filters, fmt.Sprintf("%s %s %s", filter.Name, filter.Operator, filter.Value))
		}
		sqlFilters := " WHERE " + strings.Join(filters, " AND ")
		sqlCommand = sqlCommand + sqlFilters
	}

	logger.Verbose(sqlCommand)
	_, err := database.sql.Exec(sqlCommand)
	if err != nil {
		response.Success = false
		response.Error = fmt.Errorf("[Database] [UpdateRequest] [Error] failed database request: %s", err)
	}

	return response
}

func (database *Database) UpdateHandler(responseWriter http.ResponseWriter, request *http.Request) {
	dbRequest := &dbi.Request{}
	if err := json.NewDecoder(request.Body).Decode(dbRequest); err != nil {
		logger.Error(fmt.Errorf("[Database] [UpdateHandler] [Error] failed decode json request: %s", err))
		return
	}

	response := database.UpdateRequest(dbRequest)
	if err := json.NewEncoder(responseWriter).Encode(response); err != nil {
		logger.Error(fmt.Errorf("[Database] [UpdateHandler] [Error] failed encode json response: %s", err))
		return
	}

	json.NewEncoder(responseWriter).Encode(response)
}
