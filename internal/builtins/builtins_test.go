package builtins

import (
	"testing"

	"github.com/verygoodsoftwarenotvirus/naff/internal/config"
)

func TestBuiltinDomains(t *testing.T) {
	t.Parallel()

	t.Run("all defaults enabled", func(t *testing.T) {
		t.Parallel()

		features := config.Features{}
		domains := BuiltinDomains(features)

		// Always-on: audit
		// Default-on: webhooks, issuereports, comments, notifications, uploadedmedia, dataprivacy
		// Default-off: waitlists, payments
		// auth/identity/oauth are NOT in builtins — shipped as hand-rolled templates via
		// pipeline.planDomainExtrasFiles(). See Phase 3b note in builtins.go.
		expectedCount := 7 // audit + 6 default-on
		if len(domains) != expectedCount {
			t.Errorf("expected %d domains with defaults, got %d", expectedCount, len(domains))
			for _, d := range domains {
				t.Logf("  - %s", d.Name)
			}
		}

		names := make(map[string]bool)
		for _, d := range domains {
			names[d.Name] = true
		}

		for _, expected := range []string{"audit", "webhooks", "issuereports", "comments", "notifications", "uploadedmedia", "dataprivacy"} {
			if !names[expected] {
				t.Errorf("expected domain %q to be present", expected)
			}
		}
		for _, notExpected := range []string{"waitlists", "payments"} {
			if names[notExpected] {
				t.Errorf("expected domain %q to NOT be present by default", notExpected)
			}
		}
	})

	t.Run("opt-in domains enabled", func(t *testing.T) {
		t.Parallel()

		features := config.Features{
			Waitlists: boolPtr(true),
			Payments:  boolPtr(true),
		}
		domains := BuiltinDomains(features)

		names := make(map[string]bool)
		for _, d := range domains {
			names[d.Name] = true
		}

		if !names["waitlists"] {
			t.Error("expected waitlists to be present when enabled")
		}
		if !names["payments"] {
			t.Error("expected payments to be present when enabled")
		}
	})

	t.Run("domains can be disabled", func(t *testing.T) {
		t.Parallel()

		features := config.Features{
			Webhooks:     boolPtr(false),
			IssueReports: boolPtr(false),
			Comments:     boolPtr(false),
		}
		domains := BuiltinDomains(features)

		names := make(map[string]bool)
		for _, d := range domains {
			names[d.Name] = true
		}

		if names["webhooks"] {
			t.Error("webhooks should be disabled")
		}
		if names["issuereports"] {
			t.Error("issuereports should be disabled")
		}
		if names["comments"] {
			t.Error("comments should be disabled")
		}
		// audit is always on
		if !names["audit"] {
			t.Error("audit should always be present")
		}
	})

	t.Run("all domains have valid entities", func(t *testing.T) {
		t.Parallel()

		features := config.Features{
			Waitlists: boolPtr(true),
			Payments:  boolPtr(true),
		}
		domains := BuiltinDomains(features)

		for _, d := range domains {
			if d.Name == "" {
				t.Error("domain has empty name")
			}
			if len(d.Entities) == 0 {
				t.Errorf("domain %q has no entities", d.Name)
			}
			for _, e := range d.Entities {
				if e.Name == "" {
					t.Errorf("domain %q has entity with empty name", d.Name)
				}
				if len(e.Fields) == 0 {
					t.Errorf("domain %q entity %q has no fields", d.Name, e.Name)
				}
				for _, f := range e.Fields {
					if f.Name == "" {
						t.Errorf("domain %q entity %q has field with empty name", d.Name, e.Name)
					}
					if f.Type == "" {
						t.Errorf("domain %q entity %q field %q has empty type", d.Name, e.Name, f.Name)
					}
				}
			}
		}
	})
}
