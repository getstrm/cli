package common

import "strings"

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

var DefaultConfigFileContents = []byte(`# The following configuration options are reflected in the CLI's flags
# api-host: localhost
`)

const GetCommandName = "get"
const ListCommandName = "list"
const CreateCommandName = "create"
const DeleteCommandName = "delete"
const UpdateCommandName = "update"
const OutputFormatJson = "json"
const OutputFormatJsonRaw = "json-raw"
const OutputFormatTable = "table"
const OutputFormatPlain = "plain"

const OutputFormatFlag = "output"
const OutputFormatFlagShort = "o"

const ProcessingPlatformFlag = "processing-platform"
const ProcessingPlatformFlagShort = "p"

const CatalogFlag = "catalog"
const CatalogFlagShort = "c"

const DatabaseFlag = "database"
const DatabaseFlagShort = "d"

const SchemaFlag = "schema"
const SchemaFlagShort = "s"

var OutputFormatFlagAllowedValues = []string{OutputFormatJson, OutputFormatJsonRaw, OutputFormatTable, OutputFormatPlain}
var OutputFormatFlagAllowedValuesText = strings.Join(OutputFormatFlagAllowedValues, ", ")
