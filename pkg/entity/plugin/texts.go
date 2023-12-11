package plugin

import . "pace/pace/pkg/util"

var invokeLongDocs = LongDocs(`
Invoke an action for a plugin with the provided payload (JSON or YAML).
The payload file is checked for validity. The result is plugin-dependent.
`)

var invokeExample = LongDocs(`
!pace invoke plugin openai GENERATE_DATA_POLICY --payload example.yaml
`)

var listLongDocs = LongDocs(`
List all available plugins.
`)

var listExample = LongDocs(`
!pace list plugins
`)
