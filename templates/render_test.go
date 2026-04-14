package templates

import (
	"strings"
	"testing"

	"github.com/verygoodsoftwarenotvirus/naff/internal/config"
	"github.com/verygoodsoftwarenotvirus/naff/internal/pipeline"
	"github.com/verygoodsoftwarenotvirus/naff/internal/renderer"
)

func testProject() *config.Project {
	return &config.Project{
		ProjectMeta: config.ProjectMeta{
			Name:           "TestProject",
			Module:         "github.com/example/testproject",
			PlatformModule: "github.com/primandproper/platform",
		},
	}
}

func testDomain() *config.Domain {
	return &config.Domain{
		Name:     "issuereports",
		TypeName: "IssueReports",
	}
}

func testEntity() *config.Entity {
	falseVal := false
	return &config.Entity{
		Name:             "IssueReport",
		BelongsToAccount: true,
		CreatedByUser:    true,
		Fields: []config.Field{
			{Name: "IssueType", Type: "string"},
			{Name: "Details", Type: "string"},
			{Name: "RelevantTable", Type: "string", Required: &falseVal, Omitempty: true},
			{Name: "RelevantRecordID", Type: "string", Required: &falseVal, Omitempty: true},
		},
	}
}

func TestEntityTemplate(t *testing.T) {
	t.Parallel()

	r := renderer.New(FS)
	ctx := pipeline.TemplateContext{
		Project: testProject(),
		Domain:  testDomain(),
		Entity:  testEntity(),
	}

	output, err := r.Render("domain/entity.go.tmpl", ctx)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	// Verify key patterns exist in the output.
	checks := []string{
		"package issuereports",
		"IssueReportCreatedServiceEventType",
		"IssueReport struct",
		"IssueReportCreationRequestInput struct",
		"IssueReportDatabaseCreationInput struct",
		"IssueReportUpdateRequestInput struct",
		"IssueReportDataManager interface",
		"func (x *IssueReport) Update(input *IssueReportUpdateRequestInput)",
		"ValidateWithContext",
		"gob.Register",
		"CreatedByUser",
		"BelongsToAccount",
		`json:"issueType"`,
		`json:"relevantTable,omitempty"`,
		"validation.Required",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("expected output to contain %q\n\nOutput:\n%s", check, output)
		}
	}

	// Verify the output is valid Go (can be formatted).
	_, err = renderer.FormatGo([]byte(output))
	if err != nil {
		t.Errorf("output is not valid Go: %v\n\nOutput:\n%s", err, output)
	}
}

func TestRepositoryTemplate(t *testing.T) {
	t.Parallel()

	r := renderer.New(FS)
	ctx := pipeline.TemplateContext{
		Project: testProject(),
		Domain:  testDomain(),
		Entity:  testEntity(),
	}

	output, err := r.Render("domain/repository.go.tmpl", ctx)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	checks := []string{
		"package issuereports",
		"Repository interface",
		"IssueReportDataManager",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("expected output to contain %q", check)
		}
	}
}

func TestKeysTemplate(t *testing.T) {
	t.Parallel()

	r := renderer.New(FS)
	ctx := pipeline.TemplateContext{
		Project: testProject(),
		Domain:  testDomain(),
		Entity:  testEntity(),
	}

	output, err := r.Render("domain_keys/keys.go.tmpl", ctx)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	checks := []string{
		"package keys",
		"IssueReportIDKey",
		`"issue_report"`,
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("expected output to contain %q", check)
		}
	}
}

func TestConvertersTemplate(t *testing.T) {
	t.Parallel()

	r := renderer.New(FS)
	ctx := pipeline.TemplateContext{
		Project: testProject(),
		Domain:  testDomain(),
		Entity:  testEntity(),
	}

	output, err := r.Render("domain_converters/converters.go.tmpl", ctx)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	checks := []string{
		"package converters",
		"ConvertIssueReportToIssueReportUpdateRequestInput",
		"ConvertIssueReportCreationRequestInputToIssueReportDatabaseCreationInput",
		"identifiers.New()",
		"userID string",
		"accountID string",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("expected output to contain %q", check)
		}
	}
}

func TestMockRepositoryTemplate(t *testing.T) {
	t.Parallel()

	r := renderer.New(FS)
	ctx := pipeline.TemplateContext{
		Project: testProject(),
		Domain:  testDomain(),
		Entity:  testEntity(),
	}

	output, err := r.Render("domain_mock/repository.go.tmpl", ctx)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	checks := []string{
		"package mock",
		"Repository struct",
		"mock.Mock",
		"GetIssueReport",
		"GetIssueReports",
		"GetIssueReportsForAccount",
		"CreateIssueReport",
		"UpdateIssueReport",
		"ArchiveIssueReport",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("expected output to contain %q", check)
		}
	}
}

