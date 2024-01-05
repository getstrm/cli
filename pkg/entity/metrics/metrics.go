package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"pace/pace/pkg/common"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	TelemetryUploadIntervalSeconds = "telemetry-upload-interval-seconds"
)

// Telemetry represents all metrics that are collected about the getSTRM CLI.
// Various metrics are collected about how the CLI is used, in order to improve it.
// The metrics are collected in a way that does not identify the user, nor anything the CLI interacts with.
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
	if commandPath == "" {
		return
	}
	telemetry := readTelemetryFileContents()
	updateTelemetry(commandPath, err, telemetry)
	_ = storeTelemetry(telemetry)
}

/*
UploadTelemetry possibly uploads telemetry contents to Google Cloud function.
the interval is defines the number of seconds that have to have passed
since the last time telemetry was updated. A value of -1 will NOT send
any data to the cloud function.

`ch` is a channel that wait for true in the main thread to indicate the function
call has succeeded (or at most 2 seconds).
*/
func UploadTelemetry(telemetryUploadIntervalSeconds int64, ch chan bool) {
	if telemetryUploadIntervalSeconds < 0 { // opt out
		ch <- true
		return
	}
	telemetry := readTelemetryFileContents()
	now := time.Now().Unix()
	if (now - getTimestampFileContents()) < telemetryUploadIntervalSeconds {
		ch <- true
		return
	}
	go func() {
		updateFilestampContents(now)
		sendTelemetry(telemetry)
		ch <- true
	}()
}

/*
sendTelemetry uploads a Telemetry instance to the public cloud function.
*/
func sendTelemetry(telemetry Telemetry) {
	marshalled, _ := json.Marshal(telemetry)
	req, err := http.NewRequest("POST", "https://cli.getstrm.com/telemetry", bytes.NewReader(marshalled))
	if err != nil {
		log.Errorf("Error creating telemetry request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil || response.StatusCode != 200 {
		log.Errorf("Error sending telemetry: %v", response.Status)
		if response.Body != nil {
			responseBody, _ := io.ReadAll(response.Body)

			log.Errorf("Cloud function response:%s\n", responseBody)
		}
	} else {
		responseBody, _ := io.ReadAll(response.Body)
		if strings.TrimSpace(string(responseBody)) != "OK" {
			log.Errorf("Telemetry sent, received unexpected '%s'\n", responseBody)
		}
	}
}

/*
updateTelemetry updates the telemetry argument based on the command path
*/
func updateTelemetry(commandPath string, err error, telemetry Telemetry) {
	telemetry.CliVersion = common.Version

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

/*
determineExitCode determines process exit code based on the error.

This will become somewhat smarter in the future :-)
*/
func determineExitCode(err error) uint32 {
	if err != nil {
		return 1
	}
	return 0
}

func getTimestampFileContents() int64 {
	configPath, _ := common.ConfigPath()
	content, err := os.ReadFile(path.Join(configPath, common.TelemetryTimestampFileName))
	if err != nil {
		return 0
	}
	ts, err := strconv.ParseInt(strings.TrimSpace(string(content)), 10, 64)
	if err != nil {
		return 0
	}
	return ts
}

func updateFilestampContents(ts int64) {
	configPath, _ := common.ConfigPath()
	_ = os.WriteFile(
		path.Join(configPath, common.TelemetryTimestampFileName),
		[]byte(fmt.Sprintf("%d", ts)),
		0644,
	)
}

func storeTelemetry(telemetry Telemetry) error {
	configPath, err := common.ConfigPath()
	if err != nil {
		return err
	}
	telemetryFile := path.Join(configPath, common.DefaultTelemetryFilename)
	telemetryYaml, err := yaml.Marshal(telemetry)
	if err != nil {
		log.Errorf("Error marshalling telemetry: %v", err)
		return err
	}

	return os.WriteFile(telemetryFile, telemetryYaml, 0644)
}

func readTelemetryFileContents() Telemetry {
	configPath, err := common.ConfigPath()
	defaultTelemetry := Telemetry{
		CliVersion:   common.Version,
		OsVersion:    runtime.GOOS,
		Id:           uuid.New().String(),
		MetricPoints: make(map[string]map[uint32]Metric),
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
