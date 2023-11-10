package globaltransform

import . "pace/pace/pkg/util"

var upsertLongDocs = LongDocs(`
Insert or update a global transform. The 'ref' field in the global policy is the identifier of the global policy.
If it is the same as that of an already existing one, the existing one will be replaced.
`)

var upsertExample = Example(`
!pace upsert global-transform global-tag-transform.yaml

	transform:
	  description: This is a global transform ...
	  ref: PII_EMAIL
	  tag_transform:
		tag_content: PII_EMAIL
		transforms:
		- principals:
		  - group: FRAUD_AND_RISK
		  regexp:
			regexp: ^.*(@.*)$
			replacement: '****\1'
		- nullify: {}
`)

var getLongDoc = LongDocs(`
Returns a global transform from Pace, by transform reference and transform type.
`)

var getExample = Example(`
!pace get global-transform PII_EMAIL

	description: ...
	ref: PII_EMAIL
	tag_transform:
	  tag_content: PII_EMAIL
	  transforms:
	  - principals:
		- group: FRAUD_AND_RISK
		regexp:
		  regexp: ^.*(@.*)$
		  replacement: '****\1'
	  - nullify: {}
`)

var listLongDoc = LongDocs(`
Lists all the global transforms stored in Pace.
`)

var listExample = Example(`
!pace list global-transforms

	global_transforms:
	- description: This ...
	  ref: pii-email
	  tag_transform:
		tag_content: pii-email
	...
`)

var deleteExample = Example(`
!pace delete global-transform PII_EMAIL

	deleted_count: 1
`)