func TestFakesTemplate(t *testing.T) {
	t.Parallel()

	r := renderer.New(FS)
	ctx := pipeline.TemplateContext{
		Project: testProject(),
		Domain:  testDomain(),
		Entity:  testEntity(),
	}

	output, err := r.Render("domain_fakes/fake.go.tmpl", ctx)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	checks := []string{
		"package fakes",
		"BuildFakeIssueReport()",
		"BuildFakeIssueReportCreationRequestInput()",
		"BuildFakeIssueReportDatabaseCreationInput()",
		"BuildFakeIssueReportUpdateRequestInput()",
		"BuildFakeID()",
		"BuildFakeTime()",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("expected output to contain %q", check)
		}
	}
}

func TestMigrationTemplate(t *testing.T) {
	t.Parallel()

	r := renderer.New(FS)
	ctx := pipeline.TemplateContext{
		Project: testProject(),
		Domain:  testDomain(),
		Entity:  testEntity(),
	}

	output, err := r.Render("migrations/migration.sql.tmpl", ctx)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	checks := []string{
		"CREATE TABLE IF NOT EXISTS issue_reports",
		"id TEXT NOT NULL PRIMARY KEY",
		"issue_type TEXT NOT NULL",
		"details TEXT NOT NULL",
		"relevant_table TEXT",
		"created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()",
		"last_updated_at TIMESTAMP WITH TIME ZONE",
		"archived_at TIMESTAMP WITH TIME ZONE",
		"created_by_user TEXT NOT NULL REFERENCES users",
		"belongs_to_account TEXT NOT NULL REFERENCES accounts",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("expected output to contain %q\n\nOutput:\n%s", check, output)
		}
	}
}

func TestQueryCodegenTemplate(t *testing.T) {
	t.Parallel()

	r := renderer.New(FS)
	ctx := pipeline.TemplateContext{
		Project: testProject(),
		Domain:  testDomain(),
		Entity:  testEntity(),
	}

	output, err := r.Render("codegen/queries.go.tmpl", ctx)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	checks := []string{
		"package main",
		`issueReportTableName = "issue_reports"`,
		"issueReportColumns",
		"buildIssueReportsQueries",
		"CreateIssueReport",
		"UpdateIssueReport",
		"ArchiveIssueReport",
		"GetIssueReport",
		"CheckIssueReportExistence",
		"GetIssueReports",
		"GetIssueReportsForAccount",
		"builq.Builder",
		"filterForInsert",
		"filterForUpdate",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("expected output to contain %q\n\nOutput:\n%s", check, output)
		}
	}
}

func TestRepositoryClientTemplate(t *testing.T) {
	t.Parallel()

	r := renderer.New(FS)
	ctx := pipeline.TemplateContext{
		Project: testProject(),
		Domain:  testDomain(),
		Entity:  testEntity(),
	}

	output, err := r.Render("repository/client.go.tmpl", ctx)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	checks := []string{
		"package issuereports",
		`o11yName = "issue_report_db_client"`,
		"type repository struct",
		"database.Client",
		"generatedQuerier",
		"auditLogEntryRepo",
		"ProvideIssueReportsRepository",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("expected output to contain %q\n\nOutput:\n%s", check, output)
		}
	}
}

func TestRepositoryEntityTemplate(t *testing.T) {
	t.Parallel()

	r := renderer.New(FS)
	ctx := pipeline.TemplateContext{
		Project: testProject(),
		Domain:  testDomain(),
		Entity:  testEntity(),
	}

	output, err := r.Render("repository/entity.go.tmpl", ctx)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	checks := []string{
		"package issuereports",
		"GetIssueReport",
		"GetIssueReports",
		"GetIssueReportsForAccount",
		"CreateIssueReport",
		"UpdateIssueReport",
		"ArchiveIssueReport",
		"r.generatedQuerier",
		"r.auditLogEntryRepo.CreateAuditLogEntry",
		"tx.Commit()",
		"r.RollbackTransaction",
		"audit.AuditLogEventTypeCreated",
		"audit.AuditLogEventTypeUpdated",
		"audit.AuditLogEventTypeArchived",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("expected output to contain %q", check)
		}
	}
}

func TestRepositoryDoTemplate(t *testing.T) {
	t.Parallel()

	r := renderer.New(FS)
	ctx := pipeline.TemplateContext{
		Project: testProject(),
		Domain:  testDomain(),
		Entity:  testEntity(),
	}

	output, err := r.Render("repository/do.go.tmpl", ctx)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	checks := []string{
		"package issuereports",
		"RegisterIssueReportsRepository",
		"do.Provide",
		"do.MustInvoke",
		"samber/do/v2",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("expected output to contain %q", check)
		}
	}
}

func TestSqlcBlockTemplate(t *testing.T) {
	t.Parallel()

	r := renderer.New(FS)
	ctx := pipeline.TemplateContext{
		Project: testProject(),
		Domain:  testDomain(),
		Entity:  testEntity(),
	}

	output, err := r.Render("sqlc/block.yaml.tmpl", ctx)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	checks := []string{
		"engine: \"postgresql\"",
		"internal/repositories/postgres/issuereports/sqlc_queries",
		"internal/repositories/postgres/issuereports/generated",
		"emit_interface: true",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("expected output to contain %q", check)
		}
	}
}
