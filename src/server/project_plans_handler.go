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

func (server *Server) ProjectPlansHandler(responseWriter http.ResponseWriter, request *http.Request) {
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
	server.PageHandler(responseWriter, settings.InterfaceProjectPlansHTML, PageDescriptor, token)
}

func (server *Server) ProjectPlansSelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]

	projectPlansRequest := &SearchRequest{}
	json.NewDecoder(request.Body).Decode(projectPlansRequest)
	projectPlansRequest.SearchFilter = strings.Trim(projectPlansRequest.SearchFilter, " ")

	searchFilters := strings.Split(projectPlansRequest.SearchFilter, ";")
	if projectPlansRequest.SearchFilter == "" {
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
						Value:    fmt.Sprintf("'%s'", settings.ObjectTypeTestPlan),
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
				PlanId := strings.Split(record.Fields["ObjectId"], ";")[1]
				//PlanId := record.Fields["ObjectId"]
				ids = append(ids, PlanId)
			}

			dbFilter := dbi.Filter{
				Name:     "Id",
				Operator: "IN",
				Value:    fmt.Sprintf("('%s')", strings.Join(ids, "','")),
			}

			dbFilters = append(dbFilters, dbFilter)
		}
	}

	projectTestCaseResponse := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_TestPlan",
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

	if projectTestCaseResponse.Error != nil {
		json.NewEncoder(responseWriter).Encode(projectTestCaseResponse)
		return
	}

	for _, record := range projectTestCaseResponse.Records {
		testPlanId := record.Fields["Id"]

		testCaseResponse := server.db.SelectRequest(&dbi.Request{
			Table: "TSM_TestPlanTestCase",
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
				{
					Name:     "TestPlanId",
					Operator: "=",
					Value:    testPlanId,
				},
			},
		})

		if testCaseResponse.Error != nil {
			json.NewEncoder(responseWriter).Encode(projectTestCaseResponse)
			return
		}

		record.Fields["TestCaseCount"] = testCaseResponse.Records[0].Fields["Count(*)"]
	}

	json.NewEncoder(responseWriter).Encode(projectTestCaseResponse)
}

func (server *Server) ProjectPlansInsertHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	token := params["token"]
	projectId := params["id"]

	if server.GetUserProjectRole(token, projectId) < roleTester {
		responseWriter.Write([]byte("Permission denied"))
		return
	}

	name, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Error("%s", err)
		responseWriter.Write([]byte(err.Error()))
		return
	}

	projectsResponse := server.db.InsertRequest(&dbi.Request{
		Table: "TSM_TestPlan",
		Fields: []dbi.Field{
			{
				Name:  "Id",
				Value: fmt.Sprintf("(select COALESCE((MAX(Id) + 1), 1) from TSM_TestPlan where ProjectId = %s)", projectId),
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
