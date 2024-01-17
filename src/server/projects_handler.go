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

func (server *Server) ProjectsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	token := params["token"]
	server.PageHandler(responseWriter, settings.InterfaceProjectsHTML, PageDescriptor{}, token)
}

func (server *Server) ProjectsSelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	token := params["token"]
	username := server.FindUsernameByToken(token)

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
						Value:    fmt.Sprintf("'%s'", settings.ObjectTypeProject),
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
				ids = append(ids, record.Fields["ObjectId"])
			}

			dbFilter := dbi.Filter{
				Name:     "Id",
				Operator: "IN",
				Value:    fmt.Sprintf("('%s')", strings.Join(ids, "','")),
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
		logger.Error("%s", projectsResponse.Error)
		json.NewEncoder(responseWriter).Encode(projectsResponse)
		return
	}

	projects := []dbi.Record{}
	for _, record := range projectsResponse.Records {
		projectId := record.Fields["Id"]

		role := server.FindUserProjectRole(username, projectId)
		if role == "" {
			continue
		}

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
					Value:    fmt.Sprintf("'%s'", projectId),
				},
			},
		})

		if testCaseResponse.Error != nil {
			logger.Error("%s", testCaseResponse.Error)
			json.NewEncoder(responseWriter).Encode(projectsResponse)
			return
		}

		record.Fields["TestCaseCount"] = testCaseResponse.Records[0].Fields["Count(*)"]
		projects = append(projects, record)
	}
	projectsResponse.Records = projects

	json.NewEncoder(responseWriter).Encode(projectsResponse)
}

func (server *Server) ProjectsInsertHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	token := params["token"]
	username := server.FindUsernameByToken(token)

	projectName, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Error("%s", err)
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
			{
				Name:  "Archived",
				Value: "'false'",
			},
		},
	})

	if projectsResponse.Error != nil {
		logger.Error("%s", projectsResponse.Error)
		json.NewEncoder(responseWriter).Encode(projectsResponse)
		return
	}

	projectUserResponse := server.db.InsertRequest(&dbi.Request{
		Table: "TSM_ProjectUsers",
		Fields: []dbi.Field{
			{
				Name:  "Username",
				Value: fmt.Sprintf("'%s'", username),
			},
			{
				Name:  "ProjectId",
				Value: "(select MAX(Id) from Projects)",
			},
			{
				Name:  "Role",
				Value: fmt.Sprintf("'%s'", settings.ProjectRoleCreator),
			},
		},
	})

	if projectUserResponse.Error != nil {
		logger.Error("%s", projectUserResponse.Error)
		json.NewEncoder(responseWriter).Encode(projectUserResponse)
		return
	}

	idResponse := server.db.SelectRequest(&dbi.Request{
		Table: "Projects",
		Fields: []dbi.Field{{
			Name: "MAX(Id) as Id",
		}},
	})

	if idResponse.Error != nil {
		logger.Error("%s", idResponse.Error)
		json.NewEncoder(responseWriter).Encode(idResponse)
		return
	}

	newId := idResponse.Records[0].Fields["Id"]
	responseWriter.Write([]byte(newId))
}

func (server *Server) ProjectsDeleteHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]

	response := server.db.DeleteRequest(&dbi.Request{
		Table: "Projects",
		Filters: []dbi.Filter{
			{
				Name:     "Id",
				Operator: "=",
				Value:    projectId,
			},
		},
	})

	if response.Error != nil {
		logger.Error(response.Error.Error())
		responseWriter.Write([]byte(response.Error.Error()))
		return
	}

	response = server.db.DeleteRequest(&dbi.Request{
		Table: "TSM_TestCase",
		Filters: []dbi.Filter{
			{
				Name:     "ProjectId",
				Operator: "=",
				Value:    projectId,
			},
		},
	})

	if response.Error != nil {
		logger.Error(response.Error.Error())
		responseWriter.Write([]byte(response.Error.Error()))
		return
	}

	response = server.db.DeleteRequest(&dbi.Request{
		Table: "TSM_TestPlan",
		Filters: []dbi.Filter{
			{
				Name:     "ProjectId",
				Operator: "=",
				Value:    projectId,
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
		},
	})

	if response.Error != nil {
		logger.Error(response.Error.Error())
		responseWriter.Write([]byte(response.Error.Error()))
		return
	}

	response = server.db.DeleteRequest(&dbi.Request{
		Table: "TSM_Stat",
		Filters: []dbi.Filter{
			{
				Name:     "ProjectId",
				Operator: "=",
				Value:    projectId,
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
				Value:    projectId,
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
				Operator: "like",
				Value:    fmt.Sprintf("'%s;%%'", projectId),
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
		},
	})

	if response.Error != nil {
		logger.Error(response.Error.Error())
		responseWriter.Write([]byte(response.Error.Error()))
		return
	}
}
