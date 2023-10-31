package processingplatform

import (
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/entities/v1alpha"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/processing_platforms/v1alpha"
	"fmt"
	"github.com/elliotchance/orderedmap/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/samber/lo"
	"pace/pace/pkg/common"
)

var printer common.Printer

func availablePrinters() orderedmap.OrderedMap[string, common.Printer] {
	printers := common.DefaultPrinters.Copy()
	printers.Set(common.OutputFormatTable, listTablePrinter{})
	printers.Set(common.OutputFormatPlain, listPlainPrinter{})
	return *printers
}

type listTablePrinter struct{}
type listPlainPrinter struct{}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*ListProcessingPlatformsResponse)
	printTable(listResponse.ProcessingPlatforms)
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*ListProcessingPlatformsResponse)
	for _, platform := range listResponse.ProcessingPlatforms {

		fmt.Println(platform.Id, platform.PlatformType)
	}
}

func printTable(platforms []*DataPolicy_ProcessingPlatform) {
	rows := lo.Map(platforms, func(platform *DataPolicy_ProcessingPlatform, _ int) table.Row {
		return table.Row{
			platform.Id,
			platform.PlatformType,
		}
	})
	headers := table.Row{
		"ID",
		"Type",
	}
	common.RenderTable(headers, rows)
}
