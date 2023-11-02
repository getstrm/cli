package group

import . "pace/pace/pkg/util"

var listExample = LongDocs(`
!pace list groups --processing-platform bigquery-dev --output table
 NAME

 marketing
 foo-bar
 fraud-detection
`)

var listLongDocs = LongDocs(`
list the groups that exist in a processing platform.

These groups are needed in the rule sets to determine group membership of the
entity executing the query on the view in the rule set.
`)
