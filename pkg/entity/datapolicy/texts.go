package datapolicy

import . "pace/pace/pkg/util"

var upsertLongDocs = LongDocs(`
Upserts (inserts or updates) a data policy into PACE AND
optionally applies it to the target platform (default false).

The file to upsert is checked for validity, a transformation is generated
for the processing platform, and then applied on it.

By default, the version does not need to be set in the metadata, PACE will
auto-increment it. If, however, PACE has been configured to not do so, then 
when updating an existing policy, the latest existing version should be set
in the metadata. When creating a new policy, no version needs to be specified.
This is the case when the property °app.data-policies.auto-increment-version°
is set to false in the PACE configuration file.
`)

var upsertExample = PlainExample(`
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

var applyLongDocs = LongDocs(`
Applies an existing data policy to the target platform.
`)

var applyExample = PlainExample(`
!pace apply data-policy public.demo --processing-platform bigquery-dev
data_policy:
  id: public.demo
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

var evaluateLongDocs = LongDocs(`
Evaluates an existing data policy by applying it to sample data provided in a csv file.
You can use this to test the correctness of your field transforms and filters.
The csv file should contain a header row with the column names, matching the fields in the data policy.
A comma should be used as the delimiter.
Currently, only standard SQL data types are supported. For platform-specific transforms, test on the platform itself.
`)

var evaluateExample = PlainExample(`
!pace evaluate data-policy public.demo --processing-platform example-platform --sample-data sample.csv
Results for rule set with target: public.demo_view
group: administrator

 TRANSACTIONID   USERID   EMAIL                      AGE   BRAND    TRANSACTIONAMOUNT 
                                                                          
 534704584       870941   acole@gmail.com            4     HP       7                 
 807835672       867943   knappjeremy@hotmail.com    49    Acer     10                
 467414030       251481   morriserin@hotmail.com     6     Acer     277               
 994186205       500392   wgolden@yahoo.com          68    Lenovo   160               
 217127008       143855   nelsondaniel@hotmail.com   28    Lenovo   263               
 142409570       567637   meganriley@gmail.com       56    Acer     296               

group: marketing

 TRANSACTIONID   USERID   EMAIL              AGE   BRAND   TRANSACTIONAMOUNT 
                                                                 
 807835672       0        ****@hotmail.com   49    Other   10                
 994186205       0        ****@yahoo.com     68    Other   160               
 217127008       0        ****@hotmail.com   28    Other   263               
 142409570       0        ****@gmail.com     56    Other   296               

All other principals

 TRANSACTIONID   USERID   EMAIL   AGE   BRAND   TRANSACTIONAMOUNT 
                                                      
 807835672       0        ****    49    Other   10                
 994186205       0        ****    68    Other   160               
 217127008       0        ****    28    Other   263               
 142409570       0        ****    56    Other   296  
`)

var getLongDoc = LongDocs(`
retrieves a DataPolicy from PACE.

A blueprint policy is a policy that can be retrieved from a data catalog or a
processing platform with 0 or more rule sets. This means we use the table information in the platform to
build the °source° part of a data policy. We must either provide a platform or a catalog
id to make the call succeed.

Without a °--processing-platform° or a °--catalog° it just means we interact with the PACE
database and retrieve successfully applied data policies.
`)

var getExample = PlainExample(`
# get a blueprint policy without rule sets from Catalog Collibra
!pace get data-policy --catalog COLLIBRA-testdrive \
	--blueprint \
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

# get a blueprint policy without rule sets from Processing Platform BigQuery
!pace get data-policy \
	--blueprint \
	--processing-platform bigquery-dev \
	--database stream-machine-development \
	--schema dynamic_view_poc \ 
	gddemo
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

# get a datapolicy from a processing platform via its fully qualified name

!pace get data-policy \	
	--processing-platform bigquery-dev    \
	--blueprint \
	--fqn=true  \
	stream-machine-development.data_lineage_demo.total_green_trips_22_21

...	

# get a complete datapolicy (with rulesets) from the PACE database
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
        value: redacted
  target:
`)

var listExample = PlainExample(`
!pace list data-policies --output table
 PLATFORM       SOURCE                                                  TAGS

 bigquery-dev   stream-machine-development.dynamic_views.cdc_diabetes
`)

var listLongDoc = LongDocs(`
lists all the active policies defined and applied by PACE.

These will always include at least one rule set.
`)
