package server

import (
	"encoding/json"
	"fmt"
	"strconv"
	"net/http"
	"tsm/src/db/dbi"
	"tsm/src/logger"
	"tsm/src/settings"

	"github.com/gorilla/mux"
	"golang.org/x/exp/slices"
)

func (server *Server) ProjectPlanRunHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]
	testPlanId := params["planId"]
	token := params["token"]
	testRunId := 0

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

	projectTestRunResponse := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_Stat",
		Fields: []dbi.Field{
			{
				Name: "TestRunId",
			},
		},
		Filters: []dbi.Filter{
			{
				Name:     "TestRunId",
				Operator: "=",
				Value:    fmt.Sprintf("(select COALESCE((MAX(TestRunId)), 1) from TSM_Stat)"),
			},
		},
	})

	if projectTestRunResponse.Error != nil {
		logger.Error("%s", projectTestRunResponse.Error)
		responseWriter.Write([]byte(projectTestRunResponse.Error.Error()))
		return
	}

	if len(projectTestRunResponse.Records) > 0  {
		testRunId,_ = strconv.Atoi(projectTestRunResponse.Records[0].Fields["TestRunId"])
		testRunId = testRunId + 1
	} else{
		testRunId = 1
	}

	PageDescriptor := PageDescriptor{
		ProjectId:    projectId,
		TestPlanId:   testPlanId,
		TestPlanName: response.Records[0].Fields["Name"],
		TestRunId:    fmt.Sprintf("%d", testRunId),
	}
	server.PageHandler(responseWriter, settings.InterfacePlanRunHTML, PageDescriptor, token)
}


func (server *Server) ProjectPlanRunSelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
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
	data := TestPlanRunDescriptor{
		Description: &description,
		TestCases:   []TestPlanRunDescriptorCases{},
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
				{
					Name: "Scenario",
				},
				{
					Name: "Description",
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

		data.TestCases = append(data.TestCases, TestPlanRunDescriptorCases{
			Id:   testCaseId,
			Name: testCaseResponse.Records[0].Fields["Name"],
			Scenario: testCaseResponse.Records[0].Fields["Scenario"],
			Description: testCaseResponse.Records[0].Fields["Description"],
		})
	}

	json.NewEncoder(responseWriter).Encode(data)

}

func (server *Server) ProjectPlanRunInsertHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]
	testPlanId := params["planId"]

	data := &TestPlanStatDescriptor{}
	json.NewDecoder(request.Body).Decode(data)


	for _, descriptor := range data.TestCases {
		response := server.db.InsertRequest(&dbi.Request{
			Table: "TSM_Stat",
			Fields: []dbi.Field{{
				Name:  "Result",
				Value:  fmt.Sprintf("'%s'",descriptor.Result),
			},
			{
				Name: "DateTime",
				Value:  fmt.Sprintf("'%s'",descriptor.DateTime),
			},	
			
			{
				Name:     "ProjectId",
				Value:    projectId,
			},
			{
				Name:     "TestPlanId",
				Value:    testPlanId,
			},
			{
				Name:     "TestCaseId",
				Value:    fmt.Sprintf("%s", descriptor.Id),
			},
			{
				Name: "TestRunId",
				Value: fmt.Sprintf("%s", descriptor.RunId),
			},
			{
				Name: "Comment",
				Value: fmt.Sprintf("'%s'", descriptor.Comment),
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