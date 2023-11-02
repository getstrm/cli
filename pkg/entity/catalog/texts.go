package catalog

import "pace/pace/pkg/util"

var listCatalogsDocs = util.LongDocs(`
Shows all the catalogs that have been configured on this Pace instance.

Catalogs are connected via configuration, and only read upon startup of the Pace server.
`)

var listCatalogsExample = util.LongDocs(`
!pace list catalogs --output table
 ID                       TYPE

 COLLIBRA-testdrive   COLLIBRA
 datahub-on-dev        DATAHUB

!pace list catalogs
catalogs:
- id: COLLIBRA-testdrive
  type: COLLIBRA
- id: datahub-on-dev
  type: DATAHUB
- id: odd
  type: ODD
`)
