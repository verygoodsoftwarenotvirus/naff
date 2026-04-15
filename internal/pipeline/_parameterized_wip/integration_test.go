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

func TestPlanFilesIncludesIOSFiles(t *testing.T) {
	t.Parallel()

	falseVal := false
	cfg := &config.Project{
		ProjectMeta: config.ProjectMeta{
			Name:          "TestProject",
			IOSModuleName: "TestProject",
			IOSBundleID:   "com.example.testproject",
		},
		Features: config.Features{},
		Targets: config.Targets{
			Backend: &falseVal,
			IOS:     true,
		},
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

	if len(planned) == 0 {
		t.Fatal("expected planned files, got none")
	}

	// Check iOS project files.
	var hasPackageSwift, hasAppSwift, hasAuthManager, hasGRPCClient bool
	// Check iOS entity files.
	var hasWidgetModel, hasWidgetInput, hasWidgetProto bool
	// Check iOS domain files.
	var hasWidgetService bool
	// Should NOT have backend files.
	var hasBackendFile bool

	for _, pf := range planned {
		if strings.HasSuffix(pf.OutputPath, "Package.swift") {
			hasPackageSwift = true
		}
		if strings.HasSuffix(pf.OutputPath, "App.generated.swift") {
			hasAppSwift = true
		}
		if strings.HasSuffix(pf.OutputPath, "AuthManager.generated.swift") {
			hasAuthManager = true
		}
		if strings.HasSuffix(pf.OutputPath, "GRPCClient.generated.swift") {
			hasGRPCClient = true
		}
		if strings.HasSuffix(pf.OutputPath, "Widget.generated.swift") {
			hasWidgetModel = true
		}
		if strings.HasSuffix(pf.OutputPath, "WidgetInput.generated.swift") {
			hasWidgetInput = true
		}
		if strings.Contains(pf.OutputPath, "Proto/") && strings.Contains(pf.OutputPath, "widget") {
			hasWidgetProto = true
		}
		if strings.HasSuffix(pf.OutputPath, "WidgetsService.generated.swift") {
			hasWidgetService = true
		}
		if strings.HasSuffix(pf.OutputPath, ".go") {
			hasBackendFile = true
		}
	}

	if !hasPackageSwift {
		t.Error("expected Package.swift in iOS plan")
	}
	if !hasAppSwift {
		t.Error("expected App.generated.swift in iOS plan")
	}
	if !hasAuthManager {
		t.Error("expected AuthManager.generated.swift in iOS plan")
	}
	if !hasGRPCClient {
		t.Error("expected GRPCClient.generated.swift in iOS plan")
	}
	if !hasWidgetModel {
		t.Error("expected Widget.generated.swift in iOS plan")
	}
	if !hasWidgetInput {
		t.Error("expected WidgetInput.generated.swift in iOS plan")
	}
	if !hasWidgetProto {
		t.Error("expected widget proto files in iOS plan")
	}
	if !hasWidgetService {
		t.Error("expected WidgetsService.generated.swift in iOS plan")
	}
	if hasBackendFile {
		t.Error("should NOT have backend .go files when only iOS target is enabled")
	}
}

func TestPlanFilesIncludesConsumerApp(t *testing.T) {
	t.Parallel()

	trueVal := true
	cfg := &config.Project{
		ProjectMeta: config.ProjectMeta{
			Name:           "TestProject",
			Module:         "github.com/example/test",
			PlatformModule: "github.com/primandproper/platform",
		},
		Features: config.Features{
			ConsumerApp: &trueVal,
		},
		Domains: []config.Domain{
			{
				Name: "widgets",
				Entities: []config.Entity{
					{
						Name:             "Widget",
						ConsumerEditable: true,
						Fields: []config.Field{
							{Name: "Name", Type: "string"},
						},
					},
					{
						Name: "Gadget",
						Fields: []config.Field{
							{Name: "Label", Type: "string"},
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

	var (
		hasPackageJSON    bool
		hasConsumerLayout bool
		hasNav            bool
		hasTypesTS        bool
		hasAPITS          bool
		hasWidgetList     bool
		hasWidgetDetail   bool
		hasWidgetCreate   bool
		hasWidgetEdit     bool
		hasGadgetList     bool
		hasGadgetCreate   bool
	)

	for _, pf := range planned {
		switch {
		case pf.OutputPath == "frontend/consumer/package.json":
			hasPackageJSON = true
		case pf.OutputPath == "frontend/consumer/src/routes/+layout.svelte":
			hasConsumerLayout = true
		case strings.HasSuffix(pf.OutputPath, "Nav.generated.svelte"):
			hasNav = true
		case pf.OutputPath == "frontend/consumer/src/lib/api/widgets/types.generated.ts":
			hasTypesTS = true
		case pf.OutputPath == "frontend/consumer/src/lib/api/widgets/api.generated.ts":
			hasAPITS = true
		case pf.OutputPath == "frontend/consumer/src/routes/widgets/+page.svelte":
			hasWidgetList = true
		case pf.OutputPath == "frontend/consumer/src/routes/widgets/[id]/+page.svelte":
			hasWidgetDetail = true
		case pf.OutputPath == "frontend/consumer/src/routes/widgets/new/+page.svelte":
			hasWidgetCreate = true
		case pf.OutputPath == "frontend/consumer/src/routes/widgets/[id]/edit/+page.svelte":
			hasWidgetEdit = true
		case pf.OutputPath == "frontend/consumer/src/routes/gadgets/+page.svelte":
			hasGadgetList = true
		case pf.OutputPath == "frontend/consumer/src/routes/gadgets/new/+page.svelte":
			hasGadgetCreate = true
		}
	}

	if !hasPackageJSON {
		t.Error("expected consumer package.json")
	}
	if !hasConsumerLayout {
		t.Error("expected consumer +layout.svelte")
	}
	if !hasNav {
		t.Error("expected consumer Nav.generated.svelte")
	}
	if !hasTypesTS {
		t.Error("expected consumer types.generated.ts for widgets domain")
	}
	if !hasAPITS {
		t.Error("expected consumer api.generated.ts for widgets domain")
	}
	if !hasWidgetList {
		t.Error("expected consumer widget list page")
	}
	if !hasWidgetDetail {
		t.Error("expected consumer widget detail page")
	}
	if !hasWidgetCreate {
		t.Error("expected consumer widget create page (consumer_editable)")
	}
	if !hasWidgetEdit {
		t.Error("expected consumer widget edit page (consumer_editable)")
	}
	if !hasGadgetList {
		t.Error("expected consumer gadget list page")
	}
	if hasGadgetCreate {
		t.Error("gadget should NOT have create page (not consumer_editable)")
	}
}

func TestPlanFilesIncludesAdminApp(t *testing.T) {
	t.Parallel()

	trueVal := true
	cfg := &config.Project{
		ProjectMeta: config.ProjectMeta{
			Name:           "TestProject",
			Module:         "github.com/example/test",
			PlatformModule: "github.com/primandproper/platform",
		},
		Features: config.Features{
			AdminApp: &trueVal,
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

	var (
		hasPackageJSON bool
		hasAdminLayout bool
		hasWidgetList  bool
		hasWidgetNew   bool
		hasWidgetEdit  bool
	)

	for _, pf := range planned {
		switch {
		case pf.OutputPath == "frontend/admin/package.json":
			hasPackageJSON = true
		case pf.OutputPath == "frontend/admin/src/routes/+layout.svelte":
			hasAdminLayout = true
		case pf.OutputPath == "frontend/admin/src/routes/widgets/+page.svelte":
			hasWidgetList = true
		case pf.OutputPath == "frontend/admin/src/routes/widgets/new/+page.svelte":
			hasWidgetNew = true
		case pf.OutputPath == "frontend/admin/src/routes/widgets/[id]/edit/+page.svelte":
			hasWidgetEdit = true
		}
	}

	if !hasPackageJSON {
		t.Error("expected admin package.json")
	}
	if !hasAdminLayout {
		t.Error("expected admin +layout.svelte")
	}
	if !hasWidgetList {
		t.Error("expected admin widget list page")
	}
	if !hasWidgetNew {
		t.Error("expected admin widget create page (admin always has CRUD)")
	}
	if !hasWidgetEdit {
		t.Error("expected admin widget edit page (admin always has CRUD)")
	}
}

func TestPlanFilesExcludesFrontendWhenDisabled(t *testing.T) {
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
		if strings.HasPrefix(pf.OutputPath, "frontend/") {
			t.Errorf("should NOT have frontend files when consumer_app/admin_app are disabled, found: %s", pf.OutputPath)
		}
	}
}

func TestPlanFilesBackendOnlyExcludesIOS(t *testing.T) {
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
		if strings.Contains(pf.OutputPath, "ios/") {
			t.Errorf("should NOT have iOS files when iOS target is disabled, found: %s", pf.OutputPath)
		}
	}
}
