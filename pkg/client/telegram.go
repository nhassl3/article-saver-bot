package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/nhassl3/article-saver-bot/pkg/e"
	"github.com/nhassl3/article-saver-bot/pkg/entities"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func NewClient(host, token string) *Client {
	return &Client{
		host:     host,
		basePath: "bot" + token,
		client:   http.Client{},
	}
}

func (c *Client) Updates(offset, limit int) ([]entities.Updates, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset)) // Integer to ASCII
	q.Add("limit", strconv.Itoa(limit))

	// Request
	data, err := c.DoRequest(q, getUpdatesMethod)
	if err != nil {
		return nil, err
	}

	var res entities.UpdatesResponse
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatId int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.DoRequest(q, sendMessageMethod)
	return e.WrapIfErr("can't send message", err)
}

func (c *Client) DoRequest(query url.Values, method string) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
