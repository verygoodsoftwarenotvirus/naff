package fake

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_webhookDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := webhookDotGo(proj)

		expected := `
package example

import (
	v5 "github.com/brianvoe/gofakeit/v5"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// BuildFakeWebhook builds a faked Webhook.
func BuildFakeWebhook() *v1.Webhook {
	return &v1.Webhook{
		ID:            v5.Uint64(),
		Name:          v5.Word(),
		ContentType:   v5.FileMimeType(),
		URL:           v5.URL(),
		Method:        v5.HTTPMethod(),
		Events:        []string{v5.Word()},
		DataTypes:     []string{v5.Word()},
		Topics:        []string{v5.Word()},
		CreatedOn:     uint64(uint32(v5.Date().Unix())),
		ArchivedOn:    nil,
		BelongsToUser: v5.Uint64(),
	}
}

// BuildFakeWebhookList builds a faked WebhookList.
func BuildFakeWebhookList() *v1.WebhookList {
	exampleWebhook1 := BuildFakeWebhook()
	exampleWebhook2 := BuildFakeWebhook()
	exampleWebhook3 := BuildFakeWebhook()
	return &v1.WebhookList{
		Pagination: v1.Pagination{
			Page:  1,
			Limit: 20,
		},
		Webhooks: []v1.Webhook{
			*exampleWebhook1,
			*exampleWebhook2,
			*exampleWebhook3,
		},
	}
}

// BuildFakeWebhookUpdateInputFromWebhook builds a faked WebhookUpdateInput.
func BuildFakeWebhookUpdateInputFromWebhook(webhook *v1.Webhook) *v1.WebhookUpdateInput {
	return &v1.WebhookUpdateInput{
		Name:          webhook.Name,
		ContentType:   webhook.ContentType,
		URL:           webhook.URL,
		Method:        webhook.Method,
		Events:        webhook.Events,
		DataTypes:     webhook.DataTypes,
		Topics:        webhook.Topics,
		BelongsToUser: webhook.BelongsToUser,
	}
}

// BuildFakeWebhookCreationInput builds a faked WebhookCreationInput.
func BuildFakeWebhookCreationInput() *v1.WebhookCreationInput {
	webhook := BuildFakeWebhook()
	return BuildFakeWebhookCreationInputFromWebhook(webhook)
}

// BuildFakeWebhookCreationInputFromWebhook builds a faked WebhookCreationInput.
func BuildFakeWebhookCreationInputFromWebhook(webhook *v1.Webhook) *v1.WebhookCreationInput {
	return &v1.WebhookCreationInput{
		Name:          webhook.Name,
		ContentType:   webhook.ContentType,
		URL:           webhook.URL,
		Method:        webhook.Method,
		Events:        webhook.Events,
		DataTypes:     webhook.DataTypes,
		Topics:        webhook.Topics,
		BelongsToUser: webhook.BelongsToUser,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildFakeWebhook(proj)

		expected := `
package example

import (
	v5 "github.com/brianvoe/gofakeit/v5"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// BuildFakeWebhook builds a faked Webhook.
func BuildFakeWebhook() *v1.Webhook {
	return &v1.Webhook{
		ID:            v5.Uint64(),
		Name:          v5.Word(),
		ContentType:   v5.FileMimeType(),
		URL:           v5.URL(),
		Method:        v5.HTTPMethod(),
		Events:        []string{v5.Word()},
		DataTypes:     []string{v5.Word()},
		Topics:        []string{v5.Word()},
		CreatedOn:     uint64(uint32(v5.Date().Unix())),
		ArchivedOn:    nil,
		BelongsToUser: v5.Uint64(),
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeWebhookList(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildFakeWebhookList(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// BuildFakeWebhookList builds a faked WebhookList.
func BuildFakeWebhookList() *v1.WebhookList {
	exampleWebhook1 := BuildFakeWebhook()
	exampleWebhook2 := BuildFakeWebhook()
	exampleWebhook3 := BuildFakeWebhook()
	return &v1.WebhookList{
		Pagination: v1.Pagination{
			Page:  1,
			Limit: 20,
		},
		Webhooks: []v1.Webhook{
			*exampleWebhook1,
			*exampleWebhook2,
			*exampleWebhook3,
		},
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeWebhookUpdateInputFromWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildFakeWebhookUpdateInputFromWebhook(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// BuildFakeWebhookUpdateInputFromWebhook builds a faked WebhookUpdateInput.
func BuildFakeWebhookUpdateInputFromWebhook(webhook *v1.Webhook) *v1.WebhookUpdateInput {
	return &v1.WebhookUpdateInput{
		Name:          webhook.Name,
		ContentType:   webhook.ContentType,
		URL:           webhook.URL,
		Method:        webhook.Method,
		Events:        webhook.Events,
		DataTypes:     webhook.DataTypes,
		Topics:        webhook.Topics,
		BelongsToUser: webhook.BelongsToUser,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeWebhookCreationInput(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildFakeWebhookCreationInput(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// BuildFakeWebhookCreationInput builds a faked WebhookCreationInput.
func BuildFakeWebhookCreationInput() *v1.WebhookCreationInput {
	webhook := BuildFakeWebhook()
	return BuildFakeWebhookCreationInputFromWebhook(webhook)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeWebhookCreationInputFromWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildFakeWebhookCreationInputFromWebhook(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// BuildFakeWebhookCreationInputFromWebhook builds a faked WebhookCreationInput.
func BuildFakeWebhookCreationInputFromWebhook(webhook *v1.Webhook) *v1.WebhookCreationInput {
	return &v1.WebhookCreationInput{
		Name:          webhook.Name,
		ContentType:   webhook.ContentType,
		URL:           webhook.URL,
		Method:        webhook.Method,
		Events:        webhook.Events,
		DataTypes:     webhook.DataTypes,
		Topics:        webhook.Topics,
		BelongsToUser: webhook.BelongsToUser,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
