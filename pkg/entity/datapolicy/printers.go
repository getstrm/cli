package datapolicy

import (
	api "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/data_policies/v1alpha"
	entities "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/entities/v1alpha"
	"fmt"
	"github.com/elliotchance/orderedmap/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/samber/lo"
	"pace/pace/pkg/common"
	"strings"
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
	listResponse, _ := (data).(*api.ListDataPoliciesResponse)
	common.RenderTable(table.Row{
		"Platform",
		"Source",
		"Tags",
	}, lo.Map(listResponse.DataPolicies, func(policy *entities.DataPolicy, _ int) table.Row {
		return table.Row{
			policy.Platform.Id,
			policy.Source.Ref,
			strings.Join(policy.Metadata.Tags, ","),
		}
	}))
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*api.ListDataPoliciesResponse)
	for _, policy := range listResponse.DataPolicies {
		fmt.Println(
			policy.Platform.Id,
			policy.Source.Ref,
			strings.Join(policy.Metadata.Tags, ","),
		)
	}
}
