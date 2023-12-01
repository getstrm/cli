package plugin

import . "pace/pace/pkg/util"

var invokeLongDocs = LongDocs(`
Invoke a plugin with the provided payload (JSON or YAML).
The payload file is checked for validity. The result is plugin-dependent.
`)

var invokeExample = LongDocs(`
!pace invoke plugin openai-data-policy-generator --payload example.yaml
`)
