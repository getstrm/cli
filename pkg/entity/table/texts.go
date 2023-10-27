package table

import "pace/pace/pkg/util"

var listExample = util.LongDocs(`
!pace list tables --catalog COLLIBRA-testdrive \
	--database 99379294-6e87-4e26-9f09-21c6bf86d415 \
	--schema c0a8b864-83e7-4dd1-a71d-0c356c1ae9be

!pace list tables --processing-platform bigquery-dev
`)
