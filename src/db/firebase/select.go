package firebase

import (
	"context"
	"fmt"
	"tsm/src/db/dbi"
)

func (database *Database) SelectAll(table string) ([]map[string]string, error) {
	result := []map[string]string{}
	err := database.reference.Child(table).Get(context.Background(), &result)
	if err != nil {
		return nil, fmt.Errorf("[Database] [SelectRequest] [Error] failed database request: %s", err)
	}
	return result, nil
}

func (database *Database) SelectWithFilter(table string, field string, value string) (map[string]string, error) {
	result := map[string]string{}
	err := database.reference.Child(table).OrderByChild(field).EqualTo(value).Get(context.Background(), &result)
	if err != nil {
		return nil, fmt.Errorf("[Database] [Error] select failed: %s", err)
	}
	return result, nil
}

func (database *Database) SelectRequest(request *dbi.Request) *dbi.Response {
	response := &dbi.Response{
		Records: []dbi.Record{},
		Success: true,
		Error:   nil,
	}

	records, err := database.SelectAll(request.Table)
	if err != nil {
		response.Success = false
		response.Error = err
		return response
	}

	if len(request.Filters) != 0 {

	} else {
		for _, record := range records {
			response.Records = append(response.Records, dbi.Record{Fields: record})
		}
	}

	return response
}
