package authorization

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authorizationTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestAuthorizations").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("service user"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					append([]jen.Code{
						jen.ID("t").Dot("Parallel").Call(),
						jen.Newline(),
					},
						buildServiceUserTests(proj)...,
					)...,
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("service admin"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					append([]jen.Code{
						jen.ID("t").Dot("Parallel").Call(),
						jen.Newline(),
					},
						buildServiceAdminTests(proj)...,
					)...,
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("account admin"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					append([]jen.Code{
						jen.ID("t").Dot("Parallel").Call(),
						jen.Newline(),
					},
						buildAccountAdminTests(proj)...,
					)...,
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("account member"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					append([]jen.Code{
						jen.ID("t").Dot("Parallel").Call(),
						jen.Newline(),
					},
						buildAccountMemberTests(proj)...,
					)...,
				),
			),
		),
		jen.Newline(),
	)

	return code
}

func buildServiceUserTests(proj *models.Project) []jen.Code {
	out := []jen.Code{}

	for _, typ := range proj.DataTypes {
		pn := typ.Name.Plural()

		out = append(out,
			jen.ID("assert").Dot("False").Call(
				jen.ID("t"),
				jen.IDf("CanCreate%s", pn).Call(jen.ID("serviceUserRoleName")),
			),
			jen.ID("assert").Dot("False").Call(
				jen.ID("t"),
				jen.IDf("CanSee%s", pn).Call(jen.ID("serviceUserRoleName")),
			),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.IDf("CanSearch%s", pn).Call(jen.ID("serviceUserRoleName")),
					)
				}
				return jen.Null()
			}(),
			jen.ID("assert").Dot("False").Call(
				jen.ID("t"),
				jen.IDf("CanUpdate%s", pn).Call(jen.ID("serviceUserRoleName")),
			),
			jen.ID("assert").Dot("False").Call(
				jen.ID("t"),
				jen.IDf("CanDelete%s", pn).Call(jen.ID("serviceUserRoleName")),
			),
		)
	}

	return out
}

func buildServiceAdminTests(proj *models.Project) []jen.Code {
	out := []jen.Code{}

	for _, typ := range proj.DataTypes {
		pn := typ.Name.Plural()

		out = append(out,
			jen.ID("assert").Dot("True").Call(
				jen.ID("t"),
				jen.IDf("CanCreate%s", pn).Call(jen.ID("serviceAdminRoleName")),
			),
			jen.ID("assert").Dot("True").Call(
				jen.ID("t"),
				jen.IDf("CanSee%s", pn).Call(jen.ID("serviceAdminRoleName")),
			),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.IDf("CanSearch%s", pn).Call(jen.ID("serviceAdminRoleName")),
					)
				}
				return jen.Null()
			}(),
			jen.ID("assert").Dot("True").Call(
				jen.ID("t"),
				jen.IDf("CanUpdate%s", pn).Call(jen.ID("serviceAdminRoleName")),
			),
			jen.ID("assert").Dot("True").Call(
				jen.ID("t"),
				jen.IDf("CanDelete%s", pn).Call(jen.ID("serviceAdminRoleName")),
			),
		)
	}

	return out
}

func buildAccountAdminTests(proj *models.Project) []jen.Code {
	out := []jen.Code{}

	for _, typ := range proj.DataTypes {
		pn := typ.Name.Plural()

		out = append(out,
			jen.ID("assert").Dot("True").Call(
				jen.ID("t"),
				jen.IDf("CanCreate%s", pn).Call(jen.ID("accountAdminRoleName")),
			),
			jen.ID("assert").Dot("True").Call(
				jen.ID("t"),
				jen.IDf("CanSee%s", pn).Call(jen.ID("accountAdminRoleName")),
			),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.IDf("CanSearch%s", pn).Call(jen.ID("accountAdminRoleName")),
					)
				}
				return jen.Null()
			}(),
			jen.ID("assert").Dot("True").Call(
				jen.ID("t"),
				jen.IDf("CanUpdate%s", pn).Call(jen.ID("accountAdminRoleName")),
			),
			jen.ID("assert").Dot("True").Call(
				jen.ID("t"),
				jen.IDf("CanDelete%s", pn).Call(jen.ID("accountAdminRoleName")),
			),
		)
	}

	return out
}

func buildAccountMemberTests(proj *models.Project) []jen.Code {
	out := []jen.Code{}

	for _, typ := range proj.DataTypes {
		pn := typ.Name.Plural()

		out = append(out,
			jen.ID("assert").Dot("True").Call(
				jen.ID("t"),
				jen.IDf("CanCreate%s", pn).Call(jen.ID("accountMemberRoleName")),
			),
			jen.ID("assert").Dot("True").Call(
				jen.ID("t"),
				jen.IDf("CanSee%s", pn).Call(jen.ID("accountMemberRoleName")),
			),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.IDf("CanSearch%s", pn).Call(jen.ID("accountMemberRoleName")),
					)
				}
				return jen.Null()
			}(),
			jen.ID("assert").Dot("True").Call(
				jen.ID("t"),
				jen.IDf("CanUpdate%s", pn).Call(jen.ID("accountMemberRoleName")),
			),
			jen.ID("assert").Dot("True").Call(
				jen.ID("t"),
				jen.IDf("CanDelete%s", pn).Call(jen.ID("accountMemberRoleName")),
			),
		)
	}

	return out
}
