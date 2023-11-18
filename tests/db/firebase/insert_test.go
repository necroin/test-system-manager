package firebase_test

import (
	"testing"
	"tsm/src/config"
	"tsm/src/db/dbi"
	"tsm/src/db/firebase"
)

const (
	localDirPath = "../../../"
	configPath   = "config.yml"
)

func TestDatabase_InsertOnePartPK(t *testing.T) {
	config, err := config.Load(localDirPath + configPath)
	if err != nil {
		t.Fatal(err)
	}
	config.Credentials = localDirPath + config.Credentials
	config.Database.Schema = localDirPath + config.Database.Schema

	db, err := firebase.New(config)
	if err != nil {
		t.Fatal(err)
	}

	response := db.InsertRequest(&dbi.Request{
		Table: "Projects",
		Fields: []dbi.Field{
			{
				Name:  "Id",
				Value: "1",
			},
			{
				Name:  "Name",
				Value: "Project Name 1",
			},
		},
	})

	if response.Error != nil {
		t.Fatal(err)
	}

	response = db.InsertRequest(&dbi.Request{
		Table: "Projects",
		Fields: []dbi.Field{
			{
				Name:  "Id",
				Value: "2",
			},
			{
				Name:  "Name",
				Value: "Project Name 2",
			},
		},
	})

	if response.Error != nil {
		t.Fatal(err)
	}
}

func TestDatabase_InsertTwoPartPK(t *testing.T) {
	config, err := config.Load(localDirPath + configPath)
	if err != nil {
		t.Fatal(err)
	}
	config.Credentials = localDirPath + config.Credentials
	config.Database.Schema = localDirPath + config.Database.Schema

	db, err := firebase.New(config)
	if err != nil {
		t.Fatal(err)
	}

	response := db.InsertRequest(&dbi.Request{
		Table: "TSM_ProjectTestCase",
		Fields: []dbi.Field{
			{
				Name:  "ProjectId",
				Value: "1",
			},
			{
				Name:  "TestCaseId",
				Value: "1",
			},
		},
	})

	if response.Error != nil {
		t.Fatal(err)
	}

	response = db.InsertRequest(&dbi.Request{
		Table: "TSM_ProjectTestCase",
		Fields: []dbi.Field{
			{
				Name:  "ProjectId",
				Value: "2",
			},
			{
				Name:  "TestCaseId",
				Value: "1",
			},
		},
	})

	if response.Error != nil {
		t.Fatal(err)
	}
}
