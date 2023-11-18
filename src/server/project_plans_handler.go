package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"tsm/src/db/dbi"
	"tsm/src/settings"

	"github.com/gorilla/mux"
)

func (server *Server) ProjectPlansHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]

	style, _ := os.ReadFile(settings.InterfaceStylePath)
	script, _ := os.ReadFile(settings.InterfaceScriptPath)
	html, _ := os.ReadFile(settings.InterfaceProjectPlansHTML)

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
	content := fmt.Sprintf(
		string(html),
		fmt.Sprintf(`<style type="text/css">%s</style>`, style),
		fmt.Sprintf(
			fmt.Sprintf(`<script type="text/javascript">%s</script>`, script),
			server.url,
		),
		projectId,
		projectId,
		projectId,
		projectsResponse.Records[0].Fields["Name"],
	)
	responseWriter.Write([]byte(content))
}

func (server *Server) ProjectPlansSelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]

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
		Filters: []dbi.Filter{
			{
				Name:     "ProjectId",
				Operator: "=",
				Value:    projectId,
			},
		},
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
					Name:     "TestPlanId",
					Operator: "=",
					Value:    testPlanId,
				},
			},
		})

		if testCaseResponse.Error != nil {
			projectTestCaseResponse.Success = false
			json.NewEncoder(responseWriter).Encode(projectTestCaseResponse)
			return
		}

		record.Fields["TestCaseCount"] = testCaseResponse.Records[0].Fields["Count(*)"]
	}

	json.NewEncoder(responseWriter).Encode(projectTestCaseResponse)
}
