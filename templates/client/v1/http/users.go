package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

var usersBasePath = "users"

//
func (c *V1Client) BuildGetUserRequest(ctx context.Context, userID uint64) (*http.Request, error) {
	uri := c.buildVersionlessURL(nil, usersBasePath, strconv.FormatUint(userID, 10))
	return http.NewRequest(http.MethodGet, uri, nil)
}

//
func (c *V1Client) GetUser(ctx context.Context, userID uint64) (user *models.User, err error) {
	req, err := c.BuildGetUserRequest(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	err = c.retrieve(ctx, req, &user)
	return user, err
}

//
func (c *V1Client) BuildGetUsersRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	uri := c.buildVersionlessURL(filter.ToValues(), usersBasePath)
	return http.NewRequest(http.MethodGet, uri, nil)
}

//
func (c *V1Client) GetUsers(ctx context.Context, filter *models.QueryFilter) (*models.UserList, error) {
	users := &models.UserList{}
	req, err := c.BuildGetUsersRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	err = c.retrieve(ctx, req, &users)
	return users, err
}

//
func (c *V1Client) BuildCreateUserRequest(ctx context.Context, body *models.UserInput) (*http.Request, error) {
	uri := c.buildVersionlessURL(nil, usersBasePath)
	return c.buildDataRequest(http.MethodPost, uri, body)
}

//
func (c *V1Client) CreateUser(ctx context.Context, input *models.UserInput) (*models.UserCreationResponse, error) {
	user := &models.UserCreationResponse{}
	req, err := c.BuildCreateUserRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	err = c.executeUnathenticatedDataRequest(ctx, req, &user)
	return user, err
}

//
func (c *V1Client) BuildArchiveUserRequest(ctx context.Context, userID uint64) (*http.Request, error) {
	uri := c.buildVersionlessURL(nil, usersBasePath, strconv.FormatUint(userID, 10))
	return http.NewRequest(http.MethodDelete, uri, nil)
}

//
func (c *V1Client) ArchiveUser(ctx context.Context, userID uint64) error {
	req, err := c.BuildArchiveUserRequest(ctx, userID)
	if err != nil {
		return fmt.Errorf("building request", err)
	}
	return c.executeRequest(ctx, req, nil)
}

//
func (c *V1Client) BuildLoginRequest(username, password, totpToken string) (*http.Request, error) {
	body, err := createBodyFromStruct(&models.UserLoginInput{Username: username, Password: password, TOTPToken: totpToken})
	if err != nil {
		return nil, fmt.Errorf("creating body from struct", err)
	}
	uri := c.buildVersionlessURL(nil, usersBasePath, "login")
	return c.buildDataRequest(http.MethodPost, uri, body)
}

//
func (c *V1Client) Login(ctx context.Context, username, password, totpToken string) (*http.Cookie, error) {
	req, err := c.BuildLoginRequest(username, password, totpToken)
	if err != nil {
		return nil, err
	}
	res, err := c.plainClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("encountered error executing login request: %w", err)
	}
	if c.Debug {
		b, err := httputil.DumpResponse(res, true)
		if err != nil {
			c.logger.Error(err, "dumping response")
		}
		c.logger.WithValue("response", string(b)).Debug("login response received")
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			c.logger.Error(err, "closing response body")
		}
	}()
	cookies := res.Cookies()
	if len(cookies) > 0 {
		return cookies[0], nil
	}
	return nil, errors.New("no cookies returned from request")
}
