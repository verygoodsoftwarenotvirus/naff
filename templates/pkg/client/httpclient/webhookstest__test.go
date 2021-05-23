package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_webhooksTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := webhooksTestDotGo(proj)

		expected := `
package example

import (
	"context"
	"encoding/json"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestV1Client_BuildGetWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet
		exampleWebhook := fake.BuildFakeWebhook()

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildGetWebhookRequest(ctx, exampleWebhook.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleWebhook.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetWebhook(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleWebhook.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/webhooks/%d", exampleWebhook.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleWebhook))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetWebhook(ctx, exampleWebhook.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleWebhook, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()

		actual, err := buildTestClientWithInvalidURL(t).GetWebhook(ctx, exampleWebhook.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildGetWebhooksRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetWebhooksRequest(ctx, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetWebhooks(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhookList := fake.BuildFakeWebhookList()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, "/api/v1/webhooks", "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleWebhookList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetWebhooks(ctx, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleWebhookList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		actual, err := buildTestClientWithInvalidURL(t).GetWebhooks(ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildCreateWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodPost
		ts := httptest.NewTLSServer(nil)

		exampleWebhook := fake.BuildFakeWebhook()
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)

		c := buildTestClient(t, ts)
		actual, err := c.BuildCreateWebhookRequest(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_CreateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
		exampleInput.BelongsToUser = 0

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, "/api/v1/webhooks", "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					var x *v1.WebhookCreationInput
					require.NoError(t, json.NewDecoder(req.Body).Decode(&x))
					assert.Equal(t, exampleInput, x)

					require.NoError(t, json.NewEncoder(res).Encode(exampleWebhook))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.CreateWebhook(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleWebhook, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)

		actual, err := buildTestClientWithInvalidURL(t).CreateWebhook(ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildUpdateWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodPut
		exampleWebhook := fake.BuildFakeWebhook()

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildUpdateWebhookRequest(ctx, exampleWebhook)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_UpdateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/webhooks/%d", exampleWebhook.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPut)
					assert.NoError(t, json.NewEncoder(res).Encode(exampleWebhook))
				},
			),
		)

		err := buildTestClient(t, ts).UpdateWebhook(ctx, exampleWebhook)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()

		err := buildTestClientWithInvalidURL(t).UpdateWebhook(ctx, exampleWebhook)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildArchiveWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodDelete
		ts := httptest.NewTLSServer(nil)
		exampleWebhook := fake.BuildFakeWebhook()

		c := buildTestClient(t, ts)
		actual, err := c.BuildArchiveWebhookRequest(ctx, exampleWebhook.ID)

		require.NotNil(t, actual)
		require.NotNil(t, actual.URL)
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleWebhook.ID)))
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_ArchiveWebhook(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/webhooks/%d", exampleWebhook.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodDelete)
				},
			),
		)

		err := buildTestClient(t, ts).ArchiveWebhook(ctx, exampleWebhook.ID)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()

		err := buildTestClientWithInvalidURL(t).ArchiveWebhook(ctx, exampleWebhook.ID)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientBuildGetWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientBuildGetWebhookRequest(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestV1Client_BuildGetWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet
		exampleWebhook := fake.BuildFakeWebhook()

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildGetWebhookRequest(ctx, exampleWebhook.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleWebhook.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientGetWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientGetWebhook(proj)

		expected := `
package example

import (
	"context"
	"encoding/json"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestV1Client_GetWebhook(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleWebhook.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/webhooks/%d", exampleWebhook.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleWebhook))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetWebhook(ctx, exampleWebhook.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleWebhook, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()

		actual, err := buildTestClientWithInvalidURL(t).GetWebhook(ctx, exampleWebhook.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientBuildGetWebhooksRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientBuildGetWebhooksRequest(proj)

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

func TestV1Client_BuildGetWebhooksRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetWebhooksRequest(ctx, nil)

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

func Test_buildTestV1ClientGetWebhooks(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientGetWebhooks(proj)

		expected := `
package example

import (
	"context"
	"encoding/json"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_GetWebhooks(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhookList := fake.BuildFakeWebhookList()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, "/api/v1/webhooks", "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleWebhookList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetWebhooks(ctx, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleWebhookList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		actual, err := buildTestClientWithInvalidURL(t).GetWebhooks(ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientBuildCreateWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientBuildCreateWebhookRequest(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_BuildCreateWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodPost
		ts := httptest.NewTLSServer(nil)

		exampleWebhook := fake.BuildFakeWebhook()
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)

		c := buildTestClient(t, ts)
		actual, err := c.BuildCreateWebhookRequest(ctx, exampleInput)

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

func Test_buildTestV1ClientCreateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientCreateWebhook(proj)

		expected := `
package example

import (
	"context"
	"encoding/json"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_CreateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
		exampleInput.BelongsToUser = 0

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, "/api/v1/webhooks", "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					var x *v1.WebhookCreationInput
					require.NoError(t, json.NewDecoder(req.Body).Decode(&x))
					assert.Equal(t, exampleInput, x)

					require.NoError(t, json.NewEncoder(res).Encode(exampleWebhook))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.CreateWebhook(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleWebhook, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)

		actual, err := buildTestClientWithInvalidURL(t).CreateWebhook(ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientBuildUpdateWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientBuildUpdateWebhookRequest(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_BuildUpdateWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodPut
		exampleWebhook := fake.BuildFakeWebhook()

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildUpdateWebhookRequest(ctx, exampleWebhook)

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

func Test_buildTestV1ClientUpdateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientUpdateWebhook(proj)

		expected := `
package example

import (
	"context"
	"encoding/json"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_UpdateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/webhooks/%d", exampleWebhook.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPut)
					assert.NoError(t, json.NewEncoder(res).Encode(exampleWebhook))
				},
			),
		)

		err := buildTestClient(t, ts).UpdateWebhook(ctx, exampleWebhook)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()

		err := buildTestClientWithInvalidURL(t).UpdateWebhook(ctx, exampleWebhook)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientBuildArchiveWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientBuildArchiveWebhookRequest(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestV1Client_BuildArchiveWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodDelete
		ts := httptest.NewTLSServer(nil)
		exampleWebhook := fake.BuildFakeWebhook()

		c := buildTestClient(t, ts)
		actual, err := c.BuildArchiveWebhookRequest(ctx, exampleWebhook.ID)

		require.NotNil(t, actual)
		require.NotNil(t, actual.URL)
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleWebhook.ID)))
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientArchiveWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientArchiveWebhook(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_ArchiveWebhook(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/webhooks/%d", exampleWebhook.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodDelete)
				},
			),
		)

		err := buildTestClient(t, ts).ArchiveWebhook(ctx, exampleWebhook.ID)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()

		err := buildTestClientWithInvalidURL(t).ArchiveWebhook(ctx, exampleWebhook.ID)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
