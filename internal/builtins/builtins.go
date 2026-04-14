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
	domains = append(domains, authDomain())
	domains = append(domains, identityDomain())
	domains = append(domains, oauthDomain())

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

// authDomain provides the auth service scaffolding: password reset tokens and
// user sessions. DTO-shaped auth messages (UserLoginInput, PasswordUpdateInput,
// TOTPSecretRefreshInput, TOTPSecretVerificationInput, ChangeActiveAccountInput)
// are not entities — they're request shapes consumed by the hand-rolled
// services/auth/handlers/authentication package, deferred to Phase 3c.
func authDomain() config.Domain {
	return config.Domain{
		Name: "auth",
		Entities: []config.Entity{
			{
				Name: "PasswordResetToken",
				Fields: []config.Field{
					{Name: "Token", Type: "string", Editable: boolPtr(false)},
					{Name: "BelongsToUser", Type: "string", Editable: boolPtr(false)},
					{Name: "ExpiresAt", Type: "string", Editable: boolPtr(false)},
					{Name: "RedeemedAt", Type: "*string", Required: boolPtr(false), Editable: boolPtr(false), Omitempty: true},
				},
			},
			{
				Name:          "UserSession",
				CreatedByUser: true,
				Fields: []config.Field{
					{Name: "ClientIP", Type: "string", Editable: boolPtr(false)},
					{Name: "UserAgent", Type: "string", Editable: boolPtr(false)},
					{Name: "DeviceName", Type: "string"},
					{Name: "LoginMethod", Type: "string", Editable: boolPtr(false)},
					{Name: "LastActiveAt", Type: "string"},
					{Name: "ExpiresAt", Type: "string", Editable: boolPtr(false)},
					{Name: "IsCurrent", Type: "bool"},
				},
			},
		},
	}
}

// identityDomain provides user/account scaffolding: User, Account, AccountInvitation,
// AccountUserMembership, DataCollection. DTO-shaped identity messages
// (UserRegistrationInput, ModifyUserPermissionsInput, AccountOwnershipTransferInput,
// UserAccountStatusUpdateInput) are not entities and are deferred to Phase 3b's
// planDomainExtrasFiles() extension point.
//
// Target's User message embeds uploaded_media.UploadedMedia avatar and Account embeds
// repeated AccountUserMembershipWithUser members. Naff has no "embedded entity" concept
// so these relational fields are elided; they live on the repository/service layer
// at DDB-specific composition time.
func identityDomain() config.Domain {
	return config.Domain{
		Name: "identity",
		Entities: []config.Entity{
			{
				Name:          "User",
				CreatedByUser: false,
				Fields: []config.Field{
					{Name: "Username", Type: "string"},
					{Name: "EmailAddress", Type: "string"},
					{Name: "HashedPassword", Type: "string", Editable: boolPtr(false)},
					{Name: "TwoFactorSecret", Type: "string", Editable: boolPtr(false)},
					{Name: "TwoFactorSecretVerifiedAt", Type: "*string", Required: boolPtr(false), Editable: boolPtr(false), Omitempty: true},
					{Name: "PasswordLastChangedAt", Type: "*string", Required: boolPtr(false), Editable: boolPtr(false), Omitempty: true},
					{Name: "RequiresPasswordChange", Type: "bool"},
					{Name: "AccountStatus", Type: "string"},
					{Name: "AccountStatusExplanation", Type: "string"},
					{Name: "ServiceRole", Type: "string"},
					{Name: "FirstName", Type: "string"},
					{Name: "LastName", Type: "string"},
					{Name: "Birthday", Type: "*string", Required: boolPtr(false), Omitempty: true},
					{Name: "EmailAddressVerifiedAt", Type: "*string", Required: boolPtr(false), Editable: boolPtr(false), Omitempty: true},
					{Name: "LastAcceptedTermsOfService", Type: "*string", Required: boolPtr(false), Editable: boolPtr(false), Omitempty: true},
					{Name: "LastAcceptedPrivacyPolicy", Type: "*string", Required: boolPtr(false), Editable: boolPtr(false), Omitempty: true},
				},
			},
			{
				Name: "Account",
				Fields: []config.Field{
					{Name: "Name", Type: "string"},
					{Name: "BillingStatus", Type: "string"},
					{Name: "ContactPhone", Type: "string"},
					{Name: "AddressLine1", Type: "string"},
					{Name: "AddressLine2", Type: "string"},
					{Name: "City", Type: "string"},
					{Name: "State", Type: "string"},
					{Name: "ZipCode", Type: "string"},
					{Name: "Country", Type: "string"},
					{Name: "Latitude", Type: "*float64", Required: boolPtr(false), Omitempty: true},
					{Name: "Longitude", Type: "*float64", Required: boolPtr(false), Omitempty: true},
					{Name: "PaymentProcessorCustomerID", Type: "string", Editable: boolPtr(false)},
					{Name: "SubscriptionPlanID", Type: "*string", Required: boolPtr(false), Omitempty: true},
					{Name: "BelongsToUser", Type: "string", Editable: boolPtr(false)},
					{Name: "WebhookEncryptionKey", Type: "string", Editable: boolPtr(false)},
				},
			},
			{
				Name: "AccountInvitation",
				Fields: []config.Field{
					{Name: "FromUser", Type: "string", Editable: boolPtr(false)},
					{Name: "ToUser", Type: "*string", Required: boolPtr(false), Editable: boolPtr(false), Omitempty: true},
					{Name: "ToEmail", Type: "string", Editable: boolPtr(false)},
					{Name: "ToName", Type: "string", Editable: boolPtr(false)},
					{Name: "Note", Type: "string"},
					{Name: "StatusNote", Type: "string"},
					{Name: "Token", Type: "string", Editable: boolPtr(false)},
					{Name: "DestinationAccount", Type: "string", Editable: boolPtr(false)},
					{Name: "ExpiresAt", Type: "string", Editable: boolPtr(false)},
					{Name: "Status", Type: "string"},
				},
			},
			{
				Name: "AccountUserMembership",
				Fields: []config.Field{
					{Name: "BelongsToUser", Type: "string", Editable: boolPtr(false)},
					{Name: "BelongsToAccount", Type: "string", Editable: boolPtr(false)},
					{Name: "AccountRole", Type: "string"},
					{Name: "DefaultAccount", Type: "bool"},
				},
			},
			{
				Name:          "DataCollection",
				CreatedByUser: true,
				Fields: []config.Field{
					{Name: "Status", Type: "string"},
					{Name: "ReportContents", Type: "string", Editable: boolPtr(false)},
				},
			},
		},
	}
}

