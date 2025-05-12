package faceit

import (
	"fmt"
	"io"
	"net/http"
)

const baseUrl = "https://open.faceit.com/data/v4"

type Client struct {
	baseUrl    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		baseUrl:    baseUrl,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()

		return nil, fmt.Errorf("error from faceit api %s", string(body))
	}

	return resp, nil
}
