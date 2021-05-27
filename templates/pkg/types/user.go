package types

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func userDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("GoodStandingAccountStatus").ID("accountStatus").Op("=").Lit("good"),
			jen.ID("UnverifiedAccountStatus").ID("accountStatus").Op("=").Lit("unverified"),
			jen.ID("BannedUserAccountStatus").ID("accountStatus").Op("=").Lit("banned"),
			jen.ID("TerminatedUserReputation").ID("accountStatus").Op("=").Lit("terminated"),
			jen.ID("validTOTPTokenLength").Op("=").Lit(6),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("totpTokenLengthRule").Op("=").Qual("github.com/go-ozzo/ozzo-validation/v4", "Length").Call(
				jen.ID("validTOTPTokenLength"),
				jen.ID("validTOTPTokenLength"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("accountStatus").ID("string"),
			jen.ID("User").Struct(
				jen.ID("PasswordLastChangedOn").Op("*").ID("uint64"),
				jen.ID("ArchivedOn").Op("*").ID("uint64"),
				jen.ID("LastUpdatedOn").Op("*").ID("uint64"),
				jen.ID("TwoFactorSecretVerifiedOn").Op("*").ID("uint64"),
				jen.ID("AvatarSrc").Op("*").ID("string"),
				jen.ID("ExternalID").ID("string"),
				jen.ID("Username").ID("string"),
				jen.ID("ReputationExplanation").ID("string"),
				jen.ID("ServiceAccountStatus").ID("accountStatus"),
				jen.ID("TwoFactorSecret").ID("string"),
				jen.ID("HashedPassword").ID("string"),
				jen.ID("ServiceRoles").Index().ID("string"),
				jen.ID("ID").ID("uint64"),
				jen.ID("CreatedOn").ID("uint64"),
				jen.ID("RequiresPasswordChange").ID("bool"),
			),
			jen.ID("TestUserCreationConfig").Struct(
				jen.ID("Username").ID("string"),
				jen.ID("Password").ID("string"),
				jen.ID("HashedPassword").ID("string"),
				jen.ID("IsServiceAdmin").ID("bool"),
			),
			jen.ID("UserList").Struct(
				jen.ID("Users").Index().Op("*").ID("User"),
				jen.ID("Pagination"),
			),
			jen.ID("UserRegistrationInput").Struct(
				jen.ID("Username").ID("string"),
				jen.ID("Password").ID("string"),
			),
			jen.ID("UserDataStoreCreationInput").Struct(
				jen.ID("Username").ID("string"),
				jen.ID("HashedPassword").ID("string"),
				jen.ID("TwoFactorSecret").ID("string"),
			),
			jen.ID("UserCreationResponse").Struct(
				jen.ID("Username").ID("string"),
				jen.ID("AccountStatus").ID("accountStatus"),
				jen.ID("TwoFactorSecret").ID("string"),
				jen.ID("TwoFactorQRCode").ID("string"),
				jen.ID("CreatedUserID").ID("uint64"),
				jen.ID("CreatedOn").ID("uint64"),
				jen.ID("IsAdmin").ID("bool"),
			),
			jen.ID("UserLoginInput").Struct(
				jen.ID("Username").ID("string"),
				jen.ID("Password").ID("string"),
				jen.ID("TOTPToken").ID("string"),
			),
			jen.ID("PasswordUpdateInput").Struct(
				jen.ID("NewPassword").ID("string"),
				jen.ID("CurrentPassword").ID("string"),
				jen.ID("TOTPToken").ID("string"),
			),
			jen.ID("TOTPSecretRefreshInput").Struct(
				jen.ID("CurrentPassword").ID("string"),
				jen.ID("TOTPToken").ID("string"),
			),
			jen.ID("TOTPSecretVerificationInput").Struct(
				jen.ID("TOTPToken").ID("string"),
				jen.ID("UserID").ID("uint64"),
			),
			jen.ID("TOTPSecretRefreshResponse").Struct(
				jen.ID("TwoFactorQRCode").ID("string"),
				jen.ID("TwoFactorSecret").ID("string"),
			),
			jen.ID("AdminUserDataManager").Interface(jen.ID("UpdateUserReputation").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("input").Op("*").ID("UserReputationUpdateInput")).Params(jen.ID("error"))),
			jen.ID("UserDataManager").Interface(
				jen.ID("UserHasStatus").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("statuses").Op("...").ID("string")).Params(jen.ID("bool"), jen.ID("error")),
				jen.ID("GetUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("User"), jen.ID("error")),
				jen.ID("GetUserWithUnverifiedTwoFactorSecret").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("User"), jen.ID("error")),
				jen.ID("MarkUserTwoFactorSecretAsVerified").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("error")),
				jen.ID("GetUserByUsername").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("username").ID("string")).Params(jen.Op("*").ID("User"), jen.ID("error")),
				jen.ID("SearchForUsersByUsername").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("usernameQuery").ID("string")).Params(jen.Index().Op("*").ID("User"), jen.ID("error")),
				jen.ID("GetAllUsersCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")),
				jen.ID("GetUsers").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("QueryFilter")).Params(jen.Op("*").ID("UserList"), jen.ID("error")),
				jen.ID("CreateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("UserDataStoreCreationInput")).Params(jen.Op("*").ID("User"), jen.ID("error")),
				jen.ID("UpdateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("User"), jen.ID("changes").Index().Op("*").ID("FieldChangeSummary")).Params(jen.ID("error")),
				jen.ID("UpdateUserPassword").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("newHash").ID("string")).Params(jen.ID("error")),
				jen.ID("ArchiveUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("error")),
				jen.ID("GetAuditLogEntriesForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().Op("*").ID("AuditLogEntry"), jen.ID("error")),
			),
			jen.ID("UserDataService").Interface(
				jen.ID("ListHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("AuditEntryHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("CreateHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("ReadHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("SelfHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("UsernameSearchHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("NewTOTPSecretHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("TOTPSecretVerificationHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("UpdatePasswordHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("AvatarUploadHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("ArchiveHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("RegisterUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("registrationInput").Op("*").ID("UserRegistrationInput")).Params(jen.Op("*").ID("UserCreationResponse"), jen.ID("error")),
				jen.ID("VerifyUserTwoFactorSecret").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("TOTPSecretVerificationInput")).Params(jen.ID("error")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Update accepts a User as input and merges those values if they're set."),
		jen.Line(),
		jen.Func().Params(jen.ID("u").Op("*").ID("User")).ID("Update").Params(jen.ID("input").Op("*").ID("User")).Body(
			jen.If(jen.ID("input").Dot("Username").Op("!=").Lit("").Op("&&").ID("input").Dot("Username").Op("!=").ID("u").Dot("Username")).Body(
				jen.ID("u").Dot("Username").Op("=").ID("input").Dot("Username")),
			jen.If(jen.ID("input").Dot("HashedPassword").Op("!=").Lit("").Op("&&").ID("input").Dot("HashedPassword").Op("!=").ID("u").Dot("HashedPassword")).Body(
				jen.ID("u").Dot("HashedPassword").Op("=").ID("input").Dot("HashedPassword")),
			jen.If(jen.ID("input").Dot("TwoFactorSecret").Op("!=").Lit("").Op("&&").ID("input").Dot("TwoFactorSecret").Op("!=").ID("u").Dot("TwoFactorSecret")).Body(
				jen.ID("u").Dot("TwoFactorSecret").Op("=").ID("input").Dot("TwoFactorSecret")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("IsValidAccountStatus returns whether the provided string is a valid accountStatus."),
		jen.Line(),
		jen.Func().ID("IsValidAccountStatus").Params(jen.ID("s").ID("string")).Params(jen.ID("bool")).Body(
			jen.Switch(jen.ID("s")).Body(
				jen.Case(jen.ID("string").Call(jen.ID("GoodStandingAccountStatus")), jen.ID("string").Call(jen.ID("UnverifiedAccountStatus")), jen.ID("string").Call(jen.ID("BannedUserAccountStatus")), jen.ID("string").Call(jen.ID("TerminatedUserReputation"))).Body(
					jen.Return().ID("true")),
				jen.Default().Body(
					jen.Return().ID("false")),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("IsBanned is a handy helper function."),
		jen.Line(),
		jen.Func().Params(jen.ID("u").Op("*").ID("User")).ID("IsBanned").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("u").Dot("ServiceAccountStatus").Op("==").ID("BannedUserAccountStatus")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext ensures our provided TOTPSecretVerificationInput meets expectations."),
		jen.Line(),
		jen.Func().Params(jen.ID("i").Op("*").ID("TOTPSecretVerificationInput")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("i"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("i").Dot("UserID"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("i").Dot("TOTPToken"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					jen.ID("totpTokenLengthRule"),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext ensures our provided TOTPSecretVerificationInput meets expectations."),
		jen.Line(),
		jen.Func().Params(jen.ID("i").Op("*").ID("TOTPSecretVerificationInput")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("i"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("i").Dot("UserID"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("i").Dot("TOTPToken"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					jen.ID("totpTokenLengthRule"),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext ensures our provided TOTPSecretVerificationInput meets expectations."),
		jen.Line(),
		jen.Func().Params(jen.ID("i").Op("*").ID("TOTPSecretVerificationInput")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("i"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("i").Dot("UserID"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("i").Dot("TOTPToken"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					jen.ID("totpTokenLengthRule"),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("TOTPSecretRefreshInput")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext ensures our provided TOTPSecretVerificationInput meets expectations."),
		jen.Line(),
		jen.Func().Params(jen.ID("i").Op("*").ID("TOTPSecretVerificationInput")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("i"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("i").Dot("UserID"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("i").Dot("TOTPToken"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					jen.ID("totpTokenLengthRule"),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("TOTPSecretVerificationInput")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext ensures our provided TOTPSecretVerificationInput meets expectations."),
		jen.Line(),
		jen.Func().Params(jen.ID("i").Op("*").ID("TOTPSecretVerificationInput")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("i"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("i").Dot("UserID"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("i").Dot("TOTPToken"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					jen.ID("totpTokenLengthRule"),
				),
			)),
		jen.Line(),
	)

	return code
}
