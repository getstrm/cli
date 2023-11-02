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
	"pace/pace/pkg/util"
	"path"
	"path/filepath"
	"strings"
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
			util.CliExit(err)
		},
		Hidden: true,
	})
	util.CliExit(RootCmd.Execute())
}

var RootCmd = &cobra.Command{
	Use:               util.RootCommandName,
	Short:             fmt.Sprintf("Pace CLI %s", common.Version),
	PersistentPreRunE: rootCmdPreRun,
	DisableAutoGenTag: true,
}

func rootCmdPreRun(cmd *cobra.Command, _ []string) error {
	CreateConfigDirAndFileIfNotExists()
	err := bootstrap.InitializeConfig(cmd)
	log.Infoln(fmt.Sprintf("Executing command: %v", cmd.CommandPath()))
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		log.Infoln(fmt.Sprintf("flag %v=%v", flag.Name, flag.Value))
	})
	common.ApiHost = util.GetStringAndErr(cmd.Flags(), apiHostFlag)
	bootstrap.SetupServiceClients(nil)
	return err
}

func init() {
	logFile := common.LogFileName()
	log.Traceln(fmt.Sprintf("Log file can be found at %v", logFile))
	persistentFlags := RootCmd.PersistentFlags()
	persistentFlags.String(apiHostFlag, "localhost:50051", "api host")
	persistentFlags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort,
		common.DefaultPrinters.Keys()[0],
		fmt.Sprintf("output format [%v]", strings.Join(common.DefaultPrinters.Keys(), ", ")))

	err := RootCmd.RegisterFlagCompletionFunc(common.OutputFormatFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.DefaultPrinters.Keys(), cobra.ShellCompDirectiveNoFileComp
	})

	util.CliExit(err)
	bootstrap.SetupVerbs(RootCmd)
}

func CreateConfigDirAndFileIfNotExists() {
	err := os.MkdirAll(filepath.Dir(common.ConfigPath()), 0700)
	util.CliExit(err)

	configFilepath := path.Join(common.ConfigPath(), common.DefaultConfigFilename+common.DefaultConfigFileSuffix)

	if _, _ = os.Stat(configFilepath); os.IsNotExist(err) {
		writeFileError := os.WriteFile(
			configFilepath,
			common.DefaultConfigFileContents,
			0644,
		)

		util.CliExit(writeFileError)
	}
}
