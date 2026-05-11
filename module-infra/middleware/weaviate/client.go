package weaviate

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/data/replication"
)

type Client struct {
	client *weaviate.Client
}

func NewClient(scheme, host string) (*Client, error) {
	cfg := weaviate.Config{
		Scheme: scheme,
		Host:   host,
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Weaviate client: %w", err)
	}
	return &Client{client: client}, nil
}

func (c *Client) CreateObject(ctx context.Context, className string, properties map[string]interface{}) error {
	_, err := c.client.Data().Creator().
		WithClassName(className).
		WithProperties(properties).
		WithConsistencyLevel(replication.ConsistencyLevel.ONE).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to create object: %w", err)
	}
	return nil
}

func (c *Client) GetObject(ctx context.Context, className, id string) (map[string]interface{}, error) {
	result, err := c.client.Data().ObjectsGetter().
		WithClassName(className).
		WithID(id).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("object not found: %s/%s", className, id)
	}
	props, ok := result[0].Properties.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected properties type")
	}
	return props, nil
}
