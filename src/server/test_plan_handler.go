package server

import (
	"fmt"
	"net/http"
	"os"
	"tsm/src/db/dbi"
	"tsm/src/settings"

	"github.com/gorilla/mux"
)

func (server *Server) ProjectPlanHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]
	testPlanId := params["planId"]

	style, _ := os.ReadFile(settings.InterfaceStylePath)
	script, _ := os.ReadFile(settings.InterfaceScriptPath)
	html, _ := os.ReadFile(settings.InterfacePlanHTML)

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

	content := fmt.Sprintf(
		string(html),
		fmt.Sprintf(`<style type="text/css">%s</style>`, style),
		fmt.Sprintf(
			fmt.Sprintf(`<script type="text/javascript">%s</script>`, script),
			server.url,
		),
		projectId,
		testPlanId,
		projectId,
		projectId,
		response.Records[0].Fields["Name"],
	)
	responseWriter.Write([]byte(content))
}

func (server *Server) ProjectPlanSelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
}
