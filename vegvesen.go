package vegvesen

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	basePath   string
	httpClient *http.Client
}

var (
	ErrNotFound = fmt.Errorf("do: resource not found")
)

func New(httpClient *http.Client) *Client {
	return &Client{
		basePath:   "https://www.vegvesen.no/ws/no/vegvesen/kjoretoy/kjoretoyoppslag/v1/kjennemerkeoppslag/kjoretoy/",
		httpClient: httpClient,
	}
}

func (c *Client) VehicleData(ctx context.Context, registrationNumber string) (VehicleData, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", c.basePath, registrationNumber))
	if err != nil {
		return VehicleData{}, err
	}

	var v VehicleData
	if err := c.do(req, &v); err != nil {
		return VehicleData{}, err
	}

	return v, nil
}

func (c *Client) do(r *http.Request, v interface{}) error {
	res, err := c.httpClient.Do(r)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(res.Header.Get("Content-Type"), "application/json") {
		return fmt.Errorf("do: unsupported content type: %q", res.Header.Get("Content-Type"))
	}

	if res.StatusCode != http.StatusOK {
		switch res.StatusCode {
		case http.StatusNotFound:
			return ErrNotFound
		default:
			return fmt.Errorf("do: non-OK http status: %d", res.StatusCode)
		}
	}

	if v == nil {
		return nil
	}

	return json.NewDecoder(res.Body).Decode(v)
}

func (c *Client) newRequest(ctx context.Context, method, endpoint string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Connection", "keep-alive")
	return req, nil
}
