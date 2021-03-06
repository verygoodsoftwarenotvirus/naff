package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_oauth2ClientsDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := oauth2ClientsDotGo(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

var _ v1.OAuth2ClientDataManager = (*Client)(nil)

// GetOAuth2Client gets an OAuth2 client from the database.
func (c *Client) GetOAuth2Client(ctx context.Context, clientID, userID uint64) (*v1.OAuth2Client, error) {
	ctx, span := tracing.StartSpan(ctx, "GetOAuth2Client")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachOAuth2ClientDatabaseIDToSpan(span, clientID)

	logger := c.logger.WithValues(map[string]interface{}{
		"client_id": clientID,
		"user_id":   userID,
	})
	logger.Debug("GetOAuth2Client called")

	client, err := c.querier.GetOAuth2Client(ctx, clientID, userID)
	if err != nil {
		logger.Error(err, "error fetching oauth2 client from the querier")
		return nil, err
	}

	return client, nil
}

// GetOAuth2ClientByClientID fetches any OAuth2 client by client ID, regardless of ownership.
// This is used by authenticating middleware to fetch client information it needs to validate.
func (c *Client) GetOAuth2ClientByClientID(ctx context.Context, clientID string) (*v1.OAuth2Client, error) {
	ctx, span := tracing.StartSpan(ctx, "GetOAuth2ClientByClientID")
	defer span.End()

	tracing.AttachOAuth2ClientIDToSpan(span, clientID)
	logger := c.logger.WithValue("oauth2client_client_id", clientID)
	logger.Debug("GetOAuth2ClientByClientID called")

	client, err := c.querier.GetOAuth2ClientByClientID(ctx, clientID)
	if err != nil {
		logger.Error(err, "error fetching oauth2 client from the querier")
		return nil, err
	}

	return client, nil
}

// GetAllOAuth2ClientCount gets the count of OAuth2 clients that match the current filter.
func (c *Client) GetAllOAuth2ClientCount(ctx context.Context) (uint64, error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllOAuth2ClientCount")
	defer span.End()

	c.logger.Debug("GetAllOAuth2ClientCount called")

	return c.querier.GetAllOAuth2ClientCount(ctx)
}

// GetOAuth2ClientsForUser gets a list of OAuth2 clients.
func (c *Client) GetOAuth2ClientsForUser(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.OAuth2ClientList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetOAuth2ClientsForUser")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetOAuth2ClientsForUser called")

	return c.querier.GetOAuth2ClientsForUser(ctx, userID, filter)
}

// CreateOAuth2Client creates an OAuth2 client.
func (c *Client) CreateOAuth2Client(ctx context.Context, input *v1.OAuth2ClientCreationInput) (*v1.OAuth2Client, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateOAuth2Client")
	defer span.End()

	logger := c.logger.WithValues(map[string]interface{}{
		"oauth2client_client_id": input.ClientID,
		"belongs_to_user":        input.BelongsToUser,
	})

	client, err := c.querier.CreateOAuth2Client(ctx, input)
	if err != nil {
		logger.WithError(err).Debug("error writing oauth2 client to the querier")
		return nil, err
	}

	logger.WithValue("client_id", client.ID).Debug("new oauth2 client created successfully")

	return client, nil
}

// UpdateOAuth2Client updates a OAuth2 client. Note that this function expects the input's
// ID field to be valid.
func (c *Client) UpdateOAuth2Client(ctx context.Context, updated *v1.OAuth2Client) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateOAuth2Client")
	defer span.End()

	return c.querier.UpdateOAuth2Client(ctx, updated)
}

