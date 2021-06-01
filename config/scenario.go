package config

var (
	Scenario string

	scenarios = map[int]string{
		0: "NORMAL",               // Normal execution
		1: "IDLE",                 // Byzantine processes remain idle (send nothing)
		2: "Sleep",                // Byzantine processes send different messages to half the servers
		3: "TimestampChange",      // Byzantine processes send wrong bytes for BC
		4: "SmallTimestamp",       // Byzantine processes send wrong bytes for BC
		5: "SmallTimestampToSome", // Byzantine processes send wrong bytes for BC
		6: "Fork",
	}
)

func InitializeScenario(s int) {
	Scenario = scenarios[s]
}
