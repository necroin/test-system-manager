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

	"github.com/gorilla/mux"
)

func (server *Server) ProjectCasesHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	token := params["token"]

	projectId := params["id"]

	projectsResponse := server.db.SelectRequest(&dbi.Request{
		Table: "Projects",
		Fields: []dbi.Field{
			{
				Name: "Name",
			},
		},
		Filters: []dbi.Filter{
			{
				Name:     "Id",
				Operator: "=",
				Value:    projectId,
			},
		},
	})

	if projectsResponse.Error != nil {
		responseWriter.Write([]byte(projectsResponse.Error.Error()))
		return
	}

	PageDescriptor := PageDescriptor{
		ProjectId:   projectId,
		ProjectName: projectsResponse.Records[0].Fields["Name"],
	}
	server.PageHandler(responseWriter, settings.InterfaceProjectCasesHTML, PageDescriptor, token)
}

func (server *Server) ProjectCasesSelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]

	projectCasesRequest := &SearchRequest{}
	json.NewDecoder(request.Body).Decode(projectCasesRequest)
	projectCasesRequest.SearchFilter = strings.Trim(projectCasesRequest.SearchFilter, " ")

	searchFilters := strings.Split(projectCasesRequest.SearchFilter, ";")
	if projectCasesRequest.SearchFilter == "" {
		searchFilters = []string{}
	}

	dbFilters := []dbi.Filter{
		{
			Name:     "ProjectId",
			Operator: "=",
			Value:    projectId,
		},
	}

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

		if filterType == "tag" {
			tagsRepsonse := server.db.SelectRequest(&dbi.Request{
				Table: "TSM_Tags",
				Fields: []dbi.Field{
					{
						Name: "ObjectId",
					},
				},
				Filters: []dbi.Filter{
					{
						Name:     "Name",
						Operator: "=",
						Value:    fmt.Sprintf("'%s'", string(filterValue)),
					},
					{
						Name:     "ObjectType",
						Operator: "=",
						Value:    fmt.Sprintf("'%s'", settings.ObjectTypeTestCase),
					},
				},
			})

			if tagsRepsonse.Error != nil {
				logger.Error("%s", tagsRepsonse.Error)
				json.NewEncoder(responseWriter).Encode(tagsRepsonse)
				return
			}

			ids := []string{}

			for _, record := range tagsRepsonse.Records {
				caseId := strings.Split(record.Fields["ObjectId"], ";")[1]
				//caseId := record.Fields["ObjectId"]
				ids = append(ids, caseId)
			}

			dbFilter := dbi.Filter{
				Name:     "Id",
				Operator: "IN",
				Value:    fmt.Sprintf("('%s')", strings.Join(ids, "','")),
			}

			dbFilters = append(dbFilters, dbFilter)
		}
	}

	caseResponse := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_TestCase",
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

	if caseResponse.Error != nil {
		logger.Error("%s", caseResponse.Error)
		json.NewEncoder(responseWriter).Encode(caseResponse)
		return
	}

	json.NewEncoder(responseWriter).Encode(caseResponse)
}

func (server *Server) ProjectCasesInsertHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]

	name, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Error("%s", err)
		responseWriter.Write([]byte(err.Error()))
		return
	}

	projectsResponse := server.db.InsertRequest(&dbi.Request{
		Table: "TSM_TestCase",
		Fields: []dbi.Field{
			{
				Name:  "Id",
				Value: fmt.Sprintf("(select COALESCE((MAX(Id) + 1), 1) from TSM_TestCase where ProjectId = %s)", projectId),
			},
			{
				Name:  "ProjectId",
				Value: projectId,
			},
			{
				Name:  "Name",
				Value: fmt.Sprintf("'%s'", string(name)),
			},
		},
	})

	if projectsResponse.Error != nil {
		logger.Error("%s", projectsResponse.Error)
		json.NewEncoder(responseWriter).Encode(projectsResponse)
		return
	}
}

func (server *Server) ProjectCasesDeleteHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]
	testCaseId := params["caseId"]

	response := server.db.DeleteRequest(&dbi.Request{
		Table: "TSM_TestCase",
		Filters: []dbi.Filter{
			{
				Name:     "ProjectId",
				Operator: "=",
				Value:    projectId,
			},
			{
				Name:     "Id",
				Operator: "=",
				Value:    testCaseId,
			},
		},
	})

	if response.Error != nil {
		logger.Error(response.Error.Error())
		responseWriter.Write([]byte(response.Error.Error()))
		return
	}

	response = server.db.DeleteRequest(&dbi.Request{
		Table: "TSM_TestPlanTestCase",
		Filters: []dbi.Filter{
			{
				Name:     "ProjectId",
				Operator: "=",
				Value:    projectId,
			},
			{
				Name:     "TestCaseId",
				Operator: "=",
				Value:    testCaseId,
			},
		},
	})

	if response.Error != nil {
		logger.Error(response.Error.Error())
		responseWriter.Write([]byte(response.Error.Error()))
		return
	}

	response = server.db.DeleteRequest(&dbi.Request{
		Table: "TSM_Tags",
		Filters: []dbi.Filter{
			{
				Name:     "ObjectId",
				Operator: "=",
				Value:    fmt.Sprintf("'%s;%s'", projectId, testCaseId),
			},
			{
				Name:     "ObjectType",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", settings.ObjectTypeTestCase),
			},
		},
	})

	if response.Error != nil {
		logger.Error(response.Error.Error())
		responseWriter.Write([]byte(response.Error.Error()))
		return
	}

	response = server.db.DeleteRequest(&dbi.Request{
		Table: "TSM_Comments",
		Filters: []dbi.Filter{
			{
				Name:     "ProjectId",
				Operator: "=",
				Value:    projectId,
			},
			{
				Name:     "ObjectId",
				Operator: "=",
				Value:    testCaseId,
			},
			{
				Name:     "ObjectType",
				Operator: "=",
				Value:    fmt.Sprintf("'case'"),
			},
		},
	})

	if response.Error != nil {
		logger.Error(response.Error.Error())
		responseWriter.Write([]byte(response.Error.Error()))
		return
	}
}
