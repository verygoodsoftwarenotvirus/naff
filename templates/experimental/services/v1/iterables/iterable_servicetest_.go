package iterables

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterableServiceTestDotGo(typ models.DataType) *jen.File {
	ret := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(ret)

	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	cn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()

	ret.Add(
		jen.Func().ID("buildTestService").Params().Params(jen.Op("*").ID("Service")).Block(
			jen.Return().Op("&").ID("Service").Valuesln(
				jen.ID("logger").Op(":").ID("noop").Dot("ProvideNoopLogger").Call(),
				jen.ID(fmt.Sprintf("%sCounter", uvn)).Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Values(),
				jen.ID(fmt.Sprintf("%sDatabase", uvn)).Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("userIDFetcher").Op(":").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
				jen.ID(fmt.Sprintf("%sIDFetcher", uvn)).Op(":").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
				jen.ID("encoderDecoder").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Values(),
				jen.ID("reporter").Op(":").ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID(fmt.Sprintf("TestProvide%sService", pn)).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectation").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.ID("expectation")).Dot("Return").Call(),
				jen.Line(),
				jen.Var().ID("ucp").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounterProvider").Op("=").Func().Paramsln(
					jen.ID("counterName").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "CounterName"),
					jen.ID("description").ID("string"),
				).Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounter"),
					jen.ID("error")).Block(
					jen.Return().List(jen.ID("uc"), jen.ID("nil")),
				),
				jen.Line(),
				jen.ID("idm").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("idm").Dot("On").Call(jen.Lit(fmt.Sprintf("GetAll%sCount", pn)), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("expectation"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("s"), jen.ID("err")).Op(":=").ID(fmt.Sprintf("Provide%sService", pn)).Callln(
					jen.Qual("context", "Background").Call(),
					jen.ID("noop").Dot("ProvideNoopLogger").Call(),
					jen.ID("idm"),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Values(),
					jen.ID("ucp"),
					jen.ID("nil"),
				),
				jen.Line(),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("s")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error providing unit counter"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectation").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.ID("expectation")).Dot("Return").Call(),
				jen.Line(),
				jen.Var().ID("ucp").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounterProvider").Op("=").Func().Paramsln(
					jen.ID("counterName").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "CounterName"),
					jen.ID("description").ID("string"),
				).Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounter"),
					jen.ID("error")).Block(
					jen.Return().List(jen.ID("uc"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				),
				jen.Line(),
				jen.ID("idm").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("idm").Dot("On").Call(jen.Lit(fmt.Sprintf("GetAll%sCount", pn)), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("expectation"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("s"), jen.ID("err")).Op(":=").ID(fmt.Sprintf("Provide%sService", pn)).Callln(
					jen.Qual("context", "Background").Call(),
					jen.ID("noop").Dot("ProvideNoopLogger").Call(),
					jen.ID("idm"),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Values(),
					jen.ID("ucp"),
					jen.ID("nil"),
				),
				jen.Line(),
				jen.ID("require").Dot("Nil").Call(jen.ID("t"), jen.ID("s")),
				jen.ID("require").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit(fmt.Sprintf("with error fetching %s count", cn)), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectation").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.ID("expectation")).Dot("Return").Call(),
				jen.Line(),
				jen.Var().ID("ucp").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounterProvider").Op("=").Func().Paramsln(
					jen.ID("counterName").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "CounterName"),
					jen.ID("description").ID("string"),
				).Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounter"),
					jen.ID("error")).Block(
					jen.Return().List(jen.ID("uc"), jen.ID("nil")),
				),
				jen.Line(),
				jen.ID("idm").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("idm").Dot("On").Call(jen.Lit(fmt.Sprintf("GetAll%sCount", pn)), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("expectation"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("s"), jen.ID("err")).Op(":=").ID(fmt.Sprintf("Provide%sService", pn)).Callln(
					jen.Qual("context", "Background").Call(),
					jen.ID("noop").Dot("ProvideNoopLogger").Call(),
					jen.ID("idm"),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Values(),
					jen.ID("ucp"),
					jen.ID("nil"),
				),
				jen.Line(),
				jen.ID("require").Dot("Nil").Call(jen.ID("t"), jen.ID("s")),
				jen.ID("require").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			)),
		),
		jen.Line(),
	)
	return ret
}
