package server

import (
	"fmt"
	"net/http"
	"os"
	"tsm/src/db/dbi"
	"tsm/src/settings"

	"github.com/gorilla/mux"
)

func (server *Server) ProjectSettingsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]

	style, _ := os.ReadFile(settings.InterfaceStylePath)
	script, _ := os.ReadFile(settings.InterfaceScriptPath)
	html, _ := os.ReadFile(settings.InterfaceProjectSettingsHTML)

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
		projectId,
	)
	responseWriter.Write([]byte(content))
}
