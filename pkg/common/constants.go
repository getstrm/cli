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

var DefaultConfigFileContents = []byte(`# The following configuration options are reflected in the CLI's flags
# api-host: localhost:50051
# output: yaml
`)

var DefaultFileTypesCompletion = []string{"yml", "yaml", "json"}

const GetCommandName = "get"
const ListCommandName = "list"
const UpsertCommandName = "upsert"
const ApplyCommandName = "apply"
const EvaluateCommandName = "evaluate"
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

const PageSizeFlag = "page_size"
const PageSkipFlag = "skip"
