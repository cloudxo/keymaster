package workflow

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL *url.URL
	UserAgent string
	HttpClient *http.Client
	PolicyEncrypter string // TODO: key, signingmethod etc
}

type StartRequest struct {
	Requester struct {
		Name string `json:"name"`
		Username string `json:"username"`
		Email string `json:"email"`
	} `json:"requester"`
	Source struct {
		// Source Key? For gitlab?
		Description string `json:"description"`
		DetailsURI string `json:"details_uri"`
	} `json:"source"`
	Target struct {
		Name string `json:"name"`
		// Target discovery URI?
	} `json:"target"`
	// To be encrypted with workflow public key
	WorkflowPolicy string `json:"workflow_policy"`
}

type StartResponse struct {
	// Workflow id
	// Workflow nonce
}

type GetDetailsRequest struct {
	// Workflow id
}

type GetDetailsResponse struct {
	// State: pending / completed
}

type GetAssertionsRequest struct {
	// Workflow id
	// Workflow nonce
}

type GetAssertionsResponse struct {
	// Bag of SAML assertions. Could be wrapped(?)
}

func (c *Client) Start(ctx context.Context, req *StartRequest) (*StartResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := c.newRequest(ctx, "POST", "/1/workflow/create", body)
	if err != nil {
		return nil, err
	}
	var resp StartResponse
	_, err = c.do(httpReq, &resp)
	return &resp, err
}

func (c *Client) GetDetails(ctx context.Context, req *GetDetailsRequest) (*GetDetailsResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := c.newRequest(ctx, "POST", "/1/workflow/getDetails", body)
	if err != nil {
		return nil, err
	}
	var resp GetDetailsResponse
	_, err = c.do(httpReq, &resp)
	return &resp, err
}

func (c *Client) GetAssertions(ctx context.Context, req *GetAssertionsRequest) (*GetAssertionsResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := c.newRequest(ctx, "POST", "/1/workflow/getAssertions", body)
	if err != nil {
		return nil, err
	}
	var resp GetAssertionsResponse
	_, err = c.do(httpReq, &resp)
	return &resp, err
}

func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}
