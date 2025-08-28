package telegram

import (
	"PDFHub/internal/bot/clients"
	"PDFHub/internal/bot/lib/e"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
	getFileMethod     = "getFile"
	sendDocument      = "sendDocument"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: basePath(token),
		client:   http.Client{},
	}
}

func (c *Client) File(FileID string) (*http.Response, error) {
	q := url.Values{}
	q.Set("file_id", FileID)

	data, err := c.doRequest(getFileMethod, q)
	if err != nil {
		return nil, e.Wrap("can't get file", err)
	}

	var res clients.Document

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, e.Wrap("can't get file", err)
	}
	resp, err := http.Get(path.Join(c.host, "file", c.basePath, res.FilePath))
	if err != nil {
		return nil, e.Wrap("can't get file", err)
	}
	return resp, nil
}

func (c *Client) FullFilePath(FileID string) (string, error) {
	q := url.Values{}
	q.Set("file_id", FileID)
	data, err := c.doRequest(getFileMethod, q)
	if err != nil {
		return "", e.Wrap("can't get file", err)
	}

	var res clients.Document

	if err := json.Unmarshal(data, &res); err != nil {
		return "", e.Wrap("can't get file", err)
	}

	return path.Join(c.host, "file", c.basePath, res.FilePath), nil
}

func (c *Client) Updates(offset int, limit int) ([]clients.Update, error) {
	q := url.Values{}
	q.Set("offset", strconv.Itoa(offset))
	q.Set("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, e.Wrap("can't get updates", err)
	}

	var res clients.UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, e.Wrap("can't get updates", err)
	}

	return res.Result, nil
}

func basePath(token string) string {
	return "bot" + token
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Set("chat_id", strconv.Itoa(chatID))
	q.Set("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) doRequest(method string, q url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: c.host,
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap("can't create request", err)
	}

	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap("can't do request", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap("can't read response body", err)
	}

	return body, nil
}
