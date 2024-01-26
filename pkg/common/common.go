package common

import (
	pagingv1alpha "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/paging/v1alpha"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	. "pace/pace/pkg/util"
)

var ApiHost string

func AbortError(format string, args ...interface{}) error {
	if len(args) == 0 {
		return errors.New(format)
	} else {
		return fmt.Errorf(format, args...)
	}
}

func CobraCompletionError(err error) ([]string, cobra.ShellCompDirective) {
	errorMessage := fmt.Sprintf("%v", err)
	log.Errorln(errorMessage)
	cobra.CompErrorln(errorMessage)
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func NoFilesEmptyCompletion(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func ConfigPath() (string, error) {
	if configPath == "" {
		// if we set this environment variable, we work in a completely different configuration directory
		// Here you will define your flags and configuration settings.
		// Cobra supports persistent flags, which, if defined here,
		// will be global for your application.
		// set the default configuration path
		configPathEnvVar := EnvPrefix + "_CONFIG_PATH"
		configPathEnv := os.Getenv(configPathEnvVar)
		defaultConfigPath := "~/.config/pace"

		var err error

		if len(configPathEnv) != 0 {
			log.Debugln("Value for " + configPathEnvVar + " found in environment: " + configPathEnv)
			configPath, err = ExpandTilde(configPathEnv)
		} else {
			log.Debugln("No value for " + configPathEnvVar + " found. Falling back to default: " + defaultConfigPath)
			configPath, err = ExpandTilde(defaultConfigPath)
		}

		if err != nil {
			return "", err
		}
	}

	return configPath, nil
}

func LogFileName() (string, error) {
	if logFileName == "" {
		configPath, err := ConfigPath()

		if err != nil {
			return "", err
		}

		logFileName = configPath + "/" + RootCommandName + ".log"
		log.SetLevel(log.TraceLevel)
		log.SetOutput(&lumberjack.Logger{
			Filename:   logFileName,
			MaxSize:    1, // MB
			MaxBackups: 0,
		})
		log.Info(fmt.Sprintf("Config path is set to: %v", configPath))
	}

	return logFileName, nil
}

func GetCatalogCoordinates(flags *pflag.FlagSet) (string, string, string, error) {
	catalogId, _ := flags.GetString(CatalogFlag)
	databaseId, _ := flags.GetString(DatabaseFlag)
	schemaId, _ := flags.GetString(SchemaFlag)
	return catalogId, databaseId, schemaId, nil
}

func PageParameters(cmd *cobra.Command) *pagingv1alpha.PageParameters {
	flags := cmd.Flags()
	skip, _ := flags.GetUint32(PageSkipFlag)
	size, _ := flags.GetUint32(PageSizeFlag)
	token, _ := flags.GetString(PageTokenFlag)
	return &pagingv1alpha.PageParameters{
		Skip:      skip,
		PageSize:  size,
		PageToken: token,
	}
}
