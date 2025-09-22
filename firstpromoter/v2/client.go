package firstpromoter

import (
	"net/http"
)

func New(accountId string, apiKey string, client *http.Client) *Client {
	handler := Client{
		accountId: accountId,
		apiKey:    apiKey,
		client:    client,
	}

	return &handler
}

type Client struct {
	accountId string
	apiKey    string
	client    *http.Client
}
