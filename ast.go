package graphqlws

import (
	"github.com/lab259/graphql/language/ast"
)

func operationDefinitionsWithOperation(
	doc *ast.Document,
	op string,
) []*ast.OperationDefinition {
	defs := []*ast.OperationDefinition{}
	for _, node := range doc.Definitions {
		if node.GetKind() == "OperationDefinition" {
			if def, ok := node.(*ast.OperationDefinition); ok {
				if def.Operation == op {
					defs = append(defs, def)
				}
			}
		}
	}
	return defs
}

func selectionSetsForOperationDefinitions(
	defs []*ast.OperationDefinition,
) []*ast.SelectionSet {
	sets := []*ast.SelectionSet{}
	for _, def := range defs {
		if set := def.GetSelectionSet(); set != nil {
			sets = append(sets, set)
		}
	}
	return sets
}

func nameForSelectionSet(set *ast.SelectionSet) ([]string, bool) {
	if len(set.Selections) >= 1 {
		r := make([]string, len(set.Selections))
		for i, selection := range set.Selections {
			if field, ok := selection.(*ast.Field); ok {
				r[i] = field.Name.Value
			}
		}
		return r, true
	}
	return nil, false
}

func namesForSelectionSets(sets []*ast.SelectionSet) []string {
	names := make([]string, 0)
	for _, set := range sets {
		if nameList, ok := nameForSelectionSet(set); ok {
			names = append(names, nameList...)
		}
	}
	return names
}

func subscriptionFieldNamesFromDocument(doc *ast.Document) []string {
	defs := operationDefinitionsWithOperation(doc, "subscription")
	sets := selectionSetsForOperationDefinitions(defs)
	return namesForSelectionSets(sets)
}
