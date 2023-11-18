package settings

const (
	DatabaseInsertEndpoint = "/db/insert"
	DatabaseSelectEndpoint = "/db/select"
	DatabaseUpdateEndpoint = "/db/update"
	DatabaseDeleteEndpoint = "/db/delete"
)

const (
	ServerStatusEndpoint         = "/status"
	ServerProjectsEndpoint       = "/projects"
	ServerProjectsSelectEndpoint = "/projects/get"

	ServerProjectCasesEndpoint       = "/project/{id:[0-9]+}/cases"
	ServerProjectCasesSelectEndpoint = "/project/{id:[0-9]+}/cases/get"
	ServerProjectCaseEndpoint        = "/project/{id:[0-9]+}/case/{caseId:[0-9]+}"
	ServerProjectCaseSelectEndpoint  = "/project/{id:[0-9]+}/case/{caseId:[0-9]+}/get"
	ServerProjectCaseUpdateEndpoint  = "/project/{id:[0-9]+}/case/{caseId:[0-9]+}/update"

	ServerProjectPlansEndpoint       = "/project/{id:[0-9]+}/plans"
	ServerProjectPlansSelectEndpoint = "/project/{id:[0-9]+}/plans/get"
	ServerProjectPlanEndpoint        = "/project/{id:[0-9]+}/plan/{planId:[0-9]+}"
	ServerProjectPlanSelectEndpoint  = "/project/{id:[0-9]+}/plan/{planId:[0-9]+}/get"

	ServerProjectStatisticsEndpoint = "/project/{id:[0-9]+}/statistics"
)
