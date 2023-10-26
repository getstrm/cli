package datapolicy

import "pace/pace/pkg/util"

var getHelp = util.LongDocs(`
retrieves a DataPolicy from Pace.

°--bare° means we retrieve a policy without RuleSets from a Catalog or
Processing Platform. This means we use the table information in the platform to
build the °source° part of a data policy. We must either provide a platform or a catalog
id to make the call succeed.

Without °--bare° it just means we interact with the Pace database and retrieve succesfully applied
data policies.
`)

var getExample = util.LongDocs(`
# get a bare policy without rulesets from Catalog Collibra
pace get datapolicy --bare --catalog COLLIBRA-testdrive \
	--database 99379294-6e87-4e26-9f09-21c6bf86d415 \
	--schema 342f676c-341e-4229-b3c2-3e71f9ed0fcd 
	6e978083-bb8f-459d-a48b-c9a50289b327 | json-case -y
data_policy:
  info:
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

# get a bare policy without rulesets from Processing Platform BigQuery
pace get datapolicy --bare --processing-platform bigquery-dev stream-machine-development.dynamic_view_poc.gddemo | json2yaml
dataPolicy:
  info:
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


# get a complete datapolicy from the Pace database
pace get datapolicy 414c8334-08e4-4655-979a-32f1c8951817 | json2yaml
dataPolicy:
  id: 414c8334-08e4-4655-979a-32f1c8951817
  info:
    createTime: '2023-10-24T13:58:05.140442789Z'
    updateTime: '2023-10-24T13:58:05.140442789Z'
  platform:
    id: snowflake-demo
    platformType: SNOWFLAKE
  ruleSets:
  - fieldTransforms:
    - attribute:
        pathComponents:
        - HIGHCHOL
      transforms:
      - fixed:
          value: bla bla
    target:
      fullname: POC.CDC_DIABETES_VIEW
  source:
    attributes:
    - pathComponents:
      - HIGHBP
`)
