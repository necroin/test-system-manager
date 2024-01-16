package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tsm/src/db/dbi"
	"tsm/src/logger"
	"tsm/src/settings"

	"github.com/gorilla/mux"
)

func (server *Server) CaseTagsSelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]
	caseId := params["caseId"]

	caseTagsResponse := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_Tags",
		Fields: []dbi.Field{
			{
				Name: "Name",
			},
		},
		Filters: []dbi.Filter{
			{
				Name:     "ObjectId",
				Operator: "=",
				Value:    fmt.Sprintf("'%s;%s'", projectId, caseId),
			},
			{
				Name:     "ObjectType",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", settings.ObjectTypeTestCase),
			},
		},
	})

	if caseTagsResponse.Error != nil {
		logger.Error("%s", caseTagsResponse.Error)
		json.NewEncoder(responseWriter).Encode(caseTagsResponse)
		return
	}

	json.NewEncoder(responseWriter).Encode(caseTagsResponse)
}

func (server *Server) CaseTagsInsertHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]
	caseId := params["caseId"]
	token := params["token"]

	if server.GetUserProjectRole(token, projectId) < roleTester {
		responseWriter.Write([]byte("Permission denied"))
		return
	}

	tagName, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Error("%s", err)
		responseWriter.Write([]byte(err.Error()))
		return
	}

	response := server.db.InsertRequest(&dbi.Request{
		Table: "TSM_Tags",
		Fields: []dbi.Field{
			{
				Name:  "Name",
				Value: fmt.Sprintf("'%s'", string(tagName)),
			},
			{
				Name:  "ObjectId",
				Value: fmt.Sprintf("'%s;%s'", projectId, caseId),
			},
			{
				Name:  "ObjectType",
				Value: fmt.Sprintf("'%s'", settings.ObjectTypeTestCase),
			},
		},
		Filters: []dbi.Filter{},
	})

	if response.Error != nil {
		logger.Error("%s", response.Error)
	}
	json.NewEncoder(responseWriter).Encode(response)
}

func (server *Server) CaseTagsDeleteHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]
	caseId := params["caseId"]
	token := params["token"]

	if server.GetUserProjectRole(token, projectId) < roleTester {
		responseWriter.Write([]byte("Permission denied"))
		return
	}

	tagName, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Error("%s", err)
		responseWriter.Write([]byte(err.Error()))
		return
	}

	response := server.db.DeleteRequest(&dbi.Request{
		Table: "TSM_Tags",
		Filters: []dbi.Filter{
			{
				Name:     "Name",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", string(tagName)),
			},
			{
				Name:     "ObjectId",
				Operator: "=",
				Value:    fmt.Sprintf("'%s;%s'", projectId, caseId),
			},
			{
				Name:     "ObjectType",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", settings.ObjectTypeTestCase),
			},
		},
	})

	if response.Error != nil {
		logger.Error("%s", response.Error)
	}
	json.NewEncoder(responseWriter).Encode(response)
}
