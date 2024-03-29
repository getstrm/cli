package datapolicy

import (
	api "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/data_policies/v1alpha"
	entities "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/entities/v1alpha"
	"fmt"
	"github.com/elliotchance/orderedmap/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
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

func lineagePrinters() orderedmap.OrderedMap[string, common.Printer] {
	// We want to use the plain printer by default for the evaluate command,
	// so put it first in the map, then add the standard printers.
	printers := orderedmap.NewOrderedMap[string, common.Printer]()
	printers.Set(common.OutputFormatTable, listLineageTablePrinter{})
	printers.Set(common.OutputFormatSimpleYaml, listLineageSimpleYamlPrinter{})
	lo.ForEach(common.StandardPrinters.Keys(), func(key string, _ int) {
		printer, _ := common.StandardPrinters.Get(key)
		printers.Set(key, printer)
	})
	return *printers
}

func transpilePrinters() orderedmap.OrderedMap[string, common.Printer] {
	printers := orderedmap.NewOrderedMap[string, common.Printer]()
	printers.Set(common.OutputFormatPlain, transpilePlainPrinter{})
	return *printers
}

type listTablePrinter struct{}
type listPlainPrinter struct{}
type evaluateTablePrinter struct{}
type transpilePlainPrinter struct{}
type listLineageSimpleYamlPrinter struct{}
type listLineageTablePrinter struct{}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*api.ListDataPoliciesResponse)
	common.RenderTable(table.Row{
		"Platform",
		"Source",
		"Tags",
	}, lo.Map(listResponse.DataPolicies, func(policy *entities.DataPolicy, _ int) table.Row {
		platform := policy.Source.Ref.GetPlatform()

		if platform == nil {
			return nil
		}

		return table.Row{
			platform.Id,
			*policy.Source.Ref.IntegrationFqn,
			strings.Join(policy.Metadata.Tags, ","),
		}
	}))
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*api.ListDataPoliciesResponse)

	for _, policy := range listResponse.DataPolicies {
		platform := policy.Source.Ref.GetPlatform()

		if platform == nil {
			continue
		}

		fmt.Println(
			platform.Id,
			policy.Source.Ref,
			strings.Join(policy.Metadata.Tags, ","),
		)
	}
}

func (p evaluateTablePrinter) Print(data interface{}) {
	evaluateResponse, _ := (data).(*api.EvaluateDataPolicyResponse)
	lo.ForEach(evaluateResponse.GetRuleSetResults(), func(result *api.EvaluateDataPolicyResponse_RuleSetResult, _ int) {
		printRuleSetResult(result)
	})
}

func printRuleSetResult(ruleSetResult *api.EvaluateDataPolicyResponse_RuleSetResult) {
	fmt.Printf("Results for rule set with target: %s\n", *ruleSetResult.Target.Ref.IntegrationFqn)
	lo.ForEach(ruleSetResult.EvaluationResults, func(result *api.EvaluateDataPolicyResponse_RuleSetResult_EvaluationResult, _ int) {
		principal := result.Principal
		if principal == nil {
			fmt.Print("All other principals\n\n")
		} else {
			common.ProtoMessageYamlPrinter{}.Print(principal)
		}
		printCsvAsTable(result.GetCsvEvaluation().Csv)
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

func (p listLineageTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*api.ScanLineageResponse)

	fmt.Println("Lineage information for connected processing platforms (✗ = not managed by PACE, ✓ = managed by PACE)")

	common.RenderTable(table.Row{
		"Fully Qualified Name",
		"Platform Id",
		"Upstream Fqns",
		"Downstream Fqns",
	}, lo.Map(listResponse.LineageSummaries, func(s *entities.LineageSummary, _ int) table.Row {
		platform := s.ResourceRef.GetPlatform()

		if platform == nil {
			return nil
		}

		return table.Row{
			*s.ResourceRef.IntegrationFqn,
			platform.Id,
			lineageAsString(s.Upstream),
			lineageAsString(s.Downstream),
		}
	}))
}

func lineageAsString(lineage []*entities.Lineage) string {
	return strings.Join(
		lo.Map(lineage, func(l *entities.Lineage, _ int) string {
			var checkmark string

			if l.ManagedByPace {
				checkmark = "✓"
			} else {
				checkmark = "✗"
			}
			return fmt.Sprintf("%s (%s)", *l.ResourceRef.IntegrationFqn, checkmark)
		}),
		"\n",
	)
}

func yamlScalarMap(args ...any) []*yaml.Node {
	return lo.Map(args, func(a any, _ int) *yaml.Node {
		return &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: fmt.Sprintf("%v", a),
		}
	})
}

func (p listLineageSimpleYamlPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*api.ScanLineageResponse)
	deps := func(lineage []*entities.Lineage) []*yaml.Node {
		return lo.Map(lineage, func(l *entities.Lineage, _ int) *yaml.Node {
			return &yaml.Node{
				Kind:    yaml.MappingNode,
				Content: yamlScalarMap("fqn", *l.ResourceRef.IntegrationFqn, "managed_by_pace", l.ManagedByPace),
			}
		})
	}

	nodes := lo.Map(listResponse.LineageSummaries, func(s *entities.LineageSummary, _ int) *yaml.Node {
		platform := s.ResourceRef.GetPlatform()

		if platform == nil {
			return nil
		}

		return &yaml.Node{
			Kind:  yaml.MappingNode,
			Value: FqnFlag,
			Content: []*yaml.Node{
				{
					Kind:  yaml.ScalarNode,
					Value: "fqn",
				},
				{
					Kind:  yaml.ScalarNode,
					Value: *s.ResourceRef.IntegrationFqn,
				},
				{
					Kind:  yaml.ScalarNode,
					Value: "platform_id",
				},
				{
					Kind:  yaml.ScalarNode,
					Value: platform.Id,
				},
				{
					Kind:        yaml.ScalarNode,
					HeadComment: "Upstream lineage",
					Value:       "upstream",
				},
				{
					Kind:    yaml.SequenceNode,
					Content: deps(s.Upstream),
				},
				{
					Kind:        yaml.ScalarNode,
					HeadComment: "Downstream lineage",
					Value:       "downstream",
				},
				{
					Kind:    yaml.SequenceNode,
					Content: deps(s.Downstream),
				},
			},
		}
	})
	tree := &yaml.Node{
		Kind:        yaml.DocumentNode,
		HeadComment: "Lineage information PACE data-policies.",
		Content: []*yaml.Node{
			{
				Kind:    yaml.SequenceNode,
				Content: nodes,
			},
		},
	}
	yamlBytes, _ := yaml.Marshal(tree)
	fmt.Println(string(yamlBytes))
}

func (p transpilePlainPrinter) Print(data interface{}) {
	transpileResponse, _ := (data).(*api.TranspileDataPolicyResponse)
	fmt.Println(transpileResponse.GetSql())
}
