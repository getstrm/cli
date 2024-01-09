package main

import (
	"errors"
	"fmt"
	"github.com/lithammer/dedent"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/pflag"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"os"
	"pace/pace/pkg/bootstrap"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/metrics"
	. "pace/pace/pkg/util"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	apiHostFlag = "api-host"
)

var commandPath string
var commandError error

var telemetryUploaded chan bool

func main() {
	// There's buffering in the channel because otherwise sending to the channel from the
	// main thread will block the whole application
	telemetryUploaded = make(chan bool, 1)
	var rootCmd = &cobra.Command{
		Use:               RootCommandName,
		Short:             fmt.Sprintf("PACE CLI %s", common.Version),
		PersistentPreRunE: rootCmdPreRun,
		DisableAutoGenTag: true,
		SilenceUsage:      true,
		SilenceErrors:     true,
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "generate-docs",
		Short: "D",
		RunE: func(_ *cobra.Command, _ []string) error {
			return doc.GenMarkdownTree(rootCmd, "./generated_docs")
		},
		Hidden: true,
	})

	cliExit(setup(rootCmd))
	commandError = rootCmd.Execute()
	metrics.CollectTelemetry(commandPath, commandError)
	// wait for the telemetry to have been sent but no more than 2 seconds
	select {
	case <-time.After(2 * time.Second):
	case <-telemetryUploaded:
	}
	cliExit(commandError)
}

func rootCmdPreRun(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	commandPath = cmd.CommandPath()
	if err := createConfigDirAndFileIfNotExists(); err != nil {
		return err
	}

	if err := setLastSeenTimestamp(cmd); err != nil {
		return err
	}

	if err := bootstrap.InitializeConfig(cmd); err != nil {
		return err
	}

	telemetryUploadIntervalSeconds, _ := flags.GetInt64(metrics.TelemetryUploadIntervalSeconds)
	metrics.UploadTelemetry(telemetryUploadIntervalSeconds, telemetryUploaded)

	log.Infoln(fmt.Sprintf("Executing command: %v", cmd.CommandPath()))
	flags.Visit(func(flag *pflag.Flag) {
		log.Infoln(fmt.Sprintf("flag %v=%v", flag.Name, flag.Value))
	})
	common.ApiHost, _ = cmd.Flags().GetString(apiHostFlag)
	return bootstrap.SetupServiceClients()
}

func setup(rootCmd *cobra.Command) error {
	logFile, err := common.LogFileName()
	if err != nil {
		return err
	}
	log.Traceln(fmt.Sprintf("Log file can be found at %v", logFile))
	persistentFlags := rootCmd.PersistentFlags()
	persistentFlags.String(apiHostFlag, "localhost:50051", "api host")
	persistentFlags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort,
		common.StandardPrinters.Keys()[0],
		fmt.Sprintf("output format [%v]", strings.Join(common.StandardPrinters.Keys(), ", ")))
	persistentFlags.Int64(metrics.TelemetryUploadIntervalSeconds, 3600, "Upload usage statistics every so often. Use -1 to disable")

	err = rootCmd.RegisterFlagCompletionFunc(common.OutputFormatFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.StandardPrinters.Keys(), cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}
	bootstrap.SetupVerbs(rootCmd)
	return nil
}

func createConfigDirAndFileIfNotExists() error {
	configPath, err := common.ConfigPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
		return err
	}
	configFilepath := path.Join(configPath, common.DefaultConfigFilename+common.DefaultConfigFileSuffix)
	_, err = os.Stat(configFilepath)
	if errors.Is(err, syscall.ENOENT) {
		err = os.WriteFile(
			configFilepath,
			common.DefaultConfigFileContents,
			0644,
		)
	}
	return err
}

func setLastSeenTimestamp(cmd *cobra.Command) error {
	if cmd.Name() != "__complete" {
		configPath, err := common.ConfigPath()
		if err != nil {
			return err
		}
		lastSeenCommandFilepath := path.Join(configPath, common.DefaultLastSeenFilename)
		now := time.Now()
		if content, err := os.ReadFile(lastSeenCommandFilepath); err != nil {
			return updateLastSeen(lastSeenCommandFilepath, now)
		} else {
			ts, err := strconv.ParseInt(strings.Trim(string(content), "\n"), 10, 64)
			if err != nil {
				os.Remove(lastSeenCommandFilepath)
				return updateLastSeen(lastSeenCommandFilepath, now)
			}
			lastSeen := time.Unix(ts, 0)
			if lastSeen.Unix() < now.Add(-24*time.Hour).Unix() {
				return updateLastSeen(lastSeenCommandFilepath, now)
			}
		}
	}

	return nil
}

func updateLastSeen(lastSeenCommandFilepath string, now time.Time) error {
	err := os.WriteFile(
		lastSeenCommandFilepath,
		[]byte(fmt.Sprintf("%d", now.Unix())),
		0644,
	)
	if err != nil {
		return err
	}
	printWelcomeMessage()
	return nil
}

func printWelcomeMessage() {
	res, _ := http.Get("https://cli.getstrm.com/motd")
	if res.StatusCode != 200 {
		log.Warnln("Could not fetch MOTD message, non-200 response")
		return
	}

	resBody, _ := io.ReadAll(res.Body)

	fmt.Println(string(resBody))
}

func cliExit(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.WithFields(log.Fields{"file": file, "line": line}).Error(err)

		st, ok := status.FromError(err)

		if ok {
			var additionalDetails string
			if len(st.Details()) > 0 {
				details := st.Details()[0]
				yamlBytes := ProtoMessageToYaml(details.(proto.Message))
				additionalDetails = yamlBytes.String()
			} else {
				additionalDetails = ""
			}
			formattedMessage := fmt.Sprintf(dedentAndTrimMultiline(`
						Error code = %s
						Details = %s

						%s`), (*st).Code(), (*st).Message(), additionalDetails)

			_, _ = fmt.Fprintln(os.Stderr, formattedMessage)
		} else {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}

		os.Exit(1)
	}
}

func dedentAndTrimMultiline(s string) string {
	return strings.TrimLeft(dedent.Dedent(s), "\n")
}
