package metrics

import (
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"os"
	"pace/pace/pkg/common"
	"path"
	"runtime"
)

// Telemetry represents all metrics that are collected about the getSTRM CLI.
// Various metrics are collected about how the CLI is used, in order to improve it.
// The metrics are collected in a way that does not identify the user, nor anything the CLI interacts with.
// TODO add opt out mechanism.
type Telemetry struct {
	MetricPoints map[string]map[uint32]Metric `yaml:"metric_points"`
	CliVersion   string                       `yaml:"cli_version"`
	OsVersion    string                       `yaml:"operation_version"`
	Id           string                       `yaml:"id"`
}

type Metric struct {
	CumulativeCount uint32 `yaml:"cumulative_count"`
}

func CollectTelemetry(commandPath string, err error) {
	telemetry := readTelemetry()
	updateTelemetry(commandPath, err, telemetry)
	_ = writeTelemetry(telemetry)
}

func updateTelemetry(commandPath string, err error, telemetry Telemetry) {
	exitCode := determineExitCode(err)
	metricsByExitCode, ok := telemetry.MetricPoints[commandPath]

	if ok {
		metricForExitCode, ok := metricsByExitCode[exitCode]
		if ok {
			telemetry.MetricPoints[commandPath][exitCode] = Metric{
				CumulativeCount: metricForExitCode.CumulativeCount + 1,
			}
		} else {
			metricsByExitCode[exitCode] = Metric{
				CumulativeCount: 1,
			}
		}
	} else {
		telemetry.MetricPoints[commandPath] = map[uint32]Metric{
			exitCode: {
				CumulativeCount: 1,
			},
		}
	}
}

func determineExitCode(err error) uint32 {
	if err != nil {
		return 1
	}
	return 0
}
func writeTelemetry(telemetry Telemetry) error {
	configPath, err := common.ConfigPath()
	if err != nil {
		return err
	}
	telemetryFile := path.Join(configPath, common.DefaultTelemetryFilename)

	telemetryYaml, err := yaml.Marshal(telemetry)
	if err != nil {
		fmt.Println("Error marshalling telemetry")
	}

	err = os.WriteFile(telemetryFile, telemetryYaml, 0644)
	if err != nil {
		return err
	}
	return nil
}

func readTelemetry() Telemetry {
	configPath, err := common.ConfigPath()
	defaultTelemetry := Telemetry{
		CliVersion: common.Version,
		OsVersion:  runtime.GOOS,
		Id:         uuid.New().String(),
	}
	if err != nil {
		return defaultTelemetry
	}
	telemetryFile := path.Join(configPath, common.DefaultTelemetryFilename)

	telemetryYaml, err := os.ReadFile(telemetryFile)
	if err != nil {
		return defaultTelemetry
	}

	var telemetry Telemetry
	err = yaml.Unmarshal(telemetryYaml, &telemetry)
	if err != nil {
		return defaultTelemetry
	}
	return telemetry
}
