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

func (server *Server) PlanTagsSelectHandler(responseWriter http.ResponseWriter, request *http.Request) {

	params := mux.Vars(request)
	projectId := params["id"]
	PlanId := params["planId"]

	PlanTagsResponse := server.db.SelectRequest(&dbi.Request{
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
				Value:    fmt.Sprintf("'%s;%s'", projectId, PlanId),
			},
			{
				Name:     "ObjectType",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", settings.ObjectTypeTestPlan),
			},
		},
	})

	if PlanTagsResponse.Error != nil {
		logger.Error("%s", PlanTagsResponse.Error)
		json.NewEncoder(responseWriter).Encode(PlanTagsResponse)
		return
	}

	json.NewEncoder(responseWriter).Encode(PlanTagsResponse)
}

func (server *Server) PlanTagsInsertHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]
	PlanId := params["planId"]

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
				Value: fmt.Sprintf("'%s;%s'", projectId, PlanId),
			},
			{
				Name:  "ObjectType",
				Value: fmt.Sprintf("'%s'", settings.ObjectTypeTestPlan),
			},
		},
		Filters: []dbi.Filter{},
	})

	if response.Error != nil {
		logger.Error("%s", response.Error)
	}
	json.NewEncoder(responseWriter).Encode(response)
}

func (server *Server) PlanTagsDeleteHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]
	PlanId := params["planId"]

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
				Value:    fmt.Sprintf("'%s;%s'", projectId, PlanId),
			},
			{
				Name:     "ObjectType",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", settings.ObjectTypeTestPlan),
			},
		},
	})

	if response.Error != nil {
		logger.Error("%s", response.Error)
	}
	json.NewEncoder(responseWriter).Encode(response)
}
