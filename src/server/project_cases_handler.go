package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"tsm/src/db/dbi"
	"tsm/src/logger"
	"tsm/src/settings"

	"github.com/gorilla/mux"
)

func (server *Server) ProjectCasesHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]

	style, _ := os.ReadFile(settings.InterfaceStylePath)
	script, _ := os.ReadFile(settings.InterfaceScriptPath)
	html, _ := os.ReadFile(settings.InterfaceProjectCasesHTML)

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
		projectId,
		projectId,
		projectId,
		projectsResponse.Records[0].Fields["Name"],
	)
	responseWriter.Write([]byte(content))
}

func (server *Server) ProjectCasesSelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]

	response := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_TestCase",
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

	json.NewEncoder(responseWriter).Encode(response)
}

func (server *Server) ProjectCasesInsertHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]

	name, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Error(err)
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
		logger.Error(projectsResponse.Error)
		projectsResponse.Success = false
		json.NewEncoder(responseWriter).Encode(projectsResponse)
		return
	}
}
