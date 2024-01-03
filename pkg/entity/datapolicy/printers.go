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

	"encoding/csv"
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

func evaluatePrinters() orderedmap.OrderedMap[string, common.Printer] {
	// We want to use the table printer by default for the evaluate command,
	// so put it first in the map, then add the standard printers.
	printers := orderedmap.NewOrderedMap[string, common.Printer]()
	printers.Set(common.OutputFormatTable, evaluateTablePrinter{})
	lo.ForEach(common.StandardPrinters.Keys(), func(key string, _ int) {
		printer, _ := common.StandardPrinters.Get(key)
		printers.Set(key, printer)
	})
	return *printers
}

type listTablePrinter struct{}
type listPlainPrinter struct{}
type evaluateTablePrinter struct{}

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

func (p evaluateTablePrinter) Print(data interface{}) {
	evaluateResponse, _ := (data).(*api.EvaluateDataPolicyResponse)
	lo.ForEach(evaluateResponse.GetFullEvaluationResult().RuleSetResults, func(result *api.EvaluateDataPolicyResponse_FullEvaluationResult_RuleSetResult, _ int) {
		printRuleSetResult(result)
	})
}

func printRuleSetResult(ruleSetResult *api.EvaluateDataPolicyResponse_FullEvaluationResult_RuleSetResult) {
	fmt.Printf("Results for rule set with target: %s\n", ruleSetResult.Target.Fullname)
	lo.ForEach(ruleSetResult.PrincipalEvaluationResults, func(result *api.EvaluateDataPolicyResponse_FullEvaluationResult_RuleSetResult_PrincipalEvaluationResult, _ int) {
		principal := result.Principal
		if principal == nil {
			fmt.Print("All other principals\n\n")
		} else {
			common.ProtoMessageYamlPrinter{}.Print(principal)
		}
		printCsvAsTable(result.Csv)
		fmt.Println()
	})
	fmt.Println()
}

func printCsvAsTable(csvString string) {
	csvRows, _ := csv.NewReader(strings.NewReader(csvString)).ReadAll()
	headers := csvRows[0]
	common.RenderTable(common.SliceToRow(headers), lo.Map(csvRows[1:], func(row []string, _ int) table.Row {
		return common.SliceToRow(row)
	}))
}
