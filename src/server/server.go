package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"
	"tsm/src/db/dbi"
	"tsm/src/logger"
	"tsm/src/settings"

	"github.com/gorilla/mux"
)

type Server struct {
	url      string
	router   *mux.Router
	instance *http.Server
	db       dbi.Database
}

func New(url string, db dbi.Database) *Server {
	router := mux.NewRouter()
	return &Server{
		url:    url,
		router: router,
		instance: &http.Server{
			Addr:    url,
			Handler: router,
		},
		db: db,
	}
}

func (server *Server) Start() {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		server.instance.Shutdown(ctx)
	}()

	server.instance.ListenAndServe()
}

func (server *Server) AddHandler(path string, handler func(http.ResponseWriter, *http.Request), methods ...string) {
	server.router.HandleFunc(path, handler).Methods(methods...)
}

func (server *Server) WaitStart() error {
	client := http.Client{}
	for i := 0; i < settings.ServerWaitStartRepeatCount; i++ {
		request, _ := http.NewRequest(
			http.MethodGet,
			"http://"+server.url+settings.ServerStatusEndpoint,
			bytes.NewReader([]byte("")),
		)

		response, err := client.Do(request)
		if err != nil {
			time.Sleep(settings.ServerWaitStartSleepSeconds * time.Second)
			continue
		}
		data, err := io.ReadAll(response.Body)
		if err != nil {
			time.Sleep(settings.ServerWaitStartSleepSeconds * time.Second)
			continue
		}
		if string(data) == settings.ServerStatusResponse {
			return nil
		}
	}
	return fmt.Errorf("[Server] [WaitStart] [Error] failed get server status")
}

func (server *Server) PageHandler(responseWriter http.ResponseWriter, htmlPath string, pageDescriptor PageDescriptor, token string) {
	pageDescriptor.Url = server.url + "/" + token

	style, _ := os.ReadFile(settings.InterfaceStylePath)
	pageDescriptor.Style = fmt.Sprintf(`<style type="text/css">%s</style>`, style)

	script, _ := os.ReadFile(settings.InterfaceScriptPath)
	pageDescriptor.Script = fmt.Sprintf(fmt.Sprintf(`<script type="text/javascript">%s</script>`, script), pageDescriptor.Url)

	html, _ := os.ReadFile(htmlPath)
	pageTemplate, _ := template.New("HtmpPage").Parse(string(html))

	pageTemplate.Execute(responseWriter, pageDescriptor)
}

func (server *Server) FindUsernameByToken(token string) string {
	response := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_Users",
		Fields: []dbi.Field{{
			Name: "Username",
		}},
		Filters: []dbi.Filter{{
			Name:     "Token",
			Operator: "=",
			Value:    fmt.Sprintf("'%s'", token),
		}},
	})

	logger.Debug("%#v", response)

	if response.Error != nil {
		logger.Error("%s", response.Error)
		return ""
	}

	if len(response.Records) == 0 {
		return ""
	}

	username := response.Records[0].Fields["Username"]
	return username
}

func (server *Server) FindUserProjectRole(username string, projectId string) string {
	response := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_ProjectUsers",
		Fields: []dbi.Field{{
			Name: "Role",
		}},
		Filters: []dbi.Filter{
			{
				Name:     "Username",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", username),
			},
			{
				Name:     "ProjectId",
				Operator: "=",
				Value:    projectId,
			},
		},
	})

	logger.Debug("%#v", response)

	if response.Error != nil {
		logger.Error("%s", response.Error)
		return ""
	}

	if len(response.Records) == 0 {
		return ""
	}

	role := response.Records[0].Fields["Role"]
	return role
}

func (server *Server) GetUserProjectRole(token string, projectId string) int {
	username := server.FindUsernameByToken(token)
	role := server.FindUserProjectRole(username, projectId)
	if role == "Создатель" {
		return roleCreator
	}

	if role == "Тестировщик" {
		return roleTester
	}

	if role == "Аналитик" {
		return roleAnalisyst
	}
	return 0
}

var (
	roleCreator   = 3
	roleTester    = 2
	roleAnalisyst = 1
)
