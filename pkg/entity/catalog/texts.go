package catalog

import "pace/pace/pkg/util"

var docs = util.LongDocs(`
Shows all the catalogs that have been configured on this Pace instance.

Catalogs are connected via configuration, and only read upon startup of the Pace server.
`)
var listExample = util.LongDocs(`
!pace list catalogs
 ID                       TYPE

 COLLIBRA-testdrive   COLLIBRA
 datahub-on-dev        DATAHUB

# in yaml
!pace list catalogs -o yaml
catalogs:
- id: COLLIBRA-testdrive
  type: COLLIBRA
- id: datahub-on-dev
  type: DATAHUB
- id: odd
  type: ODD
`)
