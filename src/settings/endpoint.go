package settings

const (
	DatabaseInsertEndpoint = "/db/insert"
	DatabaseSelectEndpoint = "/db/select"
	DatabaseUpdateEndpoint = "/db/update"
	DatabaseDeleteEndpoint = "/db/delete"
)

const (
	ServerStatusEndpoint = "/status"

	ServerAuthEndpoint         = "/auth"
	ServerAuthRegisterEndpoint = "/auth/register"
	ServerAuthTokenEndpoint    = "/auth/token"

	ServerProjectsEndpoint       = "/{token}/projects"
	ServerProjectsSelectEndpoint = "/{token}/projects/get"
	ServerProjectsInsertEndpoint = "/{token}/projects/insert"

	ServerProjectTagsSelectEndpoint       = "/{token}/project/{id:[0-9]+}/tags/get"
	ServerProjectTagsInsertEndpoint       = "/{token}/project/{id:[0-9]+}/tags/insert"
	ServerProjectTagsDeleteEndpoint       = "/{token}/project/{id:[0-9]+}/tags/delete"
	ServerProjectStatisticsEndpoint       = "/{token}/project/{id:[0-9]+}/statistics"
	ServerProjectStatisticsSelectEndpoint = "/{token}/project/{id:[0-9]+}/statistics/get"
	ServerProjectSettingsEndpoint         = "/{token}/project/{id:[0-9]+}/settings"
	ServerProjectRenameEndpoint           = "/{token}/project/{id:[0-9]+}/rename"

	ServerProjectCollaboratorsEndpoint       = "/{token}/project/{id:[0-9]+}/collaborators"
	ServerProjectCollaboratorsAddEndpoint    = "/{token}/project/{id:[0-9]+}/collaborators/add"
	ServerProjectCollaboratorsDeleteEndpoint = "/{token}/project/{id:[0-9]+}/collaborators/delete"

	ServerProjectCasesEndpoint       = "/{token}/project/{id:[0-9]+}/cases"
	ServerProjectCasesSelectEndpoint = "/{token}/project/{id:[0-9]+}/cases/get"
	ServerProjectCasesInsertEndpoint = "/{token}/project/{id:[0-9]+}/cases/insert"
	ServerProjectCaseEndpoint        = "/{token}/project/{id:[0-9]+}/case/{caseId:[0-9]+}"
	ServerProjectCaseSelectEndpoint  = "/{token}/project/{id:[0-9]+}/case/{caseId:[0-9]+}/get"
	ServerProjectCaseUpdateEndpoint  = "/{token}/project/{id:[0-9]+}/case/{caseId:[0-9]+}/update"

	ServerProjectPlansEndpoint          = "/{token}/project/{id:[0-9]+}/plans"
	ServerProjectPlansSelectEndpoint    = "/{token}/project/{id:[0-9]+}/plans/get"
	ServerProjectPlansInsertEndpoint    = "/{token}/project/{id:[0-9]+}/plans/insert"
	ServerProjectPlanEndpoint           = "/{token}/project/{id:[0-9]+}/plan/{planId:[0-9]+}"
	ServerProjectPlanSelectEndpoint     = "/{token}/project/{id:[0-9]+}/plan/{planId:[0-9]+}/get"
	ServerProjectPlanUpdateEndpoint     = "/{token}/project/{id:[0-9]+}/plan/{planId:[0-9]+}/update"
	ServerProjectPlanCaseAppendEndpoint = "/{token}/project/{id:[0-9]+}/plan/{planId:[0-9]+}/case/append"

	ServerCaseTagsSelectEndpoint = "/{token}/project/{id:[0-9]+}/case/{caseId:[0-9]+}/tags/get"
	ServerCaseTagsInsertEndpoint = "/{token}/project/{id:[0-9]+}/case/{caseId:[0-9]+}/tags/insert"
	ServerCaseTagsDeleteEndpoint = "/{token}/project/{id:[0-9]+}/case/{caseId:[0-9]+}/tags/delete"

	ServerPlanTagsSelectEndpoint = "/{token}/project/{id:[0-9]+}/plan/{planId:[0-9]+}/tags/get"
	ServerPlanTagsInsertEndpoint = "/{token}/project/{id:[0-9]+}/plan/{planId:[0-9]+}/tags/insert"
	ServerPlanTagsDeleteEndpoint = "/{token}/project/{id:[0-9]+}/plan/{planId:[0-9]+}/tags/delete"
)
