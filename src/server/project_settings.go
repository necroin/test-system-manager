package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"tsm/src/db/dbi"
	"tsm/src/logger"
	"tsm/src/settings"

	"github.com/gorilla/mux"
)

func (server *Server) ProjectSettingsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	token := params["token"]
	projectId := params["id"]

	if server.GetUserProjectRole(token, projectId) < roleCreator {
		responseWriter.Write([]byte("Permission denied"))
		return
	}

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

	if projectsResponse.Error != nil {
		responseWriter.Write([]byte(projectsResponse.Error.Error()))
		return
	}

	PageDescriptor := PageDescriptor{
		ProjectId:   projectId,
		ProjectName: projectsResponse.Records[0].Fields["Name"],
	}
	server.PageHandler(responseWriter, settings.InterfaceProjectSettingsHTML, PageDescriptor, token)
}

func (server *Server) ProjectRenameHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	token := params["token"]
	projectId := params["id"]

	if server.GetUserProjectRole(token, projectId) < roleCreator {
		responseWriter.Write([]byte("Permission denied"))
		return
	}

	data, _ := io.ReadAll(request.Body)
	newName := string(data)
	newName = strings.Trim(newName, " ")

	response := server.db.UpdateRequest(&dbi.Request{
		Table: "Projects",
		Fields: []dbi.Field{{
			Name:  "Name",
			Value: fmt.Sprintf("'%s'", newName),
		}},
		Filters: []dbi.Filter{{
			Name:     "Id",
			Operator: "=",
			Value:    projectId,
		}},
	})
	logger.Debug("%#v", response)

	if response.Error != nil {
		logger.Error(response.Error.Error())
		responseWriter.Write([]byte(response.Error.Error()))
		return
	}
}
