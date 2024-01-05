package schema

import (
	catalogApi "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/data_catalogs/v1alpha"
	entities "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/entities/v1alpha"
	platformApi "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/processing_platforms/v1alpha"
	"fmt"
	"github.com/elliotchance/orderedmap/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/samber/lo"
	"pace/pace/pkg/common"
)

var printer common.Printer

// listPrinters
// printers that can handle the output of the list command
func listPrinters() orderedmap.OrderedMap[string, common.Printer] {
	printers := common.StandardPrinters.Copy()
	printers.Set(common.OutputFormatTable, listTablePrinter{})
	printers.Set(common.OutputFormatPlain, listPlainPrinter{})
	return *printers
}

type listTablePrinter struct{}
type listPlainPrinter struct{}

func (p listTablePrinter) Print(data interface{}) {
	rows := lo.Map(toSchemas(data), func(schema *entities.Schema, _ int) table.Row {
		return table.Row{
			schema.Id,
			schema.Name,
		}
	})
	headers := table.Row{
		"ID", "Name",
	}
	common.RenderTable(headers, rows)
}

func (p listPlainPrinter) Print(data interface{}) {
	for _, schema := range toSchemas(data) {
		fmt.Println(schema.Id, schema.Name)
	}
}

func toSchemas(data interface{}) []*entities.Schema {
	var schemas []*entities.Schema
	if listResponse, ok := (data).(*catalogApi.ListSchemasResponse); ok {
		schemas = listResponse.Schemas
	} else if listResponse, ok := (data).(*platformApi.ListSchemasResponse); ok {
		schemas = listResponse.Schemas
	} else {
		schemas = nil
	}
	return schemas
}
