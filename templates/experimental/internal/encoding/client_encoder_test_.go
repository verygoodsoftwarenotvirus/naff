package encoding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func clientEncoderTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestProvideClientEncoder").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("ProvideClientEncoder").Call(
							jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
							jen.ID("ContentTypeJSON"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_clientEncoder_Unmarshal").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with JSON"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("e").Op(":=").ID("ProvideClientEncoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.ID("expected").Op(":=").Op("&").ID("example").Valuesln(jen.ID("Name").Op(":").Lit("name")),
					jen.ID("actual").Op(":=").Op("&").ID("example").Valuesln(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("e").Dot("Unmarshal").Call(
							jen.ID("ctx"),
							jen.Index().ID("byte").Call(jen.Lit(`{"name": "name"}`)),
							jen.Op("&").ID("actual"),
						),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with XML"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("e").Op(":=").ID("ProvideClientEncoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeXML"),
					),
					jen.ID("expected").Op(":=").Op("&").ID("example").Valuesln(jen.ID("Name").Op(":").Lit("name")),
					jen.ID("actual").Op(":=").Op("&").ID("example").Valuesln(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("e").Dot("Unmarshal").Call(
							jen.ID("ctx"),
							jen.Index().ID("byte").Call(jen.Lit(`<example><name>name</name></example>`)),
							jen.Op("&").ID("actual"),
						),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("e").Op(":=").ID("ProvideClientEncoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.ID("actual").Op(":=").Op("&").ID("example").Valuesln(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("e").Dot("Unmarshal").Call(
							jen.ID("ctx"),
							jen.Index().ID("byte").Call(jen.Lit(`{"name"   `)),
							jen.Op("&").ID("actual"),
						),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("actual").Dot("Name"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_clientEncoder_Encode").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with JSON"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("e").Op(":=").ID("ProvideClientEncoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("e").Dot("Encode").Call(
							jen.ID("ctx"),
							jen.ID("res"),
							jen.Op("&").ID("example").Valuesln(jen.ID("Name").Op(":").ID("t").Dot("Name").Call()),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with XML"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("e").Op(":=").ID("ProvideClientEncoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeXML"),
					),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("e").Dot("Encode").Call(
							jen.ID("ctx"),
							jen.ID("res"),
							jen.Op("&").ID("example").Valuesln(jen.ID("Name").Op(":").ID("t").Dot("Name").Call()),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("e").Op(":=").ID("ProvideClientEncoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("e").Dot("Encode").Call(
							jen.ID("ctx"),
							jen.ID("nil"),
							jen.Op("&").ID("broken").Valuesln(jen.ID("Name").Op(":").Qual("encoding/json", "Number").Call(jen.ID("t").Dot("Name").Call())),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_clientEncoder_EncodeReader").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with JSON"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("e").Op(":=").ID("ProvideClientEncoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("e").Dot("EncodeReader").Call(
						jen.ID("ctx"),
						jen.Op("&").ID("example").Valuesln(jen.ID("Name").Op(":").ID("t").Dot("Name").Call()),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with XML"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("e").Op(":=").ID("ProvideClientEncoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeXML"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("e").Dot("EncodeReader").Call(
						jen.ID("ctx"),
						jen.Op("&").ID("example").Valuesln(jen.ID("Name").Op(":").ID("t").Dot("Name").Call()),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("e").Op(":=").ID("ProvideClientEncoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("e").Dot("EncodeReader").Call(
						jen.ID("ctx"),
						jen.Op("&").ID("broken").Valuesln(jen.ID("Name").Op(":").Qual("encoding/json", "Number").Call(jen.ID("t").Dot("Name").Call())),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
