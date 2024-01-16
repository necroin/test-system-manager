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
	server.AddHandler(settings.ServerAuthEndpoint, server.AuthHandler, "GET")
	server.AddHandler(settings.ServerAuthRegisterEndpoint, server.AuthRegisterHandler, "POST")
	server.AddHandler(settings.ServerAuthTokenEndpoint, server.AuthTokenHandler, "POST")

	server.AddHandler(settings.ServerProjectsEndpoint, server.ProjectsHandler, "GET")
	server.AddHandler(settings.ServerProjectsSelectEndpoint, server.ProjectsSelectHandler, "POST", "GET")
	server.AddHandler(settings.ServerProjectsInsertEndpoint, server.ProjectsInsertHandler, "POST", "GET")

	server.AddHandler(settings.ServerProjectTagsSelectEndpoint, server.ProjectTagsSelectHandler, "POST", "GET")
	server.AddHandler(settings.ServerProjectTagsInsertEndpoint, server.ProjectTagsInsertHandler, "POST", "GET")
	server.AddHandler(settings.ServerProjectTagsDeleteEndpoint, server.ProjectTagsDeleteHandler, "POST", "GET")

	server.AddHandler(settings.ServerProjectCommentsSelectEndpoint, server.ProjectCommentsSelectHandler, "POST", "GET")
	server.AddHandler(settings.ServerProjectCommentsInsertEndpoint, server.ProjectCommentsInsertHandler, "POST", "GET")
	server.AddHandler(settings.ServerProjectCommentsDeleteEndpoint, server.ProjectCommentsDeleteHandler, "POST", "GET")

	server.AddHandler(settings.ServerProjectCasesEndpoint, server.ProjectCasesHandler, "GET")
	server.AddHandler(settings.ServerProjectCasesSelectEndpoint, server.ProjectCasesSelectHandler, "POST", "GET")
	server.AddHandler(settings.ServerProjectCasesInsertEndpoint, server.ProjectCasesInsertHandler, "POST", "GET")

	server.AddHandler(settings.ServerCaseTagsSelectEndpoint, server.CaseTagsSelectHandler, "POST", "GET")
	server.AddHandler(settings.ServerCaseTagsInsertEndpoint, server.CaseTagsInsertHandler, "POST", "GET")
	server.AddHandler(settings.ServerCaseTagsDeleteEndpoint, server.CaseTagsDeleteHandler, "POST", "GET")

	server.AddHandler(settings.ServerProjectPlansEndpoint, server.ProjectPlansHandler, "GET")
	server.AddHandler(settings.ServerProjectPlansSelectEndpoint, server.ProjectPlansSelectHandler, "POST", "GET")
	server.AddHandler(settings.ServerProjectPlansInsertEndpoint, server.ProjectPlansInsertHandler, "POST", "GET")

	server.AddHandler(settings.ServerPlanTagsSelectEndpoint, server.PlanTagsSelectHandler, "POST", "GET")
	server.AddHandler(settings.ServerPlanTagsInsertEndpoint, server.PlanTagsInsertHandler, "POST", "GET")
	server.AddHandler(settings.ServerPlanTagsDeleteEndpoint, server.PlanTagsDeleteHandler, "POST", "GET")

	server.AddHandler(settings.ServerProjectSettingsEndpoint, server.ProjectSettingsHandler, "GET")
	server.AddHandler(settings.ServerProjectRenameEndpoint, server.ProjectRenameHandler, "POST", "GET")

	server.AddHandler(settings.ServerProjectCollaboratorsEndpoint, server.ProjectCollaboratorsHandler, "POST", "GET")
	server.AddHandler(settings.ServerProjectCollaboratorsAddEndpoint, server.ProjectCollaboratorsAddHandler, "POST", "GET")
	server.AddHandler(settings.ServerProjectCollaboratorsDeleteEndpoint, server.ProjectCollaboratorsDeleteHandler, "POST", "GET")
	server.AddHandler(settings.ServerProjectCollaboratorsUpdateEndpoint, server.ProjectCollaboratorsUpdateHandler, "POST", "GET")

	server.AddHandler(settings.ServerProjectStatisticsEndpoint, server.ProjectStatisticsHandler, "GET")
	server.AddHandler(settings.ServerProjectStatisticsSelectEndpoint, server.ProjectStatisticsSelectHandler, "POST", "GET")

	server.AddHandler(settings.ServerProjectCaseEndpoint, server.ProjectCaseHandler, "GET")
	server.AddHandler(settings.ServerProjectCaseSelectEndpoint, server.ProjectCaseSelectHandler, "POST", "GET")
	server.AddHandler(settings.ServerProjectCaseUpdateEndpoint, server.ProjectCaseUpdateHandler, "POST", "GET")

	server.AddHandler(settings.ServerProjectPlanEndpoint, server.ProjectPlanHandler, "GET")
	server.AddHandler(settings.ServerProjectPlanSelectEndpoint, server.ProjectPlanSelectHandler, "POST", "GET")
	server.AddHandler(settings.ServerProjectPlanUpdateEndpoint, server.ProjectPlanUpdateHandler, "POST", "GET")
	server.AddHandler(settings.ServerProjectPlanCaseAppendEndpoint, server.ProjectPlanCaseAppendHandler, "POST", "GET")

	server.AddHandler(settings.ServerGetProjectUserRoleEndpoint, server.GetProjectUserRoleHandler, "POST", "GET")

	server.AddHandler(settings.ServerProjectPlanRunEndpoint, server.ProjectPlanRunHandler, "GET")
	server.AddHandler(settings.ServerProjectPlanRunGetEndpoint, server.ProjectPlanRunSelectHandler, "POST", "GET")
	server.AddHandler(settings.ServerProjectPlanRunInsertEndpoint, server.ProjectPlanRunInsertHandler, "POST", "GET")

	server.AddHandler(settings.ServerProjectDeleteEndpoint, server.ProjectsDeleteHandler, "POST", "GET")
	server.AddHandler(settings.ServerProjectCaseDeleteEndpoint, server.ProjectCasesDeleteHandler, "POST", "GET")
	server.AddHandler(settings.ServerProjectPlanDeleteEndpoint, server.ProjectPlansDeleteHandler, "POST", "GET")

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
