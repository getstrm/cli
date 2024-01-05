package plugin

import (
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/plugins/v1alpha"
	"fmt"
	"github.com/elliotchance/orderedmap/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/samber/lo"
	"pace/pace/pkg/common"
	"strings"
)

var printer common.Printer

func listPrinters() orderedmap.OrderedMap[string, common.Printer] {
	printers := common.StandardPrinters.Copy()
	printers.Set(common.OutputFormatTable, listTablePrinter{})
	printers.Set(common.OutputFormatPlain, listPlainPrinter{})
	return *printers
}

type listTablePrinter struct{}
type listPlainPrinter struct{}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*ListPluginsResponse)
	rows := lo.Map(listResponse.Plugins, func(plugin *Plugin, _ int) table.Row {
		actions := lo.Map(plugin.Actions, func(action *Action, _ int) string {
			return action.Type.String()
		})
		return table.Row{
			plugin.Id,
			strings.Join(actions, ", "),
			plugin.Implementation,
		}
	})
	headers := table.Row{
		"ID", "Actions", "Implementation",
	}
	common.RenderTable(headers, rows)
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*ListPluginsResponse)

	for _, plugin := range listResponse.Plugins {
		actions := lo.Map(plugin.Actions, func(action *Action, _ int) string {
			return action.Type.String()
		})
		fmt.Printf("%s (%s)\n", plugin.Id, strings.Join(actions, ", "))
	}
}
