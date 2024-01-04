package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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
	StatsUploadInterval = "stats-interval"
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
	telemetry := readTelemetry()
	updateTelemetry(commandPath, err, telemetry)
	_ = storeTelemetry(telemetry)
}

func UploadTelemetry(statsInterval int64, ch chan bool) {
	if statsInterval < 0 {
		ch <- true
		return
	}
	telemetry := readTelemetry()
	now := time.Now().Unix()
	if (now - getUploadedTimestamp()) < statsInterval {
		ch <- true
		return
	}
	go func() {
		updateUploadedTimestamp(now)
		sendTelemetry(telemetry)
		ch <- true
	}()
}

func sendTelemetry(telemetry Telemetry) {
	marshalled, err := json.Marshal(telemetry)
	req, err := http.NewRequest("POST", "https://cli.getstrm.com/telemetry", bytes.NewReader(marshalled))
	if err != nil {
		fmt.Println("Error creating telemetry request")
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil || response.StatusCode != 200 {
		fmt.Println("Error sending telemetry", response.Status)
		if response.Body != nil {
			b, _ := io.ReadAll(response.Body)
			fmt.Println("Cloud function response:\n", string(b))
		}
	} else {
		b, _ := io.ReadAll(response.Body)
		if strings.TrimSpace(string(b)) != "OK" {
			fmt.Println("Telemetry sent, received unexpected '", string(b), "'")

		}
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

func getUploadedTimestamp() int64 {
	configPath, _ := common.ConfigPath()
	lastSeenCommandFilepath := path.Join(configPath, common.TelemetryTimestampFileName)
	content, err := os.ReadFile(lastSeenCommandFilepath)
	if err != nil {
		return 0
	}
	ts, _ := strconv.ParseInt(strings.Trim(string(content), "\n"), 10, 64)
	return ts
}

func updateUploadedTimestamp(now int64) {
	configPath, _ := common.ConfigPath()
	lastSeenCommandFilepath := path.Join(configPath, common.TelemetryTimestampFileName)
	os.WriteFile(
		lastSeenCommandFilepath,
		[]byte(fmt.Sprintf("%d", now)),
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
		fmt.Println("Error marshalling telemetry")
		return err
	}

	return os.WriteFile(telemetryFile, telemetryYaml, 0644)
}

func readTelemetry() Telemetry {
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
