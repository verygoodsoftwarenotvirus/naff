package webhooks

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func httpRoutesTestDotGo() *jen.File {
	ret := jen.NewFile("webhooks")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestWebhooksService_List").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"WebhookList",
			).Valuesln(jen.ID("Webhooks").Op(":").Index().ID("models").Dot(
				"Webhook",
			).Valuesln(jen.Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")))),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("GetWebhooks"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("requestingUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("s").Dot(
				"ListHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusOK")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with no rows returned"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("GetWebhooks"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("requestingUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"WebhookList",
			)).Call(jen.ID("nil")), jen.Qual("database/sql", "ErrNoRows")),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("s").Dot(
				"ListHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusOK")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error fetching webhooks from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("GetWebhooks"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("requestingUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"WebhookList",
			)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("s").Dot(
				"ListHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusInternalServerError")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error encoding response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"WebhookList",
			).Valuesln(jen.ID("Webhooks").Op(":").Index().ID("models").Dot(
				"Webhook",
			).Valuesln()),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("GetWebhooks"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("requestingUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("s").Dot(
				"ListHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusOK")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestValidateWebhook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"WebhookCreationInput",
			).Valuesln(jen.ID("Method").Op(":").Qual("net/http", "MethodPost"), jen.ID("URL").Op(":").Lit("https://todo.verygoodsoftwarenotvirus.ru")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("validateWebhook").Call(jen.ID("exampleInput"))),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with invalid method"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"WebhookCreationInput",
			).Valuesln(jen.ID("Method").Op(":").Lit(` MEATLOAF `), jen.ID("URL").Op(":").Lit("https://todo.verygoodsoftwarenotvirus.ru")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("validateWebhook").Call(jen.ID("exampleInput"))),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with invalid url"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"WebhookCreationInput",
			).Valuesln(jen.ID("Method").Op(":").Qual("net/http", "MethodPost"), jen.ID("URL").Op(":").Lit("%zzzzz")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("validateWebhook").Call(jen.ID("exampleInput"))),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestWebhooksService_Create").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")),
			jen.ID("mc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Valuesln(),
			jen.ID("mc").Dot(
				"On",
			).Call(jen.Lit("Increment"), jen.ID("mock").Dot(
				"Anything",
			)),
			jen.ID("s").Dot(
				"webhookCounter",
			).Op("=").ID("mc"),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("CreateWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"WebhookCreationInput",
			).Valuesln(jen.ID("Name").Op(":").ID("expected").Dot(
				"Name",
			), jen.ID("Method").Op(":").Qual("net/http", "MethodPatch")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("CreateMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("s").Dot(
				"CreateHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusCreated")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with invalid webhook request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")),
			jen.ID("mc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Valuesln(),
			jen.ID("mc").Dot(
				"On",
			).Call(jen.Lit("Increment"), jen.ID("mock").Dot(
				"Anything",
			)),
			jen.ID("s").Dot(
				"webhookCounter",
			).Op("=").ID("mc"),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("CreateWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"WebhookCreationInput",
			).Valuesln(jen.ID("Method").Op(":").Qual("net/http", "MethodPost"), jen.ID("URL").Op(":").Lit("%zzzzz")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("CreateMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("s").Dot(
				"CreateHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusBadRequest")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("without input attached"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("s").Dot(
				"CreateHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusBadRequest")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error creating webhook"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("CreateWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"Webhook",
			)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"WebhookCreationInput",
			).Valuesln(jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"), jen.ID("Name").Op(":").ID("expected").Dot(
				"Name",
			)),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("CreateMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("s").Dot(
				"CreateHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusInternalServerError")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error encoding response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")),
			jen.ID("mc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Valuesln(),
			jen.ID("mc").Dot(
				"On",
			).Call(jen.Lit("Increment"), jen.ID("mock").Dot(
				"Anything",
			)),
			jen.ID("s").Dot(
				"webhookCounter",
			).Op("=").ID("mc"),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("CreateWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"WebhookCreationInput",
			).Valuesln(jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"), jen.ID("Name").Op(":").ID("expected").Dot(
				"Name",
			)),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("CreateMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("s").Dot(
				"CreateHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusCreated")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestWebhooksService_Read").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("s").Dot(
				"webhookIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("expected").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("GetWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"ID",
			), jen.ID("requestingUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("s").Dot(
				"ReadHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusOK")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with no such webhook in database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("s").Dot(
				"webhookIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("expected").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("GetWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"ID",
			), jen.ID("requestingUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"Webhook",
			)).Call(jen.ID("nil")), jen.Qual("database/sql", "ErrNoRows")),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("s").Dot(
				"ReadHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusNotFound")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error fetching webhook from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("s").Dot(
				"webhookIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("expected").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("GetWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"ID",
			), jen.ID("requestingUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"Webhook",
			)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("s").Dot(
				"ReadHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusInternalServerError")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error encoding response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("s").Dot(
				"webhookIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("expected").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("GetWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"ID",
			), jen.ID("requestingUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("s").Dot(
				"ReadHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusOK")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestWebhooksService_Update").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")),
			jen.ID("mc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Valuesln(),
			jen.ID("mc").Dot(
				"On",
			).Call(jen.Lit("Increment"), jen.ID("mock").Dot(
				"Anything",
			)),
			jen.ID("s").Dot(
				"webhookCounter",
			).Op("=").ID("mc"),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("s").Dot(
				"webhookIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("expected").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("GetWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"ID",
			), jen.ID("requestingUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("UpdateWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"WebhookUpdateInput",
			).Valuesln(jen.ID("Name").Op(":").ID("expected").Dot(
				"Name",
			)),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UpdateMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("s").Dot(
				"UpdateHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusOK")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("without update input"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("s").Dot(
				"UpdateHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusBadRequest")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with no rows fetching webhook"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("s").Dot(
				"webhookIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("expected").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("GetWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"ID",
			), jen.ID("requestingUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"Webhook",
			)).Call(jen.ID("nil")), jen.Qual("database/sql", "ErrNoRows")),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"WebhookUpdateInput",
			).Valuesln(jen.ID("Name").Op(":").ID("expected").Dot(
				"Name",
			)),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UpdateMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("s").Dot(
				"UpdateHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusNotFound")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error fetching webhook"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("s").Dot(
				"webhookIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("expected").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("GetWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"ID",
			), jen.ID("requestingUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"Webhook",
			)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"WebhookUpdateInput",
			).Valuesln(jen.ID("Name").Op(":").ID("expected").Dot(
				"Name",
			)),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UpdateMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("s").Dot(
				"UpdateHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusInternalServerError")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error updating webhook"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")),
			jen.ID("mc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Valuesln(),
			jen.ID("mc").Dot(
				"On",
			).Call(jen.Lit("Increment"), jen.ID("mock").Dot(
				"Anything",
			)),
			jen.ID("s").Dot(
				"webhookCounter",
			).Op("=").ID("mc"),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("s").Dot(
				"webhookIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("expected").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("GetWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"ID",
			), jen.ID("requestingUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("UpdateWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"WebhookUpdateInput",
			).Valuesln(jen.ID("Name").Op(":").ID("expected").Dot(
				"Name",
			)),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UpdateMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("s").Dot(
				"UpdateHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusInternalServerError")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error encoding response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")),
			jen.ID("mc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Valuesln(),
			jen.ID("mc").Dot(
				"On",
			).Call(jen.Lit("Increment"), jen.ID("mock").Dot(
				"Anything",
			)),
			jen.ID("s").Dot(
				"webhookCounter",
			).Op("=").ID("mc"),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("s").Dot(
				"webhookIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("expected").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("GetWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"ID",
			), jen.ID("requestingUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("UpdateWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"WebhookUpdateInput",
			).Valuesln(jen.ID("Name").Op(":").ID("expected").Dot(
				"Name",
			)),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UpdateMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("s").Dot(
				"UpdateHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusOK")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestWebhooksService_Archive").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")),
			jen.ID("mc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Valuesln(),
			jen.ID("mc").Dot(
				"On",
			).Call(jen.Lit("Decrement")).Dot(
				"Return",
			).Call(),
			jen.ID("s").Dot(
				"webhookCounter",
			).Op("=").ID("mc"),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("s").Dot(
				"webhookIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("expected").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("ArchiveWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"ID",
			), jen.ID("requestingUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("s").Dot(
				"ArchiveHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusNoContent")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with no webhook in database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("s").Dot(
				"webhookIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("expected").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("ArchiveWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"ID",
			), jen.ID("requestingUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.Qual("database/sql", "ErrNoRows")),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("s").Dot(
				"ArchiveHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusNotFound")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error reading from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("requestingUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("requestingUser").Dot(
					"ID",
				),
			),
			jen.ID("s").Dot(
				"webhookIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("expected").Dot(
					"ID",
				),
			),
			jen.ID("id").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("id").Dot(
				"On",
			).Call(jen.Lit("ArchiveWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"ID",
			), jen.ID("requestingUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"webhookDatabase",
			).Op("=").ID("id"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("s").Dot(
				"ArchiveHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusInternalServerError")),
		)),
	),

		jen.Line(),
	)
	return ret
}
