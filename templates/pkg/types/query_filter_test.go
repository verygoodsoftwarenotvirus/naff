package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_queryFilterDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := queryFilterDotGo(proj)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	// MaxLimit is the maximum value for list queries.
	MaxLimit = 250
	// DefaultLimit represents how many results we return in a response by default.
	DefaultLimit = 20

	// SearchQueryKey is the query param key we use to find search queries in requests
	SearchQueryKey = "q"
	// LimitQueryKey is the query param key we use to specify a limit in a query
	LimitQueryKey = "limit"

	pageQueryKey          = "page"
	createdBeforeQueryKey = "createdBefore"
	createdAfterQueryKey  = "createdAfter"
	updatedBeforeQueryKey = "updatedBefore"
	updatedAfterQueryKey  = "updatedAfter"
	sortByQueryKey        = "sortBy"
)

// QueryFilter represents all the filters a user could apply to a list query.
type QueryFilter struct {
	Page          uint64   ` + "`" + `json:"page"` + "`" + `
	Limit         uint8    ` + "`" + `json:"limit"` + "`" + `
	CreatedAfter  uint64   ` + "`" + `json:"createdBefore,omitempty"` + "`" + `
	CreatedBefore uint64   ` + "`" + `json:"createdAfter,omitempty"` + "`" + `
	UpdatedAfter  uint64   ` + "`" + `json:"updatedBefore,omitempty"` + "`" + `
	UpdatedBefore uint64   ` + "`" + `json:"updatedAfter,omitempty"` + "`" + `
	SortBy        sortType ` + "`" + `json:"sortBy"` + "`" + `
}

// DefaultQueryFilter builds the default query filter.
func DefaultQueryFilter() *QueryFilter {
	return &QueryFilter{
		Page:   1,
		Limit:  DefaultLimit,
		SortBy: SortAscending,
	}
}

// FromParams overrides the core QueryFilter values with values retrieved from url.Params
func (qf *QueryFilter) FromParams(params url.Values) {
	if i, err := strconv.ParseUint(params.Get(pageQueryKey), 10, 64); err == nil {
		qf.Page = uint64(math.Max(float64(i), 1))
	}

	if i, err := strconv.ParseUint(params.Get(LimitQueryKey), 10, 64); err == nil {
		qf.Limit = uint8(math.Min(math.Max(float64(i), 0), MaxLimit))
	}

	if i, err := strconv.ParseUint(params.Get(createdBeforeQueryKey), 10, 64); err == nil {
		qf.CreatedBefore = uint64(math.Max(float64(i), 0))
	}

	if i, err := strconv.ParseUint(params.Get(createdAfterQueryKey), 10, 64); err == nil {
		qf.CreatedAfter = uint64(math.Max(float64(i), 0))
	}

	if i, err := strconv.ParseUint(params.Get(updatedBeforeQueryKey), 10, 64); err == nil {
		qf.UpdatedBefore = uint64(math.Max(float64(i), 0))
	}

	if i, err := strconv.ParseUint(params.Get(updatedAfterQueryKey), 10, 64); err == nil {
		qf.UpdatedAfter = uint64(math.Max(float64(i), 0))
	}

	switch strings.ToLower(params.Get(sortByQueryKey)) {
	case string(SortAscending):
		qf.SortBy = SortAscending
	case string(SortDescending):
		qf.SortBy = SortDescending
	}
}

// SetPage sets the current page with certain constraints.
func (qf *QueryFilter) SetPage(page uint64) {
	qf.Page = uint64(math.Max(1, float64(page)))
}

// QueryPage calculates a query page from the current filter values.
func (qf *QueryFilter) QueryPage() uint64 {
	return uint64(qf.Limit) * (qf.Page - 1)
}

// ToValues returns a url.Values from a QueryFilter
func (qf *QueryFilter) ToValues() url.Values {
	if qf == nil {
		return DefaultQueryFilter().ToValues()
	}

	v := url.Values{}
	if qf.Page != 0 {
		v.Set(pageQueryKey, strconv.FormatUint(qf.Page, 10))
	}
	if qf.Limit != 0 {
		v.Set(LimitQueryKey, strconv.FormatUint(uint64(qf.Limit), 10))
	}
	if qf.SortBy != "" {
		v.Set(sortByQueryKey, string(qf.SortBy))
	}
	if qf.CreatedBefore != 0 {
		v.Set(createdBeforeQueryKey, strconv.FormatUint(qf.CreatedBefore, 10))
	}
	if qf.CreatedAfter != 0 {
		v.Set(createdAfterQueryKey, strconv.FormatUint(qf.CreatedAfter, 10))
	}
	if qf.UpdatedBefore != 0 {
		v.Set(updatedBeforeQueryKey, strconv.FormatUint(qf.UpdatedBefore, 10))
	}
	if qf.UpdatedAfter != 0 {
		v.Set(updatedAfterQueryKey, strconv.FormatUint(qf.UpdatedAfter, 10))
	}

	return v
}

