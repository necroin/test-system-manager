package data_test

import (
	"fmt"
	"testing"
	"tsm/src/config"
	"tsm/src/db/dbi"
	"tsm/src/db/sql"
)

const (
	localDirPath = "../../../"
	configPath   = "config.yml"
)

func TestDatabase_GenerateTestData(t *testing.T) {
	config, err := config.Load(localDirPath + configPath)
	if err != nil {
		t.Fatal(err)
	}
	config.Database.Schema = localDirPath + config.Database.Schema
	config.Database.Storage = localDirPath + config.Database.Storage

	db, err := sql.New(config)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	for i := 1; i <= 20; i++ {
		response := db.InsertRequest(&dbi.Request{
			Table: "Projects",
			Fields: []dbi.Field{
				{
					Name:  "Id",
					Value: fmt.Sprintf("%d", i),
				},
				{
					Name:  "Name",
					Value: fmt.Sprintf("'Project Name %d'", i),
				},
			},
		})

		if response.Error != nil {
			t.Fatal(response.Error)
		}
	}

	for projectId := 1; projectId <= 20; projectId++ {
		for testCaseId := 1; testCaseId <= projectId; testCaseId++ {
			response := db.InsertRequest(&dbi.Request{
				Table: "TSM_TestCase",
				Fields: []dbi.Field{
					{
						Name:  "Id",
						Value: fmt.Sprintf("%d", testCaseId),
					},
					{
						Name:  "ProjectId",
						Value: fmt.Sprintf("%d", projectId),
					},
					{
						Name:  "Name",
						Value: fmt.Sprintf("'Test Case Name %d'", testCaseId),
					},
				},
			})

			if response.Error != nil {
				t.Fatal(response.Error)
			}
		}
	}

	for projectId := 1; projectId <= 20; projectId++ {
		for testCaseId := 1; testCaseId <= projectId; testCaseId++ {
			response := db.InsertRequest(&dbi.Request{
				Table: "TSM_TestPlan",
				Fields: []dbi.Field{
					{
						Name:  "Id",
						Value: fmt.Sprintf("%d", testCaseId),
					},
					{
						Name:  "ProjectId",
						Value: fmt.Sprintf("%d", projectId),
					},
					{
						Name:  "Name",
						Value: fmt.Sprintf("'Test Plan Name %d'", testCaseId),
					},
				},
			})

			if response.Error != nil {
				t.Fatal(response.Error)
			}
		}
	}
}
