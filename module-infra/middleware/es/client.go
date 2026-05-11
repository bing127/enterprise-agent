package es

import (
	"context"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

type Client struct {
	es *elasticsearch.Client
}

func NewClient(address string) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{address},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{es: es}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	res, err := c.es.Ping()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch ping failed: %s", res.Status())
	}
	return nil
}
