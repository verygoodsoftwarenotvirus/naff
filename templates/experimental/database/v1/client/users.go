package client

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func usersDotGo() *jen.File {
	ret := jen.NewFile("dbclient")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("models").Dot(
		"UserDataManager",
	).Op("=").Parens(jen.Op("*").ID("Client")).Call(jen.ID("nil")).Var().ID("ErrUserExists").Op("=").Qual("errors", "New").Call(jen.Lit("error: username already exists")),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("attachUsernameToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("username").ID("string")).Block(
		jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
			jen.ID("span").Dot(
				"AddAttributes",
			).Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("username"), jen.ID("username"))),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetUser fetches a user").Params(jen.ID("c").Op("*").ID("Client")).ID("GetUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("models").Dot(
		"User",
	), jen.ID("error")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetUser")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
		jen.ID("c").Dot(
			"logger",
		).Dot(
			"WithValue",
		).Call(jen.Lit("user_id"), jen.ID("userID")).Dot(
			"Debug",
		).Call(jen.Lit("GetUser called")),
		jen.Return().ID("c").Dot(
			"querier",
		).Dot(
			"GetUser",
		).Call(jen.ID("ctx"), jen.ID("userID")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetUserByUsername fetches a user by their username").Params(jen.ID("c").Op("*").ID("Client")).ID("GetUserByUsername").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("username").ID("string")).Params(jen.Op("*").ID("models").Dot(
		"User",
	), jen.ID("error")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetUserByUsername")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.ID("attachUsernameToSpan").Call(jen.ID("span"), jen.ID("username")),
		jen.ID("c").Dot(
			"logger",
		).Dot(
			"WithValue",
		).Call(jen.Lit("username"), jen.ID("username")).Dot(
			"Debug",
		).Call(jen.Lit("GetUserByUsername called")),
		jen.Return().ID("c").Dot(
			"querier",
		).Dot(
			"GetUserByUsername",
		).Call(jen.ID("ctx"), jen.ID("username")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetUserCount fetches a count of users from the database that meet a particular filter").Params(jen.ID("c").Op("*").ID("Client")).ID("GetUserCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("models").Dot(
		"QueryFilter",
	)).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetUserCount")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.ID("attachFilterToSpan").Call(jen.ID("span"), jen.ID("filter")),
		jen.ID("c").Dot(
			"logger",
		).Dot(
			"Debug",
		).Call(jen.Lit("GetUserCount called")),
		jen.Return().ID("c").Dot(
			"querier",
		).Dot(
			"GetUserCount",
		).Call(jen.ID("ctx"), jen.ID("filter")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetUsers fetches a list of users from the database that meet a particular filter").Params(jen.ID("c").Op("*").ID("Client")).ID("GetUsers").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("models").Dot(
		"QueryFilter",
	)).Params(jen.Op("*").ID("models").Dot(
		"UserList",
	), jen.ID("error")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetUsers")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.ID("attachFilterToSpan").Call(jen.ID("span"), jen.ID("filter")),
		jen.ID("c").Dot(
			"logger",
		).Dot(
			"WithValue",
		).Call(jen.Lit("filter"), jen.ID("filter")).Dot(
			"Debug",
		).Call(jen.Lit("GetUsers called")),
		jen.Return().ID("c").Dot(
			"querier",
		).Dot(
			"GetUsers",
		).Call(jen.ID("ctx"), jen.ID("filter")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// CreateUser creates a user").Params(jen.ID("c").Op("*").ID("Client")).ID("CreateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("models").Dot(
		"UserInput",
	)).Params(jen.Op("*").ID("models").Dot(
		"User",
	), jen.ID("error")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("CreateUser")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.ID("attachUsernameToSpan").Call(jen.ID("span"), jen.ID("input").Dot(
			"Username",
		)),
		jen.ID("c").Dot(
			"logger",
		).Dot(
			"WithValue",
		).Call(jen.Lit("username"), jen.ID("input").Dot(
			"Username",
		)).Dot(
			"Debug",
		).Call(jen.Lit("CreateUser called")),
		jen.Return().ID("c").Dot(
			"querier",
		).Dot(
			"CreateUser",
		).Call(jen.ID("ctx"), jen.ID("input")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// UpdateUser receives a complete User struct and updates its record in the database.").Comment("// NOTE: this function uses the ID provided in the input to make its query.").Params(jen.ID("c").Op("*").ID("Client")).ID("UpdateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("models").Dot(
		"User",
	)).Params(jen.ID("error")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("UpdateUser")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.ID("attachUsernameToSpan").Call(jen.ID("span"), jen.ID("updated").Dot(
			"Username",
		)),
		jen.ID("c").Dot(
			"logger",
		).Dot(
			"WithValue",
		).Call(jen.Lit("username"), jen.ID("updated").Dot(
			"Username",
		)).Dot(
			"Debug",
		).Call(jen.Lit("UpdateUser called")),
		jen.Return().ID("c").Dot(
			"querier",
		).Dot(
			"UpdateUser",
		).Call(jen.ID("ctx"), jen.ID("updated")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ArchiveUser archives a user").Params(jen.ID("c").Op("*").ID("Client")).ID("ArchiveUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("error")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("ArchiveUser")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
		jen.ID("c").Dot(
			"logger",
		).Dot(
			"WithValue",
		).Call(jen.Lit("user_id"), jen.ID("userID")).Dot(
			"Debug",
		).Call(jen.Lit("ArchiveUser called")),
		jen.Return().ID("c").Dot(
			"querier",
		).Dot(
			"ArchiveUser",
		).Call(jen.ID("ctx"), jen.ID("userID")),
	),

		jen.Line(),
	)
	return ret
}
