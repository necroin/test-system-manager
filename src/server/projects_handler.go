package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"tsm/src/db/dbi"
	"tsm/src/logger"
	"tsm/src/settings"
)

func (server *Server) ProjectsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	server.PageHandler(responseWriter, settings.InterfaceProjectsHTML)
}

func (server *Server) ProjectsSelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	projectsRequest := &SearchRequest{}
	json.NewDecoder(request.Body).Decode(projectsRequest)
	projectsRequest.SearchFilter = strings.Trim(projectsRequest.SearchFilter, " ")

	searchFilters := strings.Split(projectsRequest.SearchFilter, ";")
	if projectsRequest.SearchFilter == "" {
		searchFilters = []string{}
	}

	dbFilters := []dbi.Filter{}
	for _, searchFilter := range searchFilters {
		filterElements := strings.Split(searchFilter, "=")
		if len(filterElements) != 2 {
			continue
		}
		filterType := filterElements[0]
		filterValue := filterElements[1]

		filterType = strings.Trim(filterType, " ")
		filterType = strings.ToLower(filterType)

		filterValue = strings.Trim(filterValue, " ")
		filterValue = strings.Trim(filterValue, "\"")

		if filterType == "name" {
			dbFilter := dbi.Filter{
				Name:     "Name",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", filterValue),
			}
			dbFilters = append(dbFilters, dbFilter)
		}
	}

	projectsResponse := server.db.SelectRequest(&dbi.Request{
		Table: "Projects",
		Fields: []dbi.Field{
			{
				Name: "Id",
			},
			{
				Name: "Name",
			},
		},
		Filters: dbFilters,
	})

	if projectsResponse.Error != nil {
		logger.Error(projectsResponse.Error)
		json.NewEncoder(responseWriter).Encode(projectsResponse)
		return
	}

	for _, record := range projectsResponse.Records {
		projectId := record.Fields["Id"]

		testCaseResponse := server.db.SelectRequest(&dbi.Request{
			Table: "TSM_TestCase",
			Fields: []dbi.Field{
				{
					Name: "Count(*)",
				},
			},
			Filters: []dbi.Filter{
				{
					Name:     "ProjectId",
					Operator: "=",
					Value:    projectId,
				},
			},
		})

		if testCaseResponse.Error != nil {
			logger.Error(testCaseResponse.Error)
			json.NewEncoder(responseWriter).Encode(projectsResponse)
			return
		}

		record.Fields["TestCaseCount"] = testCaseResponse.Records[0].Fields["Count(*)"]
	}

	json.NewEncoder(responseWriter).Encode(projectsResponse)
}

func (server *Server) ProjectsInsertHandler(responseWriter http.ResponseWriter, request *http.Request) {
	projectName, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Error(err)
		responseWriter.Write([]byte(err.Error()))
		return
	}

	projectsResponse := server.db.InsertRequest(&dbi.Request{
		Table: "Projects",
		Fields: []dbi.Field{
			{
				Name:  "Id",
				Value: "(select COALESCE((MAX(Id) + 1), 1) from Projects)",
			},
			{
				Name:  "Name",
				Value: fmt.Sprintf("'%s'", string(projectName)),
			},
		},
	})

	if projectsResponse.Error != nil {
		logger.Error(projectsResponse.Error)
		json.NewEncoder(responseWriter).Encode(projectsResponse)
		return
	}
}
