package common

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/grpc/status"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"runtime"
	"strings"
)

var RootCommandName = "pace"

var ApiHost string

func CliExit(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.WithFields(log.Fields{"file": file, "line": line}).Error(err)

		st, ok := status.FromError(err)

		if ok {
			_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf(`Error code = %s
Details = %s`, (*st).Code(), (*st).Message()))
		} else {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}

		os.Exit(1)
	}
}

func Abort(format string, args ...interface{}) {
	if len(args) == 0 {
		CliExit(errors.New(format))
	} else {
		CliExit(errors.New(fmt.Sprintf(format, args...)))
	}
}

func GrpcRequestCompletionError(err error) ([]string, cobra.ShellCompDirective) {
	errorMessage := fmt.Sprintf("%v", err)
	log.Errorln(errorMessage)
	cobra.CompErrorln(errorMessage)

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func NoFilesEmptyCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func ConfigPath() string {
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

		CliExit(err)
	}

	return configPath
}

func LogFileName() string {
	if logFileName == "" {
		logFileName = ConfigPath() + "/" + RootCommandName + ".log"
		log.SetLevel(log.TraceLevel)
		log.SetOutput(&lumberjack.Logger{
			Filename:   LogFileName(),
			MaxSize:    1, // MB
			MaxBackups: 0,
		})
		log.Info(fmt.Sprintf("Config path is set to: %v", ConfigPath()))
	}

	return logFileName
}

func GetCatalogCoordinates(flags *pflag.FlagSet) (string, string, string) {
	catalogId, err := flags.GetString(CatalogFlag)
	databaseId, err := flags.GetString(DatabaseFlag)
	schemaId, err := flags.GetString(SchemaFlag)
	CliExit(err)
	return catalogId, databaseId, schemaId
}

func SetOutputFormats(flags *pflag.FlagSet, formats ...string) {
	outputFormatFlagAllowedValuesText := strings.Join(formats, ", ")
	flags.StringP(OutputFormatFlag, OutputFormatFlagShort, formats[0],
		fmt.Sprintf("output format [%v]", outputFormatFlagAllowedValuesText))
}
