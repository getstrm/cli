package lineage

import (
	"github.com/elliotchance/orderedmap/v2"
	"pace/pace/pkg/common"
)

var printer common.Printer

type listTablePrinter struct{}
type listPlainPrinter struct{}

// listPrinters
// printers that can handle the output of the list command
func listPrinters() orderedmap.OrderedMap[string, common.Printer] {
	printers := common.StandardPrinters.Copy()
	return *printers
}
