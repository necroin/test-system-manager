package settings

const (
	ServerWaitStartRepeatCount  = 10
	ServerWaitStartSleepSeconds = 1
	ServerStatusResponse        = "OK"
)

const (
	InterfaceProjectsHTML          = "assets/interface/projects.html"
	InterfaceProjectCasesHTML      = "assets/interface/project_cases.html"
	InterfaceProjectPlansHTML      = "assets/interface/project_plans.html"
	InterfaceProjectSettingsHTML   = "assets/interface/project_settings.html"
	InterfaceProjectStatisticsHTML = "assets/interface/project_statistics.html"
	InterfaceCaseHTML              = "assets/interface/case.html"
	InterfacePlanHTML              = "assets/interface/plan.html"
	InterfacePlanRunHTML           = "assets/interface/plan_run.html"
	InterfaceStylePath             = "assets/interface/style.css"
	InterfaceScriptPath            = "assets/interface/script.js"
)

const (
	ObjectTypeProject  = "Project"
	ObjectTypeTestCase = "TestCase"
	ObjectTypeTestPlan = "TestPlan"
)

const (
	ProjectRoleCreator = "Создатель"
	ProjectRoleAnalyst = "Аналитик"
	ProjectRoleTester  = "Тестировщик"
	ProjectRoleGuest   = "Гость"
)
