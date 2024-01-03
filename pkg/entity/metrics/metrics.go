package metrics

import (
	"fmt"
	"github.com/google/uuid"
)

// Telemetry represents all metrics that are collected about the getSTRM CLI.
// Various metrics are collected about how the CLI is used, in order to improve it.
// The metrics are collected in a way that does not identify the user, nor anything the CLI interacts with.
// TODO add opt out mechanism.
type Telemetry struct {
	metricPoints     []Metric  `yaml:"metric_points"`
	cliVersion       string    `yaml:"cli_version"`
	operationVersion string    `yaml:"operation_version"`
	id               uuid.UUID `yaml:"id"`
}

type Metric struct {
	commands        []string `yaml:"commands"`
	cumulativeCount uint32   `yaml:"cumulative_count"`
	exitCode        uint32   `yaml:"exit_code"`
}

func CollectTelemetry(commandPath string, err error) {
	fmt.Println("Collecting telemetry")
	fmt.Println(fmt.Sprintf("Command path: %s", commandPath))
	if err != nil {
		fmt.Println(fmt.Sprintf("Command error: %s", err))
	}
}
