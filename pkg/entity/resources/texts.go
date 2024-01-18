package resources

import (
	"pace/pace/pkg/util"
)

var listDocs = util.LongDocs(`
list resources using a directory path style interface, i.e. forward slash
separated path components.

The top level path component is the identifier of one of the configured data-catalogs or
processing platforms in PACE.

`)

var listExample = util.PlainExample(`
# Show all the configured integrations (data catalogs and processing platforms)
!pace list resources
 INTEGRATION           TYPE         ID

 processing-platform   DATABRICKS   dbr-pace
 processing-platform   SNOWFLAKE    sf-pace
 processing-platform   BIGQUERY     bigquery-dev
 processing-platform   POSTGRES     paceso
 data-catalog          COLLIBRA     COLLIBRA-testdrive
 data-catalog          DATAHUB      datahub-on-dev

!pace list -o table resources bigquery-dev/stream-machine-development/batch_job_demo
 TABLE              FQN

 retail_0           stream-machine-development.batch_job_demo.retail_0
 retail_encrypted   stream-machine-development.batch_job_demo.retail_encrypted
 retail_in          stream-machine-development.batch_job_demo.retail_in
 retail_keys        stream-machine-development.batch_job_demo.retail_keys
`)
