package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/pflag"
	"os"
	"pace/pace/pkg/bootstrap"
	"pace/pace/pkg/common"
	. "pace/pace/pkg/util"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	apiHostFlag = "api-host"
)

func main() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "generate-docs",
		Short: "D",
		Run: func(_ *cobra.Command, _ []string) {
			err := doc.GenMarkdownTree(RootCmd, "./generated_docs")
			CliExit(err)
		},
		Hidden: true,
	})
	CliExit(RootCmd.Execute())
}

var RootCmd = &cobra.Command{
	Use:               RootCommandName,
	Short:             fmt.Sprintf("PACE CLI %s", common.Version),
	PersistentPreRunE: rootCmdPreRun,
	DisableAutoGenTag: true,
}

func rootCmdPreRun(cmd *cobra.Command, _ []string) error {
	CreateConfigDirAndFileIfNotExists()
	CreateLastSeenCommand()
	err := bootstrap.InitializeConfig(cmd)
	log.Infoln(fmt.Sprintf("Executing command: %v", cmd.CommandPath()))
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		log.Infoln(fmt.Sprintf("flag %v=%v", flag.Name, flag.Value))
	})
	common.ApiHost = GetStringAndErr(cmd.Flags(), apiHostFlag)
	bootstrap.SetupServiceClients()
	return err
}

func init() {
	logFile := common.LogFileName()
	log.Traceln(fmt.Sprintf("Log file can be found at %v", logFile))
	persistentFlags := RootCmd.PersistentFlags()
	persistentFlags.String(apiHostFlag, "localhost:50051", "api host")
	persistentFlags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort,
		common.StandardPrinters.Keys()[0],
		fmt.Sprintf("output format [%v]", strings.Join(common.StandardPrinters.Keys(), ", ")))

	CliExit(RootCmd.RegisterFlagCompletionFunc(common.OutputFormatFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.StandardPrinters.Keys(), cobra.ShellCompDirectiveNoFileComp
	}))
	bootstrap.SetupVerbs(RootCmd)
}

func CreateConfigDirAndFileIfNotExists() {
	CliExit(os.MkdirAll(filepath.Dir(common.ConfigPath()), 0700))
	configFilepath := path.Join(common.ConfigPath(), common.DefaultConfigFilename+common.DefaultConfigFileSuffix)
	if _, _ = os.Stat(configFilepath); os.IsNotExist(os.MkdirAll(filepath.Dir(common.ConfigPath()), 0700)) {
		CliExit(os.WriteFile(
			configFilepath,
			common.DefaultConfigFileContents,
			0644,
		))
	}
}

func CreateLastSeenCommand() {
	lastSeenCommandFilepath := path.Join(common.ConfigPath(), common.DefaultLastSeenFilename)
	now := time.Now()
	if content, err := os.ReadFile(lastSeenCommandFilepath); err != nil {
		updateLastSeen(lastSeenCommandFilepath, now)
	} else {
		ts, _ := strconv.ParseInt(strings.Trim(string(content), "\n"), 10, 64)
		lastSeen := time.Unix(ts, 0)
		if lastSeen.Unix() < now.Add(-24*time.Hour).Unix() {
			updateLastSeen(lastSeenCommandFilepath, now)
		}
	}
}

func updateLastSeen(lastSeenCommandFilepath string, now time.Time) {
	CliExit(os.WriteFile(
		lastSeenCommandFilepath,
		[]byte(fmt.Sprintf("%d", now.Unix())),
		0644,
	))
	printWelcomeMessage()
}

func printWelcomeMessage() {
	asciiArt := `
---------------------------------------------------------------------
                        ____   _    ____ _____                         
                       |  _ \ / \  / ___| ____|
                       | |_) / _ \| |   |  _|  
                       |  __/ ___ \ |___| |___ 
                       |_| /_/   \_\____|_____|
---------------------------------------------------------------------
Hey there, cool you're using PACE!
We'd love to learn what your use case is or if you have any feedback.
Join our Slack: https://getstrm.com/slack
---------------------------------------------------------------------`
	fmt.Println(asciiArt)
}
