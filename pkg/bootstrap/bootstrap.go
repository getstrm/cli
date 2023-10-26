package bootstrap

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"pace/pace/pkg/cmd"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/catalog"
	"pace/pace/pkg/entity/database"
	"pace/pace/pkg/entity/datapolicy"
	"pace/pace/pkg/entity/group"
	"pace/pace/pkg/entity/processingplatform"
	"pace/pace/pkg/entity/schema"
	"pace/pace/pkg/entity/table"
	"strings"
)

const (
	cliVersionHeader = "pace-cli-version"
)

/*
*
these are the top level commands, i.e. the verbs.

Each verb sits in its own package, and will have subcommands for all the entity types
in Pace.
*/
func SetupVerbs(rootCmd *cobra.Command) {
	rootCmd.AddCommand(cmd.ListCmd)
	rootCmd.AddCommand(cmd.GetCmd)
}

func SetupServiceClients(accessToken *string) {
	clientConnection, ctx := SetupGrpc(common.ApiHost)

	processingplatform.SetupClient(clientConnection, ctx)
	catalog.SetupClient(clientConnection, ctx)
	table.SetupClient(clientConnection, ctx)
	group.SetupClient(clientConnection, ctx)
	schema.SetupClient(clientConnection, ctx)
	database.SetupClient(clientConnection, ctx)
	datapolicy.SetupClient(clientConnection, ctx)
}

func InitializeConfig(cmd *cobra.Command) error {
	viperConfig := viper.New()

	// Set the base name of the config file, without the file extension.
	viperConfig.SetConfigName(common.DefaultConfigFilename)

	// Set as many paths as you like where viper should look for the
	// config file.
	viperConfig.AddConfigPath(common.ConfigPath())

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := viperConfig.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable STRM_NUMBER. This helps
	// avoid conflicts.
	viperConfig.SetEnvPrefix(common.EnvPrefix)

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	viperConfig.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(cmd, viperConfig)

	log.Infoln(fmt.Sprintf("Viper configuration: %v", viperConfig.AllSettings()))

	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STRM_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			err := v.BindEnv(f.Name, fmt.Sprintf("%s_%s", common.EnvPrefix, envVarSuffix))
			common.CliExit(err)
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			common.CliExit(err)
		}
	})
}

func SetupGrpc(host string) (*grpc.ClientConn, context.Context) {

	var err error
	var creds grpc.DialOption

	creds = grpc.WithTransportCredentials(insecure.NewCredentials())

	clientConnection, err := grpc.Dial(host, creds, grpc.WithUnaryInterceptor(clientInterceptor))
	common.CliExit(err)

	var mdMap = map[string]string{cliVersionHeader: common.Version}

	return clientConnection, metadata.NewOutgoingContext(context.Background(), metadata.New(mdMap))
}

func clientInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {

	var header metadata.MD
	opts = append(opts, grpc.Header(&header))
	err := invoker(ctx, method, req, reply, cc, opts...)

	return err
}
