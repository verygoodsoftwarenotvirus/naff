package naming

import (
	"testing"
)

func TestName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input              string
		singular           string
		plural             string
		snake              string
		pluralSnake        string
		camel              string
		pluralCamel        string
		kebab              string
		pluralKebab        string
		packageName        string
		abbreviation       string
		lowerAbbreviation  string
		humanReadable      string
		pluralHumanReadable string
		article            string
		withArticle        string
	}{
		{
			input:              "IssueReport",
			singular:           "IssueReport",
			plural:             "IssueReports",
			snake:              "issue_report",
			pluralSnake:        "issue_reports",
			camel:              "issueReport",
			pluralCamel:        "issueReports",
			kebab:              "issue-report",
			pluralKebab:        "issue-reports",
			packageName:        "issuereports",
			abbreviation:       "IR",
			lowerAbbreviation:  "ir",
			humanReadable:      "issue report",
			pluralHumanReadable: "issue reports",
			article:            "an",
			withArticle:        "an issue report",
		},
		{
			input:              "Widget",
			singular:           "Widget",
			plural:             "Widgets",
			snake:              "widget",
			pluralSnake:        "widgets",
			camel:              "widget",
			pluralCamel:        "widgets",
			kebab:              "widget",
			pluralKebab:        "widgets",
			packageName:        "widgets",
			abbreviation:       "W",
			lowerAbbreviation:  "w",
			humanReadable:      "widget",
			pluralHumanReadable: "widgets",
			article:            "a",
			withArticle:        "a widget",
		},
		{
			input:              "RecipeStep",
			singular:           "RecipeStep",
			plural:             "RecipeSteps",
			snake:              "recipe_step",
			pluralSnake:        "recipe_steps",
			camel:              "recipeStep",
			pluralCamel:        "recipeSteps",
			kebab:              "recipe-step",
			pluralKebab:        "recipe-steps",
			packageName:        "recipesteps",
			abbreviation:       "RS",
			lowerAbbreviation:  "rs",
			humanReadable:      "recipe step",
			pluralHumanReadable: "recipe steps",
			article:            "a",
			withArticle:        "a recipe step",
		},
		{
			input:              "Account",
			singular:           "Account",
			plural:             "Accounts",
			snake:              "account",
			pluralSnake:        "accounts",
			camel:              "account",
			pluralCamel:        "accounts",
			kebab:              "account",
			pluralKebab:        "accounts",
			packageName:        "accounts",
			abbreviation:       "A",
			lowerAbbreviation:  "a",
			humanReadable:      "account",
			pluralHumanReadable: "accounts",
			article:            "an",
			withArticle:        "an account",
		},
		{
			input:              "UserIngredientPreference",
			singular:           "UserIngredientPreference",
			plural:             "UserIngredientPreferences",
			snake:              "user_ingredient_preference",
			pluralSnake:        "user_ingredient_preferences",
			camel:              "userIngredientPreference",
			pluralCamel:        "userIngredientPreferences",
			kebab:              "user-ingredient-preference",
			pluralKebab:        "user-ingredient-preferences",
			packageName:        "useringredientpreferences",
			abbreviation:       "UIP",
			lowerAbbreviation:  "uip",
			humanReadable:      "user ingredient preference",
			pluralHumanReadable: "user ingredient preferences",
			article:            "a",
			withArticle:        "a user ingredient preference",
		},
		{
			input:              "UnitOfMeasure",
			singular:           "UnitOfMeasure",
			plural:             "UnitOfMeasures",
			snake:              "unit_of_measure",
			pluralSnake:        "unit_of_measures",
			camel:              "unitOfMeasure",
			pluralCamel:        "unitOfMeasures",
			kebab:              "unit-of-measure",
			pluralKebab:        "unit-of-measures",
			packageName:        "unitofmeasures",
			abbreviation:       "UOM",
			lowerAbbreviation:  "uom",
			humanReadable:      "unit of measure",
			pluralHumanReadable: "unit of measures",
			article:            "a",
			withArticle:        "a unit of measure",
		},
		{
			input:              "HourLog",
			singular:           "HourLog",
			plural:             "HourLogs",
			snake:              "hour_log",
			pluralSnake:        "hour_logs",
			camel:              "hourLog",
			pluralCamel:        "hourLogs",
			kebab:              "hour-log",
			pluralKebab:        "hour-logs",
			packageName:        "hourlogs",
			abbreviation:       "HL",
			lowerAbbreviation:  "hl",
			humanReadable:      "hour log",
			pluralHumanReadable: "hour logs",
			article:            "an",
			withArticle:        "an hour log",
		},
		{
			input:              "Upload",
			singular:           "Upload",
			plural:             "Uploads",
			snake:              "upload",
			pluralSnake:        "uploads",
			camel:              "upload",
			pluralCamel:        "uploads",
			kebab:              "upload",
			pluralKebab:        "uploads",
			packageName:        "uploads",
			abbreviation:       "U",
			lowerAbbreviation:  "u",
			humanReadable:      "upload",
			pluralHumanReadable: "uploads",
			article:            "an",
			withArticle:        "an upload",
		},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()

			n := New(tc.input)

			check := func(label, got, want string) {
				t.Helper()
				if got != want {
					t.Errorf("%s: got %q, want %q", label, got, want)
				}
			}

			check("Singular", n.Singular(), tc.singular)
			check("Plural", n.Plural(), tc.plural)
			check("Snake", n.Snake(), tc.snake)
			check("PluralSnake", n.PluralSnake(), tc.pluralSnake)
			check("Camel", n.Camel(), tc.camel)
			check("PluralCamel", n.PluralCamel(), tc.pluralCamel)
			check("Kebab", n.Kebab(), tc.kebab)
			check("PluralKebab", n.PluralKebab(), tc.pluralKebab)
			check("PackageName", n.PackageName(), tc.packageName)
			check("Abbreviation", n.Abbreviation(), tc.abbreviation)
			check("LowerAbbreviation", n.LowerAbbreviation(), tc.lowerAbbreviation)
			check("HumanReadable", n.HumanReadable(), tc.humanReadable)
			check("PluralHumanReadable", n.PluralHumanReadable(), tc.pluralHumanReadable)
			check("Article", n.Article(), tc.article)
			check("WithArticle", n.WithArticle(), tc.withArticle)
		})
	}
}

func TestNewWithPlural(t *testing.T) {
	t.Parallel()

	n := NewWithPlural("Person", "People")

	if n.Singular() != "Person" {
		t.Errorf("expected Person, got %q", n.Singular())
	}
	if n.Plural() != "People" {
		t.Errorf("expected People, got %q", n.Plural())
	}
	if n.PackageName() != "people" {
		t.Errorf("expected people, got %q", n.PackageName())
	}
	if n.PluralSnake() != "people" {
		t.Errorf("expected people, got %q", n.PluralSnake())
	}
}
