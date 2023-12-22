package database

import (
	catalogApi "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/data_catalogs/v1alpha"
	entities "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/entities/v1alpha"
	platformApi "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/processing_platforms/v1alpha"
	"errors"
	"fmt"
	"github.com/elliotchance/orderedmap/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/samber/lo"
	"pace/pace/pkg/common"
	"pace/pace/pkg/util"
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
	common.RenderTable(table.Row{
		"ID", "Name", "Type",
	}, lo.Map(toDatabases(data), func(database *entities.Database, _ int) table.Row {
		return table.Row{
			database.Id,
			database.DisplayName,
			database.Type,
		}
	}))
}

func (p listPlainPrinter) Print(data interface{}) {
	for _, database := range toDatabases(data) {
		fmt.Println(database.Id, database.DisplayName, database.Type)
	}
}

func toDatabases(data interface{}) []*entities.Database {
	var databases []*entities.Database
	if listResponse, ok := (data).(*catalogApi.ListDatabasesResponse); ok {
		databases = listResponse.Databases
	} else if listResponse, ok := (data).(*platformApi.ListDatabasesResponse); ok {
		databases = listResponse.Databases
	} else {
		util.CliExit(errors.New("could not handle server response"))
	}
	return databases
}
