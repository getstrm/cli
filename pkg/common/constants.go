package common

var configPath string
var logFileName string

var GitSha = "dev"
var Version = "dev"
var BuiltOn = "unknown"

// The environment variable prefix of all environment variables bound to our command line flags.
// For example, --api-host is bound to PACE_API_HOST
const EnvPrefix = "PACE"

const DefaultConfigFilename = "config"
const DefaultConfigFileSuffix = ".yaml"

const DefaultLastSeenFilename = "pace-cli-last-seen"
const DefaultTelemetryFilename = "telemetry.yaml"
const TelemetryTimestampFileName = "telemetry-last-upload-timestamp"

var DefaultConfigFileContents = []byte(`
# The following configuration options are reflected in the CLI's flags
# api-host: localhost:50051
# output: yaml
# change value below to -1 if you don't want cli usage statistics to be sent to getSTRM
# telemetry-upload-interval-seconds: 3600
`)

var DefaultFileTypesCompletion = []string{"yml", "yaml", "json"}

const GetCommandName = "get"
const ListCommandName = "list"
const UpsertCommandName = "upsert"
const ApplyCommandName = "apply"
const EvaluateCommandName = "evaluate"
const TranspileCommandName = "transpile"
const InvokeCommandName = "invoke"
const DeleteCommandName = "delete"
const DisableCommandName = "disable"
const ProcessingPlatformFlag = "processing-platform"
const ProcessingPlatformFlagShort = "p"
const ProcessingPlatformFlagUsage = `id of processing platform`

const CatalogFlag = "catalog"
const CatalogFlagShort = "c"
const CatalogFlagUsage = `id of catalog`

const DatabaseFlag = "database"
const DatabaseFlagShort = "d"
const DatabaseFlagUsage = "database in the catalog"

const SchemaFlag = "schema"
const SchemaFlagShort = "s"
const SchemaFlagUsage = "schema in database on catalog"

const BlueprintFlag = "blueprint"
const BlueprintFlagShort = "b"
const BlueprintFlagUsage = "fetch a blueprint data policy from a catalog or a processing platform"

const ApplyFlag = "apply"
const ApplyFlagShort = "a"
const ApplyFlagUsage = "apply a data policy to the target processing platform when upserting"

const PluginPayloadFlag = "payload"
const PluginPayloadFlagUsage = "path to a json or yaml file containing the payload to invoke a plugin with"

const SampleDataFlag = "sample-data"
const SampleDataUsage = "path to a csv file containing sample data to evaluate a data policy"

const PageSizeFlag = "page-size"
const PageSkipFlag = "skip"
const PageTokenFlag = "page-token"

const InlineDataPolicyFlag = "data-policy-file"
const InlineDataPolicyUsage = "path to a data policy file, must be a yaml or json representation of a data policy"
const DataPolicyIdFlag = "data-policy-id"
const DataPolicyIdUsage = "an id of an existing data policy (does not have to be applied)"

const PrincipalsToEvaluateFlag = "principals"
const PrincipalsToEvaluateUsage = "comma separated list of principals to evaluate the data policy for, if unspecified, all principals will be evaluated. For example, --principals user1,user2. If you want to evaluate the `other` / `fallback` principal, use value `null` / `other` or `fallback`, for example --principals null"
