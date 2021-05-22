package requests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_helpersDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := helpersDotGo(proj)

		expected := `
package example

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

// argIsNotPointer checks an argument and returns whether or not it is a pointer.
func argIsNotPointer(i interface{}) (notAPointer bool, err error) {
	if i == nil || reflect.TypeOf(i).Kind() != reflect.Ptr {
		return true, errors.New("value is not a pointer")
	}
	return false, nil
}

// argIsNotNil checks an argument and returns whether or not it is nil.
func argIsNotNil(i interface{}) (isNil bool, err error) {
	if i == nil {
		return true, errors.New("value is nil")
	}
	return false, nil
}

// argIsNotPointerOrNil does what it says on the tin. This function
// is primarily useful for detecting if a destination value is valid
// before decoding an HTTP response, for instance.
func argIsNotPointerOrNil(i interface{}) error {
	if nn, err := argIsNotNil(i); nn || err != nil {
		return err
	}

	if np, err := argIsNotPointer(i); np || err != nil {
		return err
	}

	return nil
}

// unmarshalBody takes an HTTP response and JSON decodes its
// body into a destination value. ` + "`" + `dest` + "`" + ` must be a non-nil
// pointer to an object. Ideally, response is also not nil.
// The error returned here should only ever be received in
// testing, and should never be encountered by an end-user.
func unmarshalBody(ctx context.Context, res *http.Response, dest interface{}) error {
	_, span := tracing.StartSpan(ctx, "unmarshalBody")
	defer span.End()

	if err := argIsNotPointerOrNil(dest); err != nil {
		return err
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= http.StatusBadRequest {
		apiErr := &v1.ErrorResponse{}
		if err = json.Unmarshal(bodyBytes, &apiErr); err != nil {
			return fmt.Errorf("unmarshaling error: %w", err)
		}
		return apiErr
	}

	if err = json.Unmarshal(bodyBytes, &dest); err != nil {
		return fmt.Errorf("unmarshaling body: %w", err)
	}

	return nil
}

// createBodyFromStruct takes any value in and returns an io.Reader
// for placement within http.NewRequest's last argument.
func createBodyFromStruct(in interface{}) (io.Reader, error) {
	out, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(out), nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildArgIsNotPointer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildArgIsNotPointer()

		expected := `
package example

import (
	"errors"
	"reflect"
)

// argIsNotPointer checks an argument and returns whether or not it is a pointer.
func argIsNotPointer(i interface{}) (notAPointer bool, err error) {
	if i == nil || reflect.TypeOf(i).Kind() != reflect.Ptr {
		return true, errors.New("value is not a pointer")
	}
	return false, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildArgIsNotNil(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildArgIsNotNil()

		expected := `
package example

import (
	"errors"
)

// argIsNotNil checks an argument and returns whether or not it is nil.
func argIsNotNil(i interface{}) (isNil bool, err error) {
	if i == nil {
		return true, errors.New("value is nil")
	}
	return false, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildArgIsNotPointerOrNil(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildArgIsNotPointerOrNil()

		expected := `
package example

import ()

// argIsNotPointerOrNil does what it says on the tin. This function
// is primarily useful for detecting if a destination value is valid
// before decoding an HTTP response, for instance.
func argIsNotPointerOrNil(i interface{}) error {
	if nn, err := argIsNotNil(i); nn || err != nil {
		return err
	}

	if np, err := argIsNotPointer(i); np || err != nil {
		return err
	}

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUnmarshalBody(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildUnmarshalBody(proj)

		expected := `
package example

import (
	"context"
	"encoding/json"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"io/ioutil"
	"net/http"
)

// unmarshalBody takes an HTTP response and JSON decodes its
// body into a destination value. ` + "`" + `dest` + "`" + ` must be a non-nil
// pointer to an object. Ideally, response is also not nil.
// The error returned here should only ever be received in
// testing, and should never be encountered by an end-user.
func unmarshalBody(ctx context.Context, res *http.Response, dest interface{}) error {
	_, span := tracing.StartSpan(ctx, "unmarshalBody")
	defer span.End()

	if err := argIsNotPointerOrNil(dest); err != nil {
		return err
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= http.StatusBadRequest {
		apiErr := &v1.ErrorResponse{}
		if err = json.Unmarshal(bodyBytes, &apiErr); err != nil {
			return fmt.Errorf("unmarshaling error: %w", err)
		}
		return apiErr
	}

	if err = json.Unmarshal(bodyBytes, &dest); err != nil {
		return fmt.Errorf("unmarshaling body: %w", err)
	}

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreateBodyFromStruct(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildCreateBodyFromStruct()

		expected := `
package example

import (
	"bytes"
	"encoding/json"
	"io"
)

// createBodyFromStruct takes any value in and returns an io.Reader
// for placement within http.NewRequest's last argument.
func createBodyFromStruct(in interface{}) (io.Reader, error) {
	out, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(out), nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
