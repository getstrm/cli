package database

import "pace/pace/pkg/util"

var listDatabasesLongDocs = util.LongDocs(`
Lists all databases in a certain catalog. Some catalogs (like Collibra) are hierarchical, while others
are just a flat list of tables. In that case we might  just return a single 'Database' for the whole catalog.
`)
var listDatabasesExample = util.LongDocs(`
!pace list databases --catalog COLLIBRA-testdrive  --output table
 ID                                     NAME   TYPE

 8665f375-e08a-4810-add6-7af490f748ad          Snowflake
 99379294-6e87-4e26-9f09-21c6bf86d415          CData JDBC Driver for Google BigQuery 2021
 b6e043a7-88f1-42ee-8e81-0fdc1c96f471          Snowflake

`)
