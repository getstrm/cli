package processingplatform

import "pace/pace/pkg/util"

var listExample = util.LongDocs(`
!pace list processing-platforms --output table
 ID                                 TYPE

 databricks-pim@getstrm.com   DATABRICKS
 snowflake-demo                SNOWFLAKE
 bigquery-dev                   BIGQUERY
`)
