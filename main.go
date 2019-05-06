package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Client struct
type Client struct {
	token string
	addr  string
}

// Option of client
type Option struct {
	Token string
	Addr  string
}

// New func create a client.
func New(opt Option) *Client {
	return &Client{
		token: opt.Token,
		addr:  opt.Addr,
	}
}

// Addr returns client addr.
func (c *Client) Addr() string {
	return c.addr
}

// SetAuth sets token for client.
func (c *Client) SetAuth(token string) error {
	if c == nil {
		return fmt.Errorf("client not init")
	}
	c.token = token
	return nil
}

type namespace struct {
	StatusCode int      `json:"httpstatus"`
	Msg        string   `json:"msg"`
	Data       []string `json:"data"`
}

// Namespace returns all nodes under given ns.
func (c *Client) Namespace(ns string, onlyLeaf bool) ([]string, error) {
	var res []string
	fullURL := fmt.Sprintf("%s/router/ns?ns=%s&format=list", c.Addr, ns)
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("HTTP status error: %d", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data namespace
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	res = data.Data

	if !onlyLeaf {
		tmp := make(map[string]bool)
		for _, leaf := range data.Data {
			arr := strings.SplitAfterN(leaf, ".", 2)
			if len(arr) > 1 {
				tmp[arr[1]] = true
			}
		}

		for k := range tmp {
			res = append(res, k)
		}
	}

	return res, nil
}