// ApplyToQueryBuilder applies the query filter to a query builder.
func (qf *QueryFilter) ApplyToQueryBuilder(queryBuilder squirrel.SelectBuilder, tableName string) squirrel.SelectBuilder {
	if qf == nil {
		return queryBuilder
	}

	const (
		createdOnKey = "created_on"
		updatedOnKey = "last_updated_on"
	)

	qf.SetPage(qf.Page)
	if qp := qf.QueryPage(); qp > 0 {
		queryBuilder = queryBuilder.Offset(qp)
	}

	if qf.Limit > 0 {
		queryBuilder = queryBuilder.Limit(uint64(qf.Limit))
	} else {
		queryBuilder = queryBuilder.Limit(MaxLimit)
	}

	if qf.CreatedAfter > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Gt{fmt.Sprintf("%s.%s", tableName, createdOnKey): qf.CreatedAfter})
	}

	if qf.CreatedBefore > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Lt{fmt.Sprintf("%s.%s", tableName, createdOnKey): qf.CreatedBefore})
	}

	if qf.UpdatedAfter > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Gt{fmt.Sprintf("%s.%s", tableName, updatedOnKey): qf.UpdatedAfter})
	}

	if qf.UpdatedBefore > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Lt{fmt.Sprintf("%s.%s", tableName, updatedOnKey): qf.UpdatedBefore})
	}

	return queryBuilder
}