// ArchiveOAuth2Client archives an OAuth2 client.
func (c *Client) ArchiveOAuth2Client(ctx context.Context, clientID, userID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveOAuth2Client")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachOAuth2ClientDatabaseIDToSpan(span, clientID)

	logger := c.logger.WithValues(map[string]interface{}{
		"client_id":       clientID,
		"belongs_to_user": userID,
	})

	err := c.querier.ArchiveOAuth2Client(ctx, clientID, userID)
	if err != nil {
		logger.WithError(err).Debug("error deleting oauth2 client to the querier")
		return err
	}
	logger.Debug("removed oauth2 client successfully")

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetOAuth2Client(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetOAuth2Client gets an OAuth2 client from the database.
func (c *Client) GetOAuth2Client(ctx context.Context, clientID, userID uint64) (*v1.OAuth2Client, error) {
	ctx, span := tracing.StartSpan(ctx, "GetOAuth2Client")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachOAuth2ClientDatabaseIDToSpan(span, clientID)

	logger := c.logger.WithValues(map[string]interface{}{
		"client_id": clientID,
		"user_id":   userID,
	})
	logger.Debug("GetOAuth2Client called")

	client, err := c.querier.GetOAuth2Client(ctx, clientID, userID)
	if err != nil {
		logger.Error(err, "error fetching oauth2 client from the querier")
		return nil, err
	}

	return client, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetOAuth2ClientByClientID(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetOAuth2ClientByClientID(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetOAuth2ClientByClientID fetches any OAuth2 client by client ID, regardless of ownership.
// This is used by authenticating middleware to fetch client information it needs to validate.
func (c *Client) GetOAuth2ClientByClientID(ctx context.Context, clientID string) (*v1.OAuth2Client, error) {
	ctx, span := tracing.StartSpan(ctx, "GetOAuth2ClientByClientID")
	defer span.End()

	tracing.AttachOAuth2ClientIDToSpan(span, clientID)
	logger := c.logger.WithValue("oauth2client_client_id", clientID)
	logger.Debug("GetOAuth2ClientByClientID called")

	client, err := c.querier.GetOAuth2ClientByClientID(ctx, clientID)
	if err != nil {
		logger.Error(err, "error fetching oauth2 client from the querier")
		return nil, err
	}

	return client, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetAllOAuth2ClientCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetAllOAuth2ClientCount(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// GetAllOAuth2ClientCount gets the count of OAuth2 clients that match the current filter.
func (c *Client) GetAllOAuth2ClientCount(ctx context.Context) (uint64, error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllOAuth2ClientCount")
	defer span.End()

	c.logger.Debug("GetAllOAuth2ClientCount called")

	return c.querier.GetAllOAuth2ClientCount(ctx)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetOAuth2ClientsForUser(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetOAuth2ClientsForUser gets a list of OAuth2 clients.
func (c *Client) GetOAuth2ClientsForUser(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.OAuth2ClientList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetOAuth2ClientsForUser")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetOAuth2ClientsForUser called")

	return c.querier.GetOAuth2ClientsForUser(ctx, userID, filter)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildCreateOAuth2Client(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// CreateOAuth2Client creates an OAuth2 client.
func (c *Client) CreateOAuth2Client(ctx context.Context, input *v1.OAuth2ClientCreationInput) (*v1.OAuth2Client, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateOAuth2Client")
	defer span.End()

	logger := c.logger.WithValues(map[string]interface{}{
		"oauth2client_client_id": input.ClientID,
		"belongs_to_user":        input.BelongsToUser,
	})

	client, err := c.querier.CreateOAuth2Client(ctx, input)
	if err != nil {
		logger.WithError(err).Debug("error writing oauth2 client to the querier")
		return nil, err
	}

	logger.WithValue("client_id", client.ID).Debug("new oauth2 client created successfully")

	return client, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUpdateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildUpdateOAuth2Client(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// UpdateOAuth2Client updates a OAuth2 client. Note that this function expects the input's
// ID field to be valid.
func (c *Client) UpdateOAuth2Client(ctx context.Context, updated *v1.OAuth2Client) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateOAuth2Client")
	defer span.End()

	return c.querier.UpdateOAuth2Client(ctx, updated)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildArchiveOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildArchiveOAuth2Client(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// ArchiveOAuth2Client archives an OAuth2 client.
func (c *Client) ArchiveOAuth2Client(ctx context.Context, clientID, userID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveOAuth2Client")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachOAuth2ClientDatabaseIDToSpan(span, clientID)

	logger := c.logger.WithValues(map[string]interface{}{
		"client_id":       clientID,
		"belongs_to_user": userID,
	})

	err := c.querier.ArchiveOAuth2Client(ctx, clientID, userID)
	if err != nil {
		logger.WithError(err).Debug("error deleting oauth2 client to the querier")
		return err
	}
	logger.Debug("removed oauth2 client successfully")

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
