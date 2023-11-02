package group

import (
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/processing_platforms/v1alpha"
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
	listResponse, _ := (data).(*ListGroupsResponse)
	common.RenderTable(table.Row{
		"Name",
	}, lo.Map(listResponse.Groups, func(group string, _ int) table.Row {
		return table.Row{
			group,
		}
	}))
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*ListGroupsResponse)
	for _, group := range listResponse.Groups {
		fmt.Println(group)
	}
}
