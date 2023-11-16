package datapolicy

import . "pace/pace/pkg/util"

var upsertLongDocs = LongDocs(`
Upserts (inserts or updates) a data policy into PACE AND
applies it to the target platform.

The file to upsert is checked for validity, a transformation is generated
for the processing platform, and then applied on it.

When updating an existing policy, the latest existing version should be set
in the metadata. When creating a new policy, no version needs to be specified.
`)

var upsertExample = LongDocs(`
!pace upsert data-policy examples/sample_data/bigquery-cdc.json
data_policy:
  id: fb76958d-63a9-4b5e-bf36-fdc4d7ab807f
  metadata:
    title: stream-machine-development.dynamic_views.cdc_diabetes
  platform:
    id: bigquery-dev
    platform_type: BIGQUERY
  rule_sets:
  - field_transforms:
    - attribute:
        path_components:
        - HighChol
        type: integer
      transforms:
      - fixed:
          value: "****"
...
`)

var getLongDoc = LongDocs(`
retrieves a DataPolicy from PACE.

A blueprint policy is a policy that can be retrieved from a data catalog or a
processing platform with 0 or more rule sets. This means we use the table information in the platform to
build the °source° part of a data policy. We must either provide a platform or a catalog
id to make the call succeed.

Without a °--processing-platform° or a °--catalog° it just means we interact with the PACE
database and retrieve succesfully applied data policies.
`)

var getExample = LongDocs(`
# get a blueprint policy without rulesets from Catalog Collibra
!pace get data-policy --catalog COLLIBRA-testdrive \
	--database 99379294-6e87-4e26-9f09-21c6bf86d415 \
	--schema 342f676c-341e-4229-b3c2-3e71f9ed0fcd \
	6e978083-bb8f-459d-a48b-c9a50289b327
data_policy:
  metadata:
    title: employee_yearly_income
    description: Google BigQuery
  source:
    attributes:
      - path_components:
          - employee_salary
        type: bigint
      - path_components:
          - employee_id
        type: varchar
	...

# get a blueprint policy without rulesets from Processing Platform BigQuery
!pace get data-policy \
	--processing-platform bigquery-dev \
	stream-machine-development.dynamic_view_poc.gddemo
dataPolicy:
  metadata:
    createTime: '2023-10-04T09:04:56.246Z'
    description: ''
    title: stream-machine-development.dynamic_view_poc.gddemo
    updateTime: '2023-10-04T09:04:56.246Z'
  platform:
    id: bigquery-dev
    platformType: BIGQUERY
  source:
    attributes:
    - pathComponents:
      - transactionId
      type: INTEGER
    - pathComponents:
      - userId
      type: INTEGER


# get a complete datapolicy from the PACE database
!pace get data-policy --processing-platform bigquery-dev \
	stream-machine-development.dynamic_views.cdc_diabetes

id: stream-machine-development.dynamic_views.cdc_diabetes
metadata:
  create_time: "2023-11-02T12:51:23.108319730Z"
  description: ""
  title: stream-machine-development.dynamic_views.cdc_diabetes
  update_time: "2023-11-02T12:51:23.108319730Z"
  version: 1
platform:
  id: bigquery-dev
  platform_type: BIGQUERY
rule_sets:
- field_transforms:
  - field:
      name_parts:
      - HighChol
      type: integer
    transforms:
    - fixed:
        value: blabla
  target:
`)

var listExample = LongDocs(`
!pace list data-policies --output table
 PLATFORM       SOURCE                                                  TAGS

 bigquery-dev   stream-machine-development.dynamic_views.cdc_diabetes
`)

var listLongDoc = LongDocs(`
lists all the active policies defined and applied by PACE.

These will always include at least one rule set.
`)
