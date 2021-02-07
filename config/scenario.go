package config

type Scenario int

const (
	NORMAL Scenario = iota
	STALE_VIEWS
	STALE_STATES
	BYZANTINE_PRIM
	STALE_REQUESTS
	NON_SS
)

var TestCase = NORMAL

func InitializeScenario(scenario Scenario) {
	TestCase = scenario
}