// ExtractQueryFilter can extract a QueryFilter from a request.
func ExtractQueryFilter(req *http.Request) *QueryFilter {
	qf := &QueryFilter{}
	qf.FromParams(req.URL.Query())
	return qf
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildQueryFilterConstantDeclarations0(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildQueryFilterConstantDeclarations0()

		expected := `
package example

import ()

const (
	// MaxLimit is the maximum value for list queries.
	MaxLimit = 250
	// DefaultLimit represents how many results we return in a response by default.
	DefaultLimit = 20

	// SearchQueryKey is the query param key we use to find search queries in requests
	SearchQueryKey = "q"
	// LimitQueryKey is the query param key we use to specify a limit in a query
	LimitQueryKey = "limit"

	pageQueryKey          = "page"
	createdBeforeQueryKey = "createdBefore"
	createdAfterQueryKey  = "createdAfter"
	updatedBeforeQueryKey = "updatedBefore"
	updatedAfterQueryKey  = "updatedAfter"
	sortByQueryKey        = "sortBy"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildQueryFilter(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildQueryFilter()

		expected := `
package example

import ()

// QueryFilter represents all the filters a user could apply to a list query.
type QueryFilter struct {
	Page          uint64   ` + "`" + `json:"page"` + "`" + `
	Limit         uint8    ` + "`" + `json:"limit"` + "`" + `
	CreatedAfter  uint64   ` + "`" + `json:"createdBefore,omitempty"` + "`" + `
	CreatedBefore uint64   ` + "`" + `json:"createdAfter,omitempty"` + "`" + `
	UpdatedAfter  uint64   ` + "`" + `json:"updatedBefore,omitempty"` + "`" + `
	UpdatedBefore uint64   ` + "`" + `json:"updatedAfter,omitempty"` + "`" + `
	SortBy        sortType ` + "`" + `json:"sortBy"` + "`" + `
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildDefaultQueryFilter(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildDefaultQueryFilter()

		expected := `
package example

import ()

// DefaultQueryFilter builds the default query filter.
func DefaultQueryFilter() *QueryFilter {
	return &QueryFilter{
		Page:   1,
		Limit:  DefaultLimit,
		SortBy: SortAscending,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildFromParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildFromParams()

		expected := `
package example

import (
	"math"
	"net/url"
	"strconv"
	"strings"
)

// FromParams overrides the core QueryFilter values with values retrieved from url.Params
func (qf *QueryFilter) FromParams(params url.Values) {
	if i, err := strconv.ParseUint(params.Get(pageQueryKey), 10, 64); err == nil {
		qf.Page = uint64(math.Max(float64(i), 1))
	}

	if i, err := strconv.ParseUint(params.Get(LimitQueryKey), 10, 64); err == nil {
		qf.Limit = uint8(math.Min(math.Max(float64(i), 0), MaxLimit))
	}

	if i, err := strconv.ParseUint(params.Get(createdBeforeQueryKey), 10, 64); err == nil {
		qf.CreatedBefore = uint64(math.Max(float64(i), 0))
	}

	if i, err := strconv.ParseUint(params.Get(createdAfterQueryKey), 10, 64); err == nil {
		qf.CreatedAfter = uint64(math.Max(float64(i), 0))
	}

	if i, err := strconv.ParseUint(params.Get(updatedBeforeQueryKey), 10, 64); err == nil {
		qf.UpdatedBefore = uint64(math.Max(float64(i), 0))
	}

	if i, err := strconv.ParseUint(params.Get(updatedAfterQueryKey), 10, 64); err == nil {
		qf.UpdatedAfter = uint64(math.Max(float64(i), 0))
	}

	switch strings.ToLower(params.Get(sortByQueryKey)) {
	case string(SortAscending):
		qf.SortBy = SortAscending
	case string(SortDescending):
		qf.SortBy = SortDescending
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSetPage(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildSetPage()

		expected := `
package example

import (
	"math"
)

// SetPage sets the current page with certain constraints.
func (qf *QueryFilter) SetPage(page uint64) {
	qf.Page = uint64(math.Max(1, float64(page)))
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildQueryPage(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildQueryPage()

		expected := `
package example

import ()

// QueryPage calculates a query page from the current filter values.
func (qf *QueryFilter) QueryPage() uint64 {
	return uint64(qf.Limit) * (qf.Page - 1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildToValues(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildToValues()

		expected := `
package example

import (
	"net/url"
	"strconv"
)

// ToValues returns a url.Values from a QueryFilter
func (qf *QueryFilter) ToValues() url.Values {
	if qf == nil {
		return DefaultQueryFilter().ToValues()
	}

	v := url.Values{}
	if qf.Page != 0 {
		v.Set(pageQueryKey, strconv.FormatUint(qf.Page, 10))
	}
	if qf.Limit != 0 {
		v.Set(LimitQueryKey, strconv.FormatUint(uint64(qf.Limit), 10))
	}
	if qf.SortBy != "" {
		v.Set(sortByQueryKey, string(qf.SortBy))
	}
	if qf.CreatedBefore != 0 {
		v.Set(createdBeforeQueryKey, strconv.FormatUint(qf.CreatedBefore, 10))
	}
	if qf.CreatedAfter != 0 {
		v.Set(createdAfterQueryKey, strconv.FormatUint(qf.CreatedAfter, 10))
	}
	if qf.UpdatedBefore != 0 {
		v.Set(updatedBeforeQueryKey, strconv.FormatUint(qf.UpdatedBefore, 10))
	}
	if qf.UpdatedAfter != 0 {
		v.Set(updatedAfterQueryKey, strconv.FormatUint(qf.UpdatedAfter, 10))
	}

	return v
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildApplyToQueryBuilder(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildApplyToQueryBuilder()

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// ApplyToQueryBuilder applies the query filter to a query builder.
func (qf *QueryFilter) ApplyToQueryBuilder(queryBuilder squirrel.SelectBuilder, tableName string) squirrel.SelectBuilder {
	if qf == nil {
		return queryBuilder
	}

	const (
		createdOnKey = "created_on"
		updatedOnKey = "last_updated_on"
	)

	qf.SetPage(qf.Page)
	if qp := qf.QueryPage(); qp > 0 {
		queryBuilder = queryBuilder.Offset(qp)
	}

	if qf.Limit > 0 {
		queryBuilder = queryBuilder.Limit(uint64(qf.Limit))
	} else {
		queryBuilder = queryBuilder.Limit(MaxLimit)
	}

	if qf.CreatedAfter > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Gt{fmt.Sprintf("%s.%s", tableName, createdOnKey): qf.CreatedAfter})
	}

	if qf.CreatedBefore > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Lt{fmt.Sprintf("%s.%s", tableName, createdOnKey): qf.CreatedBefore})
	}

	if qf.UpdatedAfter > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Gt{fmt.Sprintf("%s.%s", tableName, updatedOnKey): qf.UpdatedAfter})
	}

	if qf.UpdatedBefore > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Lt{fmt.Sprintf("%s.%s", tableName, updatedOnKey): qf.UpdatedBefore})
	}

	return queryBuilder
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildExtractQueryFilter(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildExtractQueryFilter()

		expected := `
package example

import (
	"net/http"
)

// ExtractQueryFilter can extract a QueryFilter from a request.
func ExtractQueryFilter(req *http.Request) *QueryFilter {
	qf := &QueryFilter{}
	qf.FromParams(req.URL.Query())
	return qf
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
