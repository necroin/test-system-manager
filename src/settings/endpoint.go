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

	ServerGetUsernameEndpoint        = "/{token}/username"
	ServerGetProjectUserRoleEndpoint = "/{token}/project/{id:[0-9]+}/user/role"

	ServerProjectsEndpoint       = "/{token}/projects"
	ServerProjectsSelectEndpoint = "/{token}/projects/get"
	ServerProjectsInsertEndpoint = "/{token}/projects/insert"
	ServerProjectsDeleteEndpoint = "/{token}/projects/delete"

	ServerProjectCommentsSelectEndpoint = "/{token}/project/{projid:[0-9]+}/{type:[a-z]+}/{objId:[0-9]+}/comments/get"
	ServerProjectCommentsInsertEndpoint = "/{token}/project/{projid:[0-9]+}/{type:[a-z]+}/{objId:[0-9]+}/comments/insert"
	ServerProjectCommentsDeleteEndpoint = "/{token}/project/{projid:[0-9]+}/{type:[a-z]+}/{objId:[0-9]+}/comments/delete/{id:[0-9]+}"

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
	ServerProjectCollaboratorsUpdateEndpoint = "/{token}/project/{id:[0-9]+}/collaborators/update"

	ServerProjectCasesEndpoint       = "/{token}/project/{id:[0-9]+}/cases"
	ServerProjectCasesSelectEndpoint = "/{token}/project/{id:[0-9]+}/cases/get"
	ServerProjectCasesInsertEndpoint = "/{token}/project/{id:[0-9]+}/cases/insert"
	ServerProjectCasesDeleteEndpoint = "/{token}/project/{id:[0-9]+}/cases/delete"
	ServerProjectCaseEndpoint        = "/{token}/project/{id:[0-9]+}/case/{caseId:[0-9]+}"
	ServerProjectCaseSelectEndpoint  = "/{token}/project/{id:[0-9]+}/case/{caseId:[0-9]+}/get"
	ServerProjectCaseUpdateEndpoint  = "/{token}/project/{id:[0-9]+}/case/{caseId:[0-9]+}/update"

	ServerProjectPlansEndpoint          = "/{token}/project/{id:[0-9]+}/plans"
	ServerProjectPlansSelectEndpoint    = "/{token}/project/{id:[0-9]+}/plans/get"
	ServerProjectPlansInsertEndpoint    = "/{token}/project/{id:[0-9]+}/plans/insert"
	ServerProjectPlansDeleteEndpoint    = "/{token}/project/{id:[0-9]+}/plans/delete"
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

	ServerProjectPlanRunEndpoint       = "/{token}/project/{id:[0-9]+}/plan/{planId:[0-9]+}/run"
	ServerProjectPlanRunGetEndpoint    = "/{token}/project/{id:[0-9]+}/plan/{planId:[0-9]+}/run/get"
	ServerProjectPlanRunInsertEndpoint = "/{token}/project/{id:[0-9]+}/plan/{planId:[0-9]+}/run/insert"
)
