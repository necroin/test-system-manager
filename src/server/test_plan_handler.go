package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tsm/src/db/dbi"
	"tsm/src/logger"
	"tsm/src/settings"

	"github.com/gorilla/mux"
	"golang.org/x/exp/slices"
)

func (server *Server) ProjectPlanHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]
	testPlanId := params["planId"]

	response := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_TestPlan",
		Fields: []dbi.Field{
			{
				Name: "Name",
			},
		},
		Filters: []dbi.Filter{
			{
				Name:     "Id",
				Operator: "=",
				Value:    testPlanId,
			},
			{
				Name:     "ProjectId",
				Operator: "=",
				Value:    projectId,
			},
		},
	})

	if response.Error != nil {
		responseWriter.Write([]byte(response.Error.Error()))
		return
	}

	PageDescriptor := PageDescriptor{
		ProjectId:    projectId,
		TestPlanId:   testPlanId,
		TestPlanName: response.Records[0].Fields["Name"],
	}
	server.PageHandler(responseWriter, settings.InterfacePlanHTML, PageDescriptor)
}

func (server *Server) ProjectPlanSelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]
	testPlanId := params["planId"]

	projectTestPlanResponse := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_TestPlanTestCase",
		Fields: []dbi.Field{
			{
				Name: "TestCaseId",
			},
			{
				Name: "Position",
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

	slices.SortFunc(projectTestPlanResponse.Records, func(record dbi.Record, other dbi.Record) int {
		recordsPosition := record.Fields["Position"]
		otherPosition := other.Fields["Position"]
		if recordsPosition > otherPosition {
			return 1
		}
		return -1
	})

	if projectTestPlanResponse.Error != nil {
		json.NewEncoder(responseWriter).Encode(projectTestPlanResponse)
		return
	}

	for _, record := range projectTestPlanResponse.Records {
		testCaseId := record.Fields["TestCaseId"]

		testCaseResponse := server.db.SelectRequest(&dbi.Request{
			Table: "TSM_TestCase",
			Fields: []dbi.Field{
				{
					Name: "Name",
				},
			},
			Filters: []dbi.Filter{
				{
					Name:     "Id",
					Operator: "=",
					Value:    testCaseId,
				},
				{
					Name:     "ProjectId",
					Operator: "=",
					Value:    projectId,
				},
			},
		})

		if testCaseResponse.Error != nil {
			json.NewEncoder(responseWriter).Encode(testCaseResponse)
			return
		}

		record.Fields["TestCaseName"] = testCaseResponse.Records[0].Fields["Name"]
	}

	json.NewEncoder(responseWriter).Encode(projectTestPlanResponse)

}

func (server *Server) ProjectPlanUpdateHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]
	testPlanId := params["planId"]

	data := &TestPlanDescriptor{}
	json.NewDecoder(request.Body).Decode(data)
	logger.Verbose("Update ids: %v", data.TestCases)

	for position, id := range data.TestCases {
		server.db.UpdateRequest(&dbi.Request{
			Table: "TSM_TestPlanTestCase",
			Fields: []dbi.Field{{
				Name:  "Position",
				Value: fmt.Sprintf("%d", position+1),
			}},
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
				{
					Name:     "TestCaseId",
					Operator: "=",
					Value:    fmt.Sprintf("%d", id),
				},
			},
		})
	}
}
