package processingplatform

import "pace/pace/pkg/util"

var listExample = util.DedentTrim(`
pace list processing-platforms
 ID                                 TYPE

 databricks-pim@getstrm.com   DATABRICKS
 snowflake-demo                SNOWFLAKE
 bigquery-dev                   BIGQUERY
`)
