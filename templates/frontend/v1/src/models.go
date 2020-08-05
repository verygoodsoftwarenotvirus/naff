package frontendsrc

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func typeToJSType(x string) string {
	typeMap := map[string]string{
		"string":   "string",
		"*string":  "string",
		"bool":     "boolean",
		"*bool":    "boolean",
		"int":      "number",
		"*int":     "number",
		"int8":     "number",
		"*int8":    "number",
		"int16":    "number",
		"*int16":   "number",
		"int32":    "number",
		"*int32":   "number",
		"int64":    "number",
		"*int64":   "number",
		"uint":     "number",
		"*uint":    "number",
		"uint8":    "number",
		"*uint8":   "number",
		"uint16":   "number",
		"*uint16":  "number",
		"uint32":   "number",
		"*uint32":  "number",
		"uint64":   "number",
		"*uint64":  "number",
		"float32":  "number",
		"*float32": "number",
		"float64":  "number",
		"*float64": "number",
	}

	return typeMap[x]
}

func typeToDefaultJSValue(x string) string {
	typeMap := map[string]string{
		"string":   `""`,
		"*string":  `""`,
		"bool":     "false",
		"*bool":    "false",
		"int":      "0",
		"*int":     "0",
		"int8":     "0",
		"*int8":    "0",
		"int16":    "0",
		"*int16":   "0",
		"int32":    "0",
		"*int32":   "0",
		"int64":    "0",
		"*int64":   "0",
		"uint":     "0",
		"*uint":    "0",
		"uint8":    "0",
		"*uint8":   "0",
		"uint16":   "0",
		"*uint16":  "0",
		"uint32":   "0",
		"*uint32":  "0",
		"uint64":   "0",
		"*uint64":  "0",
		"float32":  "0",
		"*float32": "0",
		"float64":  "0",
		"*float64": "0",
	}

	return typeMap[x]
}

func typeToFakerValue(x string) string {
	typeMap := map[string]string{
		"string":   "faker.random.word()",
		"*string":  "faker.random.word()",
		"bool":     "faker.random.boolean()",
		"*bool":    "faker.random.boolean()",
		"int":      "faker.random.number()",
		"*int":     "faker.random.number()",
		"int8":     "faker.random.number()",
		"*int8":    "faker.random.number()",
		"int16":    "faker.random.number()",
		"*int16":   "faker.random.number()",
		"int32":    "faker.random.number()",
		"*int32":   "faker.random.number()",
		"int64":    "faker.random.number()",
		"*int64":   "faker.random.number()",
		"uint":     "faker.random.number()",
		"*uint":    "faker.random.number()",
		"uint8":    "faker.random.number()",
		"*uint8":   "faker.random.number()",
		"uint16":   "faker.random.number()",
		"*uint16":  "faker.random.number()",
		"uint32":   "faker.random.number()",
		"*uint32":  "faker.random.number()",
		"uint64":   "faker.random.number()",
		"*uint64":  "faker.random.number()",
		"float32":  "faker.random.number()",
		"*float32": "faker.random.number()",
		"float64":  "faker.random.number()",
		"*float64": "faker.random.number()",
	}

	return typeMap[x]
}

func buildSomethingFrontendModels(typ models.DataType) func() string {
	sn := typ.Name.Singular()
	abbr := typ.Name.LowercaseAbbreviation()

	output := fmt.Sprintf(`import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class %s {
  id: number;
`, sn)

	for _, field := range typ.Fields {
		var nullableModifier string
		if field.Pointer {
			nullableModifier = "?"
		}
		output += fmt.Sprintf(
			"  %s%s: %s;\n",
			field.Name.UnexportedVarName(),
			nullableModifier,
			typeToJSType(field.Type),
		)
	}

	output += `  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
`

	for _, field := range typ.Fields {
		output += fmt.Sprintf(
			"    %s = %s;\n",
			field.Name.UnexportedVarName(),
			typeToDefaultJSValue(field.Type),
		)
	}

	output += fmt.Sprintf(`    this.createdOn = 0;
  }

static areEqual = function(
  %s1: %s,
  %s2: %s,
): boolean {
    return (
      %s1.id === %s2.id &&
`, abbr, sn, abbr, sn, abbr, abbr)

	for _, field := range typ.Fields {
		output += fmt.Sprintf(
			"      %s1.%s === %s2.%s &&\n",
			abbr,
			field.Name.UnexportedVarName(),
			abbr,
			field.Name.UnexportedVarName(),
		)
	}

	output += fmt.Sprintf(`      %s1.archivedOn === %s2.archivedOn
    );
  }
}

export const fakeValidIngredientFactory = Factory.Sync.makeFactory<ValidIngredient> ({
`, abbr, abbr)

	for _, field := range typ.Fields {
		output += fmt.Sprintf(
			"  %s: Factory.Sync.each(() =>  %s),\n",
			field.Name.UnexportedVarName(),
			typeToFakerValue(field.Type),
		)
	}

	output += `  ...defaultFactories,
});
`

	return func() string {
		return output
	}
}
