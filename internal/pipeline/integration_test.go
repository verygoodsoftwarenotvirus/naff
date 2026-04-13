package pipeline

import (
	"strings"
	"testing"

	"github.com/verygoodsoftwarenotvirus/naff/internal/config"
)

func TestPlanFilesIncludesBuiltins(t *testing.T) {
	t.Parallel()

	cfg := &config.Project{
		ProjectMeta: config.ProjectMeta{
			Name:           "TestProject",
			Module:         "github.com/example/test",
			PlatformModule: "github.com/primandproper/platform",
		},
		Features: config.Features{},
		Domains: []config.Domain{
			{
				Name: "widgets",
				Entities: []config.Entity{
					{
						Name:             "Widget",
						BelongsToAccount: true,
						Fields: []config.Field{
							{Name: "Name", Type: "string"},
						},
					},
				},
			},
		},
	}

	p := &Pipeline{
		config:    cfg,
		outputDir: "/tmp/test-output",
	}

	planned := p.planFiles()

	// Should have files for built-in domains + user domain.
	if len(planned) == 0 {
		t.Fatal("expected planned files, got none")
	}

	// Check that built-in domains are included.
	var hasAudit, hasWebhooks, hasIssueReports, hasWidget bool
	for _, pf := range planned {
		if strings.Contains(pf.OutputPath, "/audit/") {
			hasAudit = true
		}
		if strings.Contains(pf.OutputPath, "/webhooks/") {
			hasWebhooks = true
		}
		if strings.Contains(pf.OutputPath, "/issuereports/") {
			hasIssueReports = true
		}
		if strings.Contains(pf.OutputPath, "/widgets/") {
			hasWidget = true
		}
	}

	if !hasAudit {
		t.Error("expected audit domain files in plan")
	}
	if !hasWebhooks {
		t.Error("expected webhooks domain files in plan")
	}
	if !hasIssueReports {
		t.Error("expected issuereports domain files in plan")
	}
	if !hasWidget {
		t.Error("expected widgets (user domain) files in plan")
	}
}

func TestPlanFilesRespectsFeatureFlags(t *testing.T) {
	t.Parallel()

	falseVal := false
	cfg := &config.Project{
		ProjectMeta: config.ProjectMeta{
			Name:           "TestProject",
			Module:         "github.com/example/test",
			PlatformModule: "github.com/primandproper/platform",
		},
		Features: config.Features{
			Webhooks:     &falseVal,
			IssueReports: &falseVal,
		},
		Domains: []config.Domain{
			{
				Name: "widgets",
				Entities: []config.Entity{
					{
						Name: "Widget",
						Fields: []config.Field{
							{Name: "Name", Type: "string"},
						},
					},
				},
			},
		},
	}

	p := &Pipeline{
		config:    cfg,
		outputDir: "/tmp/test-output",
	}

	planned := p.planFiles()

	for _, pf := range planned {
		if strings.Contains(pf.OutputPath, "/webhooks/") {
			t.Error("webhooks should NOT be in plan when disabled")
		}
		if strings.Contains(pf.OutputPath, "/issuereports/") {
			t.Error("issuereports should NOT be in plan when disabled")
		}
	}

	// audit is always-on
	var hasAudit bool
	for _, pf := range planned {
		if strings.Contains(pf.OutputPath, "/audit/") {
			hasAudit = true
			break
		}
	}
	if !hasAudit {
		t.Error("audit should always be in plan")
	}
}
