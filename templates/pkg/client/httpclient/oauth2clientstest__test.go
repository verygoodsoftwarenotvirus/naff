package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_oauth2ClientsTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := oauth2ClientsTestDotGo(proj)

		expected := `
package example

import (
	"context"
	"encoding/json"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestV1Client_BuildGetOAuth2ClientRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		actual, err := c.BuildGetOAuth2ClientRequest(ctx, exampleOAuth2Client.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleOAuth2Client.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleOAuth2Client.ID))))
					assert.Equal(t, fmt.Sprintf("/api/v1/oauth2/clients/%d", exampleOAuth2Client.ID), req.URL.Path, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleOAuth2Client))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetOAuth2Client(ctx, exampleOAuth2Client.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleOAuth2Client, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetOAuth2Client(ctx, exampleOAuth2Client.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildGetOAuth2ClientsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildGetOAuth2ClientsRequest(ctx, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetOAuth2Clients(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, "/api/v1/oauth2/clients", "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleOAuth2ClientList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetOAuth2Clients(ctx, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleOAuth2ClientList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetOAuth2Clients(ctx, nil)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildCreateOAuth2ClientRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)
		req, err := c.BuildCreateOAuth2ClientRequest(ctx, &http.Cookie{}, exampleInput)

		require.NotNil(t, req)
		assert.NoError(t, err)
		assert.Equal(t, http.MethodPost, req.Method)
	})
}

func TestV1Client_CreateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, "/oauth2/client", req.URL.Path, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)
					require.NoError(t, json.NewEncoder(res).Encode(exampleOAuth2Client))
				},
			),
		)
		c := buildTestClient(t, ts)

		actual, err := c.CreateOAuth2Client(ctx, &http.Cookie{}, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateOAuth2Client(ctx, &http.Cookie{}, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response from server", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, "/oauth2/client", "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)
					_, err := res.Write([]byte("BLAH"))
					assert.NoError(t, err)
				},
			),
		)
		c := buildTestClient(t, ts)

		actual, err := c.CreateOAuth2Client(ctx, &http.Cookie{}, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("without cookie", func(t *testing.T) {
		ctx := context.Background()

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)

		_, err := c.CreateOAuth2Client(ctx, nil, nil)
		assert.Error(t, err)
	})
}

func TestV1Client_BuildArchiveOAuth2ClientRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodDelete
		ts := httptest.NewTLSServer(nil)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		c := buildTestClient(t, ts)
		actual, err := c.BuildArchiveOAuth2ClientRequest(ctx, exampleOAuth2Client.ID)

		require.NotNil(t, actual)
		require.NotNil(t, actual.URL)
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleOAuth2Client.ID)))
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_ArchiveOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/oauth2/clients/%d", exampleOAuth2Client.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodDelete)

					res.WriteHeader(http.StatusOK)
				},
			),
		)

		err := buildTestClient(t, ts).ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		err := buildTestClientWithInvalidURL(t).ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildV1Client_BuildGetOAuth2ClientRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildV1Client_BuildGetOAuth2ClientRequest(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestV1Client_BuildGetOAuth2ClientRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		actual, err := c.BuildGetOAuth2ClientRequest(ctx, exampleOAuth2Client.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleOAuth2Client.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildV1Client_GetOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildV1Client_GetOAuth2Client(proj)

		expected := `
package example

import (
	"context"
	"encoding/json"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestV1Client_GetOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleOAuth2Client.ID))))
					assert.Equal(t, fmt.Sprintf("/api/v1/oauth2/clients/%d", exampleOAuth2Client.ID), req.URL.Path, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleOAuth2Client))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetOAuth2Client(ctx, exampleOAuth2Client.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleOAuth2Client, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetOAuth2Client(ctx, exampleOAuth2Client.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildV1Client_BuildGetOAuth2ClientsRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildV1Client_BuildGetOAuth2ClientsRequest()

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_BuildGetOAuth2ClientsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildGetOAuth2ClientsRequest(ctx, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildV1Client_GetOAuth2Clients(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildV1Client_GetOAuth2Clients(proj)

		expected := `
package example

import (
	"context"
	"encoding/json"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_GetOAuth2Clients(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, "/api/v1/oauth2/clients", "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleOAuth2ClientList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetOAuth2Clients(ctx, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleOAuth2ClientList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetOAuth2Clients(ctx, nil)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildV1Client_BuildCreateOAuth2ClientRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildV1Client_BuildCreateOAuth2ClientRequest(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_BuildCreateOAuth2ClientRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)
		req, err := c.BuildCreateOAuth2ClientRequest(ctx, &http.Cookie{}, exampleInput)

		require.NotNil(t, req)
		assert.NoError(t, err)
		assert.Equal(t, http.MethodPost, req.Method)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildV1Client_CreateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildV1Client_CreateOAuth2Client(proj)

		expected := `
package example

import (
	"context"
	"encoding/json"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_CreateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, "/oauth2/client", req.URL.Path, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)
					require.NoError(t, json.NewEncoder(res).Encode(exampleOAuth2Client))
				},
			),
		)
		c := buildTestClient(t, ts)

		actual, err := c.CreateOAuth2Client(ctx, &http.Cookie{}, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateOAuth2Client(ctx, &http.Cookie{}, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response from server", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, "/oauth2/client", "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)
					_, err := res.Write([]byte("BLAH"))
					assert.NoError(t, err)
				},
			),
		)
		c := buildTestClient(t, ts)

		actual, err := c.CreateOAuth2Client(ctx, &http.Cookie{}, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("without cookie", func(t *testing.T) {
		ctx := context.Background()

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)

		_, err := c.CreateOAuth2Client(ctx, nil, nil)
		assert.Error(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildV1Client_BuildArchiveOAuth2ClientRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildV1Client_BuildArchiveOAuth2ClientRequest(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestV1Client_BuildArchiveOAuth2ClientRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodDelete
		ts := httptest.NewTLSServer(nil)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		c := buildTestClient(t, ts)
		actual, err := c.BuildArchiveOAuth2ClientRequest(ctx, exampleOAuth2Client.ID)

		require.NotNil(t, actual)
		require.NotNil(t, actual.URL)
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleOAuth2Client.ID)))
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildV1Client_ArchiveOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildV1Client_ArchiveOAuth2Client(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_ArchiveOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/oauth2/clients/%d", exampleOAuth2Client.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodDelete)

					res.WriteHeader(http.StatusOK)
				},
			),
		)

		err := buildTestClient(t, ts).ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		err := buildTestClientWithInvalidURL(t).ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
