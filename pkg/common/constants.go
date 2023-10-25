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

const SavedEntitiesDirectory = "saved-entities"

const GetCommandName = "get"
const ListCommandName = "list"
const CreateCommandName = "create"
const DeleteCommandName = "delete"
const UpdateCommandName = "update"

const RecursiveFlagName = "recursive"
const RecursiveFlagUsage = "Retrieve entities and their dependents"
const RecursiveFlagShorthand = "r"

const OutputFormatJson = "json"
const OutputFormatJsonRaw = "json-raw"
const OutputFormatTable = "table"
const OutputFormatPlain = "plain"
const OutputFormatCsv = "csv"
const OutputFormatFilepath = "path"

const OutputFormatFlag = "output"
const OutputFormatFlagShort = "o"

var OutputFormatFlagAllowedValues = []string{OutputFormatJson, OutputFormatJsonRaw, OutputFormatTable, OutputFormatPlain}
var OutputFormatFlagAllowedValuesText = strings.Join(OutputFormatFlagAllowedValues, ", ")

var UsageOutputFormatFlagAllowedValues = []string{OutputFormatCsv, OutputFormatJson, OutputFormatJsonRaw}
var UsageOutputFormatFlagAllowedValuesText = strings.Join(UsageOutputFormatFlagAllowedValues, ", ")

var ContextOutputFormatFlagAllowedValues = []string{OutputFormatJson, OutputFormatJsonRaw, OutputFormatFilepath}
var ContextOutputFormatFlagAllowedValuesText = strings.Join(ContextOutputFormatFlagAllowedValues, ", ")

var ConfigOutputFormatFlagAllowedValues = []string{OutputFormatPlain, OutputFormatJson}
var ConfigOutputFormatFlagAllowedValuesText = strings.Join(ConfigOutputFormatFlagAllowedValues, ", ")

var AccountOutputFormatFlagAllowedValues = []string{OutputFormatPlain, OutputFormatJson, OutputFormatJsonRaw}
var AccountOutputFormatFlagAllowedValuesText = strings.Join(AccountOutputFormatFlagAllowedValues, ", ")

var ProjectOutputFormatFlagAllowedValues = []string{OutputFormatPlain}
var ProjectOutputFormatFlagAllowedValuesText = strings.Join(ProjectOutputFormatFlagAllowedValues, ", ")

var MonitorFollowOutputFormatFlagAllowedValues = []string{OutputFormatPlain, OutputFormatJson, OutputFormatJsonRaw}
var MonitorFollowOutputFormatFlagAllowedValuesText = strings.Join(MonitorFollowOutputFormatFlagAllowedValues, ", ")
var MonitorOutputFormatFlagAllowedValues = []string{OutputFormatTable, OutputFormatPlain, OutputFormatJson, OutputFormatJsonRaw}
var MonitorOutputFormatFlagAllowedValuesText = strings.Join(MonitorOutputFormatFlagAllowedValues, ", ")

var LogsOutputFormatFlagAllowedValues = []string{OutputFormatPlain}
var LogsOutputFormatFlagAllowedValuesText = strings.Join(LogsOutputFormatFlagAllowedValues, ", ")
