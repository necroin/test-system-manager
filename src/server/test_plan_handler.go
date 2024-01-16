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
	"golang.org/x/exp/slices"
)

func (server *Server) ProjectPlanHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	token := params["token"]
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
	server.PageHandler(responseWriter, settings.InterfacePlanHTML, PageDescriptor, token)
}

func (server *Server) ProjectPlanSelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]
	testPlanId := params["planId"]

	projectTestPlanResponse := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_TestPlan",
		Fields: []dbi.Field{
			{
				Name: "Description",
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

	if projectTestPlanResponse.Error != nil {
		logger.Error("%s", projectTestPlanResponse.Error)
		responseWriter.Write([]byte(projectTestPlanResponse.Error.Error()))
		return
	}

	description := projectTestPlanResponse.Records[0].Fields["Description"]
	data := TestPlanDescriptor{
		Description: &description,
		TestCases:   []TestPlanDescriptorCases{},
		TestRuns:    []TestRunDescriptorCases{},
	}

	testPlanTestCaseResponse := server.db.SelectRequest(&dbi.Request{
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

	logger.Debug("%#v", testPlanTestCaseResponse)
	slices.SortFunc(testPlanTestCaseResponse.Records, func(record dbi.Record, other dbi.Record) int {
		recordsPosition := record.Fields["Position"]
		otherPosition := other.Fields["Position"]
		if recordsPosition > otherPosition {
			return 1
		}
		return -1
	})

	if testPlanTestCaseResponse.Error != nil {
		logger.Error("%s", testPlanTestCaseResponse.Error)
		json.NewEncoder(responseWriter).Encode(testPlanTestCaseResponse)
		return
	}

	for _, record := range testPlanTestCaseResponse.Records {
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
			logger.Error("%s", testCaseResponse.Error)
			json.NewEncoder(responseWriter).Encode(testCaseResponse)
			return
		}

		data.TestCases = append(data.TestCases, TestPlanDescriptorCases{
			Id:   testCaseId,
			Name: testCaseResponse.Records[0].Fields["Name"],
		})
	}

	statResponse := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_Stat",
		Fields: []dbi.Field{
			{
				Name: "TestRunId",
			},
			{
				Name: "TestCaseId",
			},
			{
				Name: "Result",
			},
			{
				Name: "Datetime",
			},
			{
				Name: "Comment",
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
	logger.Debug("%#v", statResponse)

	if statResponse.Error != nil {
		logger.Error("%s", statResponse.Error)
		json.NewEncoder(responseWriter).Encode(statResponse)
		return
	}

	for _, record := range statResponse.Records {
		data.TestRuns = append(data.TestRuns, TestRunDescriptorCases{
			Result: record.Fields["Result"],
			TestCaseId: record.Fields["TestCaseId"],
			TestRunId: record.Fields["TestRunId"],
			Datetime: record.Fields["Datetime"],
			Comment: record.Fields["Comment"],
		})
	}

	json.NewEncoder(responseWriter).Encode(data)

}

func (server *Server) ProjectPlanUpdateHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]
	testPlanId := params["planId"]

	data := &TestPlanDescriptor{}
	json.NewDecoder(request.Body).Decode(data)

	if data.Description != nil {
		response := server.db.UpdateRequest(&dbi.Request{
			Table: "TSM_TestPlan",
			Fields: []dbi.Field{{
				Name:  "Description",
				Value: fmt.Sprintf("'%s'", *data.Description),
			}},
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
			logger.Error("%s", response.Error)
			responseWriter.Write([]byte(response.Error.Error()))
			return
		}
	}

	logger.Verbose("Update ids: %v", data.TestCases)

	if len(data.TestCases) > 0 {
		response := server.db.DeleteRequest(&dbi.Request{
			Table: "TSM_TestPlanTestCase",
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

		if response.Error != nil {
			logger.Error("%s", response.Error)
			responseWriter.Write([]byte(response.Error.Error()))
			return
		}
	}

	for position, descriptor := range data.TestCases {
		response := server.db.InsertRequest(&dbi.Request{
			Table: "TSM_TestPlanTestCase",
			Fields: []dbi.Field{
				{
					Name:  "ProjectId",
					Value: projectId,
				},
				{
					Name:  "TestPlanId",
					Value: testPlanId,
				},
				{
					Name:  "TestCaseId",
					Value: descriptor.Id,
				},
				{
					Name:  "Position",
					Value: fmt.Sprintf("%d", position+1),
				},
			},
		})

		if response.Error != nil {
			logger.Error("%s", response.Error)
			responseWriter.Write([]byte(response.Error.Error()))
			return
		}
	}
}

func (server *Server) ProjectPlanCaseAppendHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]
	testPlanId := params["planId"]

	data, _ := io.ReadAll(request.Body)
	caseId := string(data)

	response := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_TestPlanTestCase",
		Fields: []dbi.Field{
			{
				Name: "COALESCE((MAX(Position) + 1), 1) as Position",
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

	if response.Error != nil {
		logger.Error("%s", response.Error)
		responseWriter.Write([]byte(response.Error.Error()))
		return
	}

	position := response.Records[0].Fields["Position"]
	logger.Verbose("[ServerProjectPlanCaseAppendHandler] append record with position %v", position)

	response = server.db.InsertRequest(&dbi.Request{
		Table: "TSM_TestPlanTestCase",
		Fields: []dbi.Field{
			{
				Name:  "ProjectId",
				Value: projectId,
			},
			{
				Name:  "TestPlanId",
				Value: testPlanId,
			},
			{
				Name:  "TestCaseId",
				Value: caseId,
			},
			{
				Name:  "Position",
				Value: position,
			},
		},
	})

	if response.Error != nil {
		logger.Error("%s", response.Error)
		responseWriter.Write([]byte(response.Error.Error()))
		return
	}
}
