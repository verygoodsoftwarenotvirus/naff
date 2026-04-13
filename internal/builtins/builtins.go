package builtins

import (
	"github.com/verygoodsoftwarenotvirus/naff/internal/config"
)

func boolPtr(v bool) *bool { return &v }

// BuiltinDomains returns the built-in infrastructure domains, filtered by the
// project's feature flags. These are merged with user-defined domains and run
// through the same template pipeline.
func BuiltinDomains(features config.Features) []config.Domain {
	var domains []config.Domain

	// Always-on domains.
	domains = append(domains, auditDomain())

	// Togglable domains.
	if features.FeatureEnabled(features.Webhooks, true) {
		domains = append(domains, webhooksDomain())
	}
	if features.FeatureEnabled(features.IssueReports, true) {
		domains = append(domains, issueReportsDomain())
	}
	if features.FeatureEnabled(features.Comments, true) {
		domains = append(domains, commentsDomain())
	}
	if features.FeatureEnabled(features.Notifications, true) {
		domains = append(domains, notificationsDomain())
	}
	if features.FeatureEnabled(features.UploadedMedia, true) {
		domains = append(domains, uploadedMediaDomain())
	}
	if features.FeatureEnabled(features.DataPrivacy, true) {
		domains = append(domains, dataPrivacyDomain())
	}
	if features.FeatureEnabled(features.Waitlists, false) {
		domains = append(domains, waitlistsDomain())
	}
	if features.FeatureEnabled(features.Payments, false) {
		domains = append(domains, paymentsDomain())
	}

	return domains
}

func auditDomain() config.Domain {
	return config.Domain{
		Name: "audit",
		Entities: []config.Entity{
			{
				Name:             "AuditLogEntry",
				BelongsToAccount: true,
				Fields: []config.Field{
					{Name: "ResourceType", Type: "string"},
					{Name: "RelevantID", Type: "string"},
					{Name: "EventType", Type: "string"},
					{Name: "Changes", Type: "string", Required: boolPtr(false), Omitempty: true},
				},
			},
		},
	}
}

func webhooksDomain() config.Domain {
	return config.Domain{
		Name: "webhooks",
		Entities: []config.Entity{
			{
				Name:             "Webhook",
				BelongsToAccount: true,
				CreatedByUser:    true,
				Fields: []config.Field{
					{Name: "Name", Type: "string"},
					{Name: "URL", Type: "string", Editable: boolPtr(false)},
					{Name: "Method", Type: "string", Editable: boolPtr(false)},
					{Name: "ContentType", Type: "string", Editable: boolPtr(false)},
				},
			},
		},
	}
}

func issueReportsDomain() config.Domain {
	return config.Domain{
		Name: "issuereports",
		Entities: []config.Entity{
			{
				Name:             "IssueReport",
				BelongsToAccount: true,
				CreatedByUser:    true,
				Fields: []config.Field{
					{Name: "IssueType", Type: "string"},
					{Name: "Details", Type: "string"},
					{Name: "RelevantTable", Type: "string", Required: boolPtr(false), Omitempty: true},
					{Name: "RelevantRecordID", Type: "string", Required: boolPtr(false), Omitempty: true},
				},
			},
		},
	}
}

func commentsDomain() config.Domain {
	return config.Domain{
		Name: "comments",
		Entities: []config.Entity{
			{
				Name:          "Comment",
				CreatedByUser: true,
				Fields: []config.Field{
					{Name: "Content", Type: "string"},
					{Name: "TargetType", Type: "string", Editable: boolPtr(false)},
					{Name: "ReferencedID", Type: "string", Editable: boolPtr(false)},
					{Name: "ParentCommentID", Type: "*string", Required: boolPtr(false), Editable: boolPtr(false), Omitempty: true},
				},
			},
		},
	}
}

func notificationsDomain() config.Domain {
	return config.Domain{
		Name: "notifications",
		Entities: []config.Entity{
			{
				Name:          "UserNotification",
				CreatedByUser: true,
				Fields: []config.Field{
					{Name: "Content", Type: "string", Editable: boolPtr(false)},
					{Name: "Status", Type: "string", Required: boolPtr(false), Omitempty: true},
				},
			},
		},
	}
}

func uploadedMediaDomain() config.Domain {
	return config.Domain{
		Name: "uploadedmedia",
		Entities: []config.Entity{
			{
				Name:          "UploadedMedia",
				CreatedByUser: true,
				Fields: []config.Field{
					{Name: "StoragePath", Type: "string"},
					{Name: "MimeType", Type: "string"},
				},
			},
		},
	}
}

func dataPrivacyDomain() config.Domain {
	return config.Domain{
		Name: "dataprivacy",
		Entities: []config.Entity{
			{
				Name:          "UserDataDisclosure",
				CreatedByUser: true,
				Fields: []config.Field{
					{Name: "Status", Type: "string", Editable: boolPtr(false)},
					{Name: "ReportID", Type: "string", Required: boolPtr(false), Editable: boolPtr(false), Omitempty: true},
				},
			},
		},
	}
}

func waitlistsDomain() config.Domain {
	return config.Domain{
		Name: "waitlists",
		Entities: []config.Entity{
			{
				Name: "Waitlist",
				Fields: []config.Field{
					{Name: "Name", Type: "string"},
					{Name: "Description", Type: "string"},
				},
			},
		},
	}
}

func paymentsDomain() config.Domain {
	return config.Domain{
		Name: "payments",
		Entities: []config.Entity{
			{
				Name:             "Subscription",
				BelongsToAccount: true,
				CreatedByUser:    true,
				Fields: []config.Field{
					{Name: "ProductID", Type: "string", Editable: boolPtr(false)},
					{Name: "Status", Type: "string"},
					{Name: "ExternalID", Type: "string", Required: boolPtr(false), Editable: boolPtr(false), Omitempty: true},
				},
			},
		},
	}
}