// oauthDomain provides OAuth2 client scaffolding. OAuth2ClientToken's proto
// representation uses google.protobuf.Duration for its *_expires_at fields;
// naff emits strings here since Duration support would require another funcmap
// entry and sqlc/postgres type mappings that aren't useful outside this domain.
// Phase 3b can refine if needed.
func oauthDomain() config.Domain {
	return config.Domain{
		Name: "oauth",
		Entities: []config.Entity{
			{
				Name: "OAuth2Client",
				Fields: []config.Field{
					{Name: "Name", Type: "string"},
					{Name: "Description", Type: "string"},
					{Name: "ClientID", Type: "string", Editable: boolPtr(false)},
					{Name: "ClientSecret", Type: "string", Editable: boolPtr(false)},
				},
			},
			{
				Name:          "OAuth2ClientToken",
				CreatedByUser: true,
				Fields: []config.Field{
					{Name: "ClientID", Type: "string", Editable: boolPtr(false)},
					{Name: "Scope", Type: "string", Editable: boolPtr(false)},
					{Name: "Code", Type: "string", Editable: boolPtr(false)},
					{Name: "CodeChallenge", Type: "string", Editable: boolPtr(false)},
					{Name: "CodeChallengeMethod", Type: "string", Editable: boolPtr(false)},
					{Name: "RedirectURI", Type: "string", Editable: boolPtr(false)},
					{Name: "Access", Type: "string", Editable: boolPtr(false)},
					{Name: "Refresh", Type: "string", Editable: boolPtr(false)},
					{Name: "CodeCreatedAt", Type: "*string", Required: boolPtr(false), Editable: boolPtr(false), Omitempty: true},
					{Name: "AccessCreatedAt", Type: "*string", Required: boolPtr(false), Editable: boolPtr(false), Omitempty: true},
					{Name: "RefreshCreatedAt", Type: "*string", Required: boolPtr(false), Editable: boolPtr(false), Omitempty: true},
					{Name: "CodeExpiresAt", Type: "string", Editable: boolPtr(false)},
					{Name: "AccessExpiresAt", Type: "string", Editable: boolPtr(false)},
					{Name: "RefreshExpiresAt", Type: "string", Editable: boolPtr(false)},
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
