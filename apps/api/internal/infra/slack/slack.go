package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"gogym-api/internal/configs"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	contactURL string
	errorsURL  string
}

func NewClient(sc configs.SlackConfig) (*Client, error) {
	return &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		contactURL: sc.ContactWebhookURL,
		errorsURL:  sc.ErrorsWebhookURL,
	}, nil
}

type payload struct {
	Text string `json:"text"`
}

func (c *Client) NotifyContact(ctx context.Context, text string) error {
	return c.postJSON(ctx, c.contactURL, payload{Text: text})
}

func (c *Client) NotifyError(ctx context.Context, text string) error {
	return c.postJSON(ctx, c.errorsURL, payload{Text: text})
}

func (c *Client) postJSON(ctx context.Context, url string, body any) error {
	if url == "" {
		return errors.New("slack webhook url is empty")
	}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return errors.New("slack webhook request failed")
	}
	return nil
}
