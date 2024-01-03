package main

import (
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
	"time"
)

const (
	apiHostFlag = "api-host"
)

var commandPath string
var commandError error

func main() {
	cobra.OnFinalize(func() {

	})
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
			err := doc.GenMarkdownTree(rootCmd, "./generated_docs")
			return err
		},
		Hidden: true,
	})

	cliExit(setup(rootCmd))
	commandError = rootCmd.Execute()
	metrics.CollectTelemetry(commandPath, commandError)
	cliExit(commandError)
}

func rootCmdPreRun(cmd *cobra.Command, _ []string) error {
	commandPath = cmd.CommandPath()
	err := createConfigDirAndFileIfNotExists()
	if err != nil {
		return err
	}

	err = setLastSeenTimestamp(cmd)
	if err != nil {
		return err
	}

	err = bootstrap.InitializeConfig(cmd)
	if err != nil {
		return err
	}

	log.Infoln(fmt.Sprintf("Executing command: %v", cmd.CommandPath()))
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		log.Infoln(fmt.Sprintf("flag %v=%v", flag.Name, flag.Value))
	})
	v, _ := cmd.Flags().GetString(apiHostFlag)
	common.ApiHost = v
	err = bootstrap.SetupServiceClients()

	return err
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
	err = os.MkdirAll(filepath.Dir(configPath), 0700)
	if err != nil {
		return err
	}
	configFilepath := path.Join(configPath, common.DefaultConfigFilename+common.DefaultConfigFileSuffix)
	if _, _ = os.Stat(configFilepath); os.IsNotExist(os.MkdirAll(filepath.Dir(configPath), 0700)) {
		err := os.WriteFile(
			configFilepath,
			common.DefaultConfigFileContents,
			0644,
		)
		if err != nil {
			return err
		}
	}

	return nil
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
				additionalDetails = string(yamlBytes.Bytes())
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
