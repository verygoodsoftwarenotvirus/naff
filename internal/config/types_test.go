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
