package resources

import (
	entities "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/entities/v1alpha"
	api "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/resources/v1alpha"
	"fmt"
	"github.com/elliotchance/orderedmap/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/samber/lo"
	"pace/pace/pkg/common"
)

var printer common.Printer

type listTablePrinter struct{}
type listPlainPrinter struct{}

// listPrinters
// printers that can handle the output of the list command
func listPrinters() orderedmap.OrderedMap[string, common.Printer] {
	// We want to use the table printer by default for the list command
	// so put it first in the map, then add the standard printers.
	printers := orderedmap.NewOrderedMap[string, common.Printer]()
	printers.Set(common.OutputFormatTable, listTablePrinter{})
	lo.ForEach(common.StandardPrinters.Keys(), func(key string, _ int) {
		printer, _ := common.StandardPrinters.Get(key)
		printers.Set(key, printer)
	})
	return *printers
}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*api.ListResourcesResponse)
	if len(listResponse.Resources) == 0 {
		fmt.Println(" no data")
		return
	}
	row1, err := lo.Last(listResponse.Resources[0].ResourcePath)
	if err != nil {
		// top level
		common.RenderTable(table.Row{
			"Integration",
			"Type",
			"Id",
		}, lo.Map(listResponse.Resources, func(urn *entities.ResourceUrn, _ int) table.Row {
			what, typ, id := parse(urn)
			return table.Row{
				what,
				typ,
				id,
			}
		}))

	} else {
		common.RenderTable(table.Row{
			row1.PlatformName,
			"DisplayName",
			"Fqn",
		}, lo.Map(listResponse.Resources, func(urn *entities.ResourceUrn, _ int) table.Row {
			name, _ := lo.Last(urn.ResourcePath)
			return table.Row{
				name.Name,
				name.DisplayName,
				urn.PlatformFqn,
			}
		}))
	}
}

func parse(urn *entities.ResourceUrn) (string, string, string) {
	switch urn.Integration.(type) {
	case *entities.ResourceUrn_Catalog:
		return "data-catalog", urn.GetCatalog().Type.String(), urn.GetCatalog().Id
	case *entities.ResourceUrn_Platform:
		return "processing-platform", urn.GetPlatform().PlatformType.String(), urn.GetPlatform().Id
	}
	return "", "", ""

}
