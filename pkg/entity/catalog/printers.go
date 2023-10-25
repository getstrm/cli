package catalog

import (
	data_policiesv1alpha "buf.build/gen/go/getstrm/daps/protocolbuffers/go/getstrm/api/data_policies/v1alpha"
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), common.OutputFormatFlag)

	p := availablePrinters()[outputFormat+command.Parent().Name()]

	if p == nil {
		common.CliExit(errors.New(fmt.Sprintf("Output format '%v' is not supported. Allowed values: %v", outputFormat, common.OutputFormatFlagAllowedValuesText)))
	}

	return p
}

func availablePrinters() map[string]util.Printer {
	return util.MergePrinterMaps(
		util.DefaultPrinters,
		map[string]util.Printer{
			common.OutputFormatTable + common.ListCommandName: listTablePrinter{},
			common.OutputFormatPlain + common.ListCommandName: listPlainPrinter{},
		},
	)
}

type listTablePrinter struct{}
type listPlainPrinter struct{}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*data_policiesv1alpha.ListCatalogsResponse)
	printTable(listResponse.Catalogs)
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*data_policiesv1alpha.ListCatalogsResponse)
	printPlain(listResponse.Catalogs)
}

func printTable(catalogs []*data_policiesv1alpha.DataCatalog) {
	rows := make([]table.Row, 0, len(catalogs))
	for _, catalog := range catalogs {

		row := table.Row{
			catalog.Id,
			catalog.Type,
		}
		rows = append(rows, row)
	}

	headers := table.Row{
		"ID",
		"Type",
	}
	util.RenderTable(headers, rows)
}
func printPlain(catalogs []*data_policiesv1alpha.DataCatalog) {
	for _, catalog := range catalogs {
		fmt.Println(catalog.Id, catalog.Type)
	}
}
