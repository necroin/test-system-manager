package main

import (
	"fmt"
	"sync"
	"tsm/src/config"
	"tsm/src/db/sql"
	"tsm/src/logger"
	"tsm/src/server"
	"tsm/src/settings"
)

const (
	configPath = "config.yml"
)

func main() {
	config, err := config.Load(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := logger.Configure(config); err != nil {
		fmt.Println(err)
		return
	}

	db, err := sql.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	wg := sync.WaitGroup{}
	server := server.New(config.Url, db)

	server.AddHandler(settings.ServerStatusEndpoint, server.StatusHandler, "GET")

	server.AddHandler(settings.ServerProjectsEndpoint, server.ProjectsHandler, "GET")
	server.AddHandler(settings.ServerProjectsSelectEndpoint, server.ProjectsSelectHandler, "POST", "GET")

	server.AddHandler(settings.ServerProjectCasesEndpoint, server.ProjectCasesHandler, "GET")
	server.AddHandler(settings.ServerProjectCasesSelectEndpoint, server.ProjectCasesSelectHandler, "POST", "GET")

	server.AddHandler(settings.ServerProjectPlansEndpoint, server.ProjectPlansHandler, "GET")
	server.AddHandler(settings.ServerProjectPlansSelectEndpoint, server.ProjectPlansSelectHandler, "POST", "GET")

	server.AddHandler(settings.ServerProjectCaseEndpoint, server.ProjectCaseHandler, "GET")
	server.AddHandler(settings.ServerProjectCaseSelectEndpoint, server.ProjectCaseSelectHandler, "POST", "GET")
	server.AddHandler(settings.ServerProjectCaseUpdateEndpoint, server.ProjectCaseUpdateHandler, "POST", "GET")

	server.AddHandler(settings.ServerProjectPlanEndpoint, server.ProjectPlanHandler, "GET")
	server.AddHandler(settings.ServerProjectPlanSelectEndpoint, server.ProjectPlanSelectHandler, "POST", "GET")

	go func() {
		wg.Add(1)
		defer wg.Done()
		server.Start()
	}()

	if err := server.WaitStart(); err != nil {
		fmt.Println(err)
		return
	}

	wg.Wait()
}
