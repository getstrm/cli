package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/pflag"
	"pace/pace/pkg/bootstrap"
	"pace/pace/pkg/common"
	"pace/pace/pkg/util"
	"path"
	"strings"
)

const (
	apiHostFlag      = "api-host"
	generateDocsFlag = "generate-docs"
)

func main() {
	flags := RootCmd.Flags()
	flags.Bool(generateDocsFlag, false, "generate docs")
	err := flags.MarkHidden(generateDocsFlag)

	if err != nil {
		return
	}

	err = RootCmd.Execute()
	if err != nil {
		common.CliExit(err)
	}

	const fmTemplate = `---
title: "%s"
hide_title: true
---
`

	linkHandler := func(name string) string {
		return "docs/04-reference/01-cli-reference/" + strings.Replace(name, "_", "/", -1)
	}

	filePrepender := func(filename string) string {
		pathArray := strings.Split(filename, "/")
		filename = pathArray[len(pathArray)-1]
		pathArray = strings.Split(filename, "_")
		name := pathArray[len(pathArray)-1]
		base := strings.TrimSuffix(name, path.Ext(name))
		return fmt.Sprintf(fmTemplate, strings.Replace(base, "_", " ", -1))
	}

	if util.GetBoolAndErr(flags, generateDocsFlag) {
		err := doc.GenMarkdownTreeCustom(RootCmd, "./generated_docs", filePrepender, linkHandler)
		common.CliExit(err)
	}
}

var RootCmd = &cobra.Command{
	Use:               common.RootCommandName,
	Short:             fmt.Sprintf("Pace CLI %s", common.Version),
	PersistentPreRunE: rootCmdPreRun,
	DisableAutoGenTag: true,
}

func rootCmdPreRun(cmd *cobra.Command, args []string) error {
	util.CreateConfigDirAndFileIfNotExists()
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
	persistentFlags.String(apiHostFlag, "localhost", "api host")
	persistentFlags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort, common.OutputFormatTable, fmt.Sprintf("output format [%v]", common.OutputFormatFlagAllowedValuesText))

	err := RootCmd.RegisterFlagCompletionFunc(common.OutputFormatFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.OutputFormatFlagAllowedValues, cobra.ShellCompDirectiveNoFileComp
	})

	common.CliExit(err)
	bootstrap.SetupVerbs(RootCmd)
}
