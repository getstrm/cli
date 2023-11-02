package processingplatform

import . "pace/pace/pkg/util"

var listLongDocs = LongDocs(`
list all the processing platforms that are connected to PACE.
`)

var listExample = LongDocs(`
!pace list processing-platforms --output table
 ID                                 TYPE

 databricks-pim@getstrm.com   DATABRICKS
 snowflake-demo                SNOWFLAKE
 bigquery-dev                   BIGQUERY
`)
