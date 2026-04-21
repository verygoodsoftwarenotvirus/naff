package config

import "testing"

func TestProjectMetaNameForms(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name                                                     string
		title, pascal, kebab, snake, titleKebab, abbrev, lAbbrev string
	}{
		{
			name:       "Dinner Done Better",
			title:      "Dinner Done Better",
			pascal:     "DinnerDoneBetter",
			kebab:      "dinner-done-better",
			snake:      "dinner_done_better",
			titleKebab: "Dinner-Done-Better",
			abbrev:     "DDB",
			lAbbrev:    "ddb",
		},
		{
			name:       "Kitchen Sink",
			title:      "Kitchen Sink",
			pascal:     "KitchenSink",
			kebab:      "kitchen-sink",
			snake:      "kitchen_sink",
			titleKebab: "Kitchen-Sink",
			abbrev:     "KS",
			lAbbrev:    "ks",
		},
		{
			name:       "Kitchen",
			title:      "Kitchen",
			pascal:     "Kitchen",
			kebab:      "kitchen",
			snake:      "kitchen",
			titleKebab: "Kitchen",
			abbrev:     "K",
			lAbbrev:    "k",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p := ProjectMeta{Name: tc.name}
			if got := p.Title(); got != tc.title {
				t.Errorf("Title() = %q, want %q", got, tc.title)
			}
			if got := p.Pascal(); got != tc.pascal {
				t.Errorf("Pascal() = %q, want %q", got, tc.pascal)
			}
			if got := p.Kebab(); got != tc.kebab {
				t.Errorf("Kebab() = %q, want %q", got, tc.kebab)
			}
			if got := p.Snake(); got != tc.snake {
				t.Errorf("Snake() = %q, want %q", got, tc.snake)
			}
			if got := p.TitleKebab(); got != tc.titleKebab {
				t.Errorf("TitleKebab() = %q, want %q", got, tc.titleKebab)
			}
			if got := p.Abbrev(); got != tc.abbrev {
				t.Errorf("Abbrev() = %q, want %q", got, tc.abbrev)
			}
			if got := p.LowerAbbrev(); got != tc.lAbbrev {
				t.Errorf("LowerAbbrev() = %q, want %q", got, tc.lAbbrev)
			}
		})
	}
}

func TestLinkDerivedNames(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name                                          string
		link                                          Link
		fieldName, columnName, goType, jsonTag, sqlOD string
	}{
		{
			name:       "required link, target only",
			link:       Link{Target: "Recipe"},
			fieldName:  "RecipeID",
			columnName: "recipe_id",
			goType:     "string",
			jsonTag:    "recipeID",
			sqlOD:      "RESTRICT",
		},
		{
			name:       "as overrides target for naming",
			link:       Link{Target: "ValidIngredient", As: "PurchasedMeasurementUnit"},
			fieldName:  "PurchasedMeasurementUnitID",
			columnName: "purchased_measurement_unit_id",
			goType:     "string",
			jsonTag:    "purchasedMeasurementUnitID",
			sqlOD:      "RESTRICT",
		},
		{
			name:       "optional link is pointer with omitempty tag",
			link:       Link{Target: "ValidIngredient", Optional: true},
			fieldName:  "ValidIngredientID",
			columnName: "valid_ingredient_id",
			goType:     "*string",
			jsonTag:    "validIngredientID,omitempty",
			sqlOD:      "RESTRICT",
		},
		{
			name:       "on_delete cascade uppercases",
			link:       Link{Target: "Recipe", OnDelete: LinkOnDeleteCascade},
			fieldName:  "RecipeID",
			columnName: "recipe_id",
			goType:     "string",
			jsonTag:    "recipeID",
			sqlOD:      "CASCADE",
		},
		{
			name:       "on_delete set_null spaces uppercased",
			link:       Link{Target: "Recipe", Optional: true, OnDelete: LinkOnDeleteSetNull},
			fieldName:  "RecipeID",
			columnName: "recipe_id",
			goType:     "*string",
			jsonTag:    "recipeID,omitempty",
			sqlOD:      "SET NULL",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.link.FieldName(); got != tc.fieldName {
				t.Errorf("FieldName() = %q, want %q", got, tc.fieldName)
			}
			if got := tc.link.ColumnName(); got != tc.columnName {
				t.Errorf("ColumnName() = %q, want %q", got, tc.columnName)
			}
			if got := tc.link.GoType(); got != tc.goType {
				t.Errorf("GoType() = %q, want %q", got, tc.goType)
			}
			if got := tc.link.JSONTag(); got != tc.jsonTag {
				t.Errorf("JSONTag() = %q, want %q", got, tc.jsonTag)
			}
			if got := tc.link.SQLOnDelete(); got != tc.sqlOD {
				t.Errorf("SQLOnDelete() = %q, want %q", got, tc.sqlOD)
			}
		})
	}

}

func TestLinkAsField(t *testing.T) {
	t.Parallel()

	t.Run("required default-non-editable link", func(t *testing.T) {
		t.Parallel()

		f := Link{Target: "Recipe"}.AsField()
		if f.Name != "RecipeID" || f.Type != "string" {
			t.Errorf("name/type = %q/%q, want RecipeID/string", f.Name, f.Type)
		}
		if !f.IsRequired() {
			t.Error("IsRequired() = false, want true")
		}
		if !f.IsCreatable() {
			t.Error("IsCreatable() = false, want true")
		}
		if f.IsEditable() {
			t.Error("IsEditable() = true, want false (links default to non-editable)")
		}
		if f.Omitempty {
			t.Error("Omitempty = true, want false for required link")
		}
	})

	t.Run("optional editable link", func(t *testing.T) {
		t.Parallel()

		yes := true
		f := Link{Target: "ValidIngredient", Optional: true, Editable: &yes}.AsField()
		if f.Name != "ValidIngredientID" || f.Type != "*string" {
			t.Errorf("name/type = %q/%q, want ValidIngredientID/*string", f.Name, f.Type)
		}
		if f.IsRequired() {
			t.Error("IsRequired() = true, want false")
		}
		if !f.IsEditable() {
			t.Error("IsEditable() = false, want true")
		}
		if !f.Omitempty {
			t.Error("Omitempty = false, want true for optional link")
		}
	})
}

func TestEntityAllFields(t *testing.T) {
	t.Parallel()

	e := Entity{
		Name: "MealComponent",
		Fields: []Field{
			{Name: "ComponentType", Type: "string"},
			{Name: "RecipeScale", Type: "float32"},
		},
		LinksTo: []Link{
			{Target: "Recipe"},
		},
	}

	all := e.AllFields()
	if len(all) != 3 {
		t.Fatalf("AllFields len = %d, want 3", len(all))
	}
	if all[0].Name != "ComponentType" || all[1].Name != "RecipeScale" {
		t.Errorf("scalar fields out of order: %+v", all[:2])
	}
	if all[2].Name != "RecipeID" || all[2].Type != "string" {
		t.Errorf("synthetic link field = %+v, want RecipeID/string", all[2])
	}

	syn := e.SyntheticLinkFields()
	if len(syn) != 1 || syn[0].Name != "RecipeID" {
		t.Errorf("SyntheticLinkFields() = %+v, want one RecipeID entry", syn)
	}
}
