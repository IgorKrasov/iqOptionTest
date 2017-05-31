package cacheclient

import (
	"encoding/json"
	"golang.org/x/net/context"
	"strconv"
)

type SetBody struct {
	Key     string
	Expired int
	Value   interface{}
}

func (c *Client) Set(ctx context.Context, body *SetBody) (map[string]string, error) {
	b, _ := json.Marshal(body)
	config := &apiConfig{
		path: "set",
	}
	response := map[string]string{}
	err := c.postJSON(ctx, config, b, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) Get(ctx context.Context, key string) (interface{}, error) {
	config := &apiConfig{
		path: "get/" + key,
	}
	var response interface{}
	err := c.getJSON(ctx, config, &response)

	if err != nil {
		return nil, err
	}

	return response, err
}

func (c *Client) Keys(ctx context.Context) ([]string, error) {
	config := &apiConfig{
		path: "keys",
	}
	var response []string
	err := c.getJSON(ctx, config, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) Unset(ctx context.Context, key string) (map[string]string, error) {
	config := &apiConfig{
		path: "unset/" + key,
	}
	var response map[string]string
	err := c.deleteJSON(ctx, config, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) Rpush(ctx context.Context, body *SetBody) (map[string]string, error) {
	b, _ := json.Marshal(body)
	config := &apiConfig{
		path: "rpush",
	}
	var response map[string]string
	err := c.postJSON(ctx, config, b, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) Pop(ctx context.Context, key string) (interface{}, error) {
	config := &apiConfig{
		path: "pop/" + key,
	}
	var response interface{}
	err := c.getJSON(ctx, config, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) Lget(ctx context.Context, key string, i int) (interface{}, error) {
	config := &apiConfig{
		path: "lget/" + key + "/" + strconv.Itoa(i),
	}
	var response interface{}
	err := c.getJSON(ctx, config, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) Lgetall(ctx context.Context, key string) ([]interface{}, error) {
	config := &apiConfig{
		path: "lgetall/" + key,
	}
	var response []interface{}
	err := c.getJSON(ctx, config, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) Hset(ctx context.Context, body *SetBody) (map[string]string, error) {
	b, _ := json.Marshal(body)
	config := &apiConfig{
		path: "hset",
	}
	var response map[string]string
	err := c.postJSON(ctx, config, b, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) Hgetall(ctx context.Context, key string) (map[string]interface{}, error) {
	config := &apiConfig{
		path: "hgetall/" + key,
	}
	var response map[string]interface{}
	err := c.getJSON(ctx, config, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) Hget(ctx context.Context, key string, i string) (interface{}, error) {
	config := &apiConfig{
		path: "hget/" + key + "/" + i,
	}
	var response interface{}
	err := c.getJSON(ctx, config, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}
