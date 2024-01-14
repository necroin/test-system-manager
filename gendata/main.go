package main

import (
	"fmt"
	"tsm/src/config"
	"tsm/src/db/dbi"
	"tsm/src/db/sql"
)

const (
	configPath = "config.yml"
)

func main() {
	fmt.Println("Read config.")
	config, err := config.Load(configPath)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connect to database.")
	db, err := sql.New(config)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Generating Project.")
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
			panic(response.Error)
		}
	}

	fmt.Println("Generating TSM_TestCase.")
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
						Value: fmt.Sprintf("'[Project %d] Test Case Name %d'", projectId, testCaseId),
					},
				},
			})

			if response.Error != nil {
				panic(response.Error)
			}
		}
	}

	fmt.Println("Generating TSM_TestPlan & TSM_TestPlanTestCase.")
	for projectId := 1; projectId <= 20; projectId++ {
		for testPlanId := 1; testPlanId <= projectId; testPlanId++ {
			response := db.InsertRequest(&dbi.Request{
				Table: "TSM_TestPlan",
				Fields: []dbi.Field{
					{
						Name:  "Id",
						Value: fmt.Sprintf("%d", testPlanId),
					},
					{
						Name:  "ProjectId",
						Value: fmt.Sprintf("%d", projectId),
					},
					{
						Name:  "Name",
						Value: fmt.Sprintf("'[Project %d] Test Plan Name %d'", projectId, testPlanId),
					},
				},
			})

			if response.Error != nil {
				panic(response.Error)
			}

			for testCaseId := 1; testCaseId <= testPlanId; testCaseId++ {
				response := db.InsertRequest(&dbi.Request{
					Table: "TSM_TestPlanTestCase",
					Fields: []dbi.Field{
						{
							Name:  "ProjectId",
							Value: fmt.Sprintf("%d", projectId),
						},
						{
							Name:  "TestPlanId",
							Value: fmt.Sprintf("%d", testPlanId),
						},
						{
							Name:  "TestCaseId",
							Value: fmt.Sprintf("%d", testCaseId),
						},
						{
							Name:  "Position",
							Value: fmt.Sprintf("%d", testCaseId),
						},
					},
				})

				if response.Error != nil {
					panic(response.Error)
				}
			}
		}
	}

	fmt.Println("Generating TSM_Tags for projects.")
	for projectId := 1; projectId <= 20; projectId++ {
		for tagId := 1; tagId <= projectId; tagId++ {
			response := db.InsertRequest(&dbi.Request{
				Table: "TSM_Tags",
				Fields: []dbi.Field{
					{
						Name:  "ObjectId",
						Value: fmt.Sprintf("'%d'", projectId),
					},
					{
						Name:  "ObjectType",
						Value: "'Project'",
					},
					{
						Name:  "Name",
						Value: fmt.Sprintf("'Generated Tag %d'", tagId),
					},
				},
			})

			if response.Error != nil {
				panic(response.Error)
			}
		}
	}

	fmt.Println("Generating TSM_Tags for test cases.")
	for projectId := 1; projectId <= 20; projectId++ {
		for testCaseId := 1; testCaseId <= projectId; testCaseId++ {
			for tagId := 1; tagId <= testCaseId; tagId++ {
				response := db.InsertRequest(&dbi.Request{
					Table: "TSM_Tags",
					Fields: []dbi.Field{
						{
							Name:  "ObjectId",
							Value: fmt.Sprintf("'%d;%d'", projectId, testCaseId),
						},
						{
							Name:  "ObjectType",
							Value: "'TestCase'",
						},
						{
							Name:  "Name",
							Value: fmt.Sprintf("'Generated Tag %d'", tagId),
						},
					},
				})

				if response.Error != nil {
					panic(response.Error)
				}
			}
		}
	}

	fmt.Println("Generating TSM_Tags for test plans.")
	for projectId := 1; projectId <= 20; projectId++ {
		for testPlanId := 1; testPlanId <= projectId; testPlanId++ {
			for tagId := 1; tagId <= testPlanId; tagId++ {
				response := db.InsertRequest(&dbi.Request{
					Table: "TSM_Tags",
					Fields: []dbi.Field{
						{
							Name:  "ObjectId",
							Value: fmt.Sprintf("'%d;%d'", projectId, testPlanId),
						},
						{
							Name:  "ObjectType",
							Value: "'TestPlan'",
						},
						{
							Name:  "Name",
							Value: fmt.Sprintf("'Generated Tag %d'", tagId),
						},
					},
				})

				if response.Error != nil {
					panic(response.Error)
				}
			}
		}
	}

	fmt.Println("Finished.")
}
