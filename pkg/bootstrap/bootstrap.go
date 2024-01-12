package bootstrap

import (
	. "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/data_catalogs/v1alpha/data_catalogsv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/data_policies/v1alpha/data_policiesv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/global_transforms/v1alpha/global_transformsv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/plugins/v1alpha/pluginsv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/processing_platforms/v1alpha/processing_platformsv1alphagrpc"
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
	"pace/pace/pkg/entity/globaltransform"
	"pace/pace/pkg/entity/group"
	"pace/pace/pkg/entity/lineage"
	"pace/pace/pkg/entity/plugin"
	"pace/pace/pkg/entity/processingplatform"
	"pace/pace/pkg/entity/schema"
	"pace/pace/pkg/entity/table"
	"strings"
)

const (
	cliVersionHeader = "pace-cli-version"
)

/*
	these are the top level commands, i.e. the verbs.

Each verb sits in its own package, and will have subcommands for all the entity types
in PACE.
*/
func SetupVerbs(rootCmd *cobra.Command) {
	rootCmd.AddCommand(
		cmd.ListCmd,
		cmd.GetCmd,
		cmd.UpsertCmd,
		cmd.ApplyCmd,
		cmd.DeleteCmd,
		cmd.VersionCmd,
		cmd.EvaluateCmd,
		cmd.InvokeCmd,
		cmd.DisableCmd,
	)
}

func SetupServiceClients() error {
	connection, ctx, err := SetupGrpc(common.ApiHost)
	if err != nil {
		return err
	}
	dataPoliciesClient := NewDataPoliciesServiceClient(connection)
	catalogsClient := NewDataCatalogsServiceClient(connection)
	processingPlatformsClient := NewProcessingPlatformsServiceClient(connection)
	globalTransformsClient := NewGlobalTransformsServiceClient(connection)
	pluginsClient := NewPluginsServiceClient(connection)
	processingplatform.SetupClient(processingPlatformsClient, ctx)
	catalog.SetupClient(catalogsClient, ctx)
	table.SetupClient(processingPlatformsClient, catalogsClient, ctx)
	group.SetupClient(processingPlatformsClient, ctx)
	schema.SetupClient(processingPlatformsClient, catalogsClient, ctx)
	database.SetupClient(processingPlatformsClient, catalogsClient, ctx)
	datapolicy.SetupClient(dataPoliciesClient, catalogsClient, processingPlatformsClient, ctx)
	globaltransform.SetupClient(globalTransformsClient, ctx)
	plugin.SetupClient(pluginsClient, ctx)
	lineage.SetupClient(processingPlatformsClient, ctx)
	return nil
}

func InitializeConfig(cmd *cobra.Command) error {
	viperConfig := viper.New()

	// Set the base name of the config file, without the file extension.
	viperConfig.SetConfigName(common.DefaultConfigFilename)

	// Set as many paths as you like where viper should look for the
	// config file.
	configPath, err := common.ConfigPath()
	if err != nil {
		return err
	}
	viperConfig.AddConfigPath(configPath)

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
			_ = v.BindEnv(f.Name, fmt.Sprintf("%s_%s", common.EnvPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

func SetupGrpc(host string) (*grpc.ClientConn, context.Context, error) {
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())
	clientConnection, err := grpc.Dial(host, creds, grpc.WithUnaryInterceptor(clientInterceptor))
	if err != nil {
		return nil, nil, err
	}
	var mdMap = map[string]string{cliVersionHeader: common.Version}
	return clientConnection, metadata.NewOutgoingContext(context.Background(), metadata.New(mdMap)), nil
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
