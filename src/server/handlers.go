package server

import (
	"net/http"
	"tsm/src/settings"
)

func (server *Server) StatusHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Write([]byte(settings.ServerStatusResponse))
}
