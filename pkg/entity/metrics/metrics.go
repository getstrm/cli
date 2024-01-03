package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"net/http"
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
	MetricPoints map[string]map[uint32]Metric `json:"metric_points" yaml:"metric_points"`
	CliVersion   string                       `json:"cli_version" yaml:"cli_version"`
	OsVersion    string                       `json:"operation_version" yaml:"operation_version"`
	Id           string                       `json:"id" yaml:"id"`
}

type Metric struct {
	CumulativeCount uint32 `json:"cumulative_count" yaml:"cumulative_count"`
}

func CollectTelemetry(commandPath string, err error) {
	telemetry := readTelemetry()
	updateTelemetry(commandPath, err, telemetry)
	_ = writeTelemetry(telemetry)
	sendTelemetry(telemetry)
}

func sendTelemetry(telemetry Telemetry) {
	marshalled, err := json.Marshal(telemetry)
	req, err := http.NewRequest("POST", "https://cli.getstrm.com/telemetry", bytes.NewReader(marshalled))
	if err != nil {
		fmt.Println("Error creating telemetry request")
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		fmt.Println("Error sending telemetry")
	} else {
		fmt.Println("Telemetry sent")
	}
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
