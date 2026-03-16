// Code generated from jsonrpc schema by rpcgen v2.5.x with golang v1.1.1; DO NOT EDIT.

package newsportal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/vmkteam/appkit"
	"github.com/vmkteam/zenrpc/v2"
)

const name = "news-newsportal"

var (
	// Always import time package. Generated models can contain time.Time fields.
	_ time.Time
)

type Client struct {
	rpcClient *rpcClient

	News *svcNews
}

func NewClient(endpoint string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: time.Second * 30}
	}
	c := &Client{
		rpcClient: newRPCClient(endpoint, httpClient),
	}

	c.News = newClientNews(c.rpcClient)

	return c
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type News struct {
	Author       string    `json:"author"`
	Category     *Category `json:"category,omitempty"`
	Content      *string   `json:"content,omitempty"`
	Created_at   string    `json:"created_at"`
	ID           int       `json:"id"`
	Preamble     string    `json:"preamble"`
	Published_at string    `json:"published_at"`
	Tags         []Tag     `json:"tags"`
	Title        string    `json:"title"`
}

type NewsFilter struct {
	// Идентификатор категории
	Category_id int `json:"category_id"`
	// Начало периода
	From string `json:"from"`
	// Идентификатор тега
	Tag_id int `json:"tag_id"`
	// Конец периода
	To string `json:"to"`
}

type Pager struct {
	// Количество на страницу
	Limit int `json:"limit"`
	// Номер страницы
	Page int `json:"page"`
}

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type svcNews struct {
	client *rpcClient
}

func newClientNews(client *rpcClient) *svcNews {
	return &svcNews{
		client: client,
	}
}

func (c *svcNews) Categories(ctx context.Context) (res []Category, err error) {
	_req := struct {
	}{}

	err = c.client.call(ctx, "news.Categories", _req, &res)

	return
}

func (c *svcNews) Count(ctx context.Context, filter NewsFilter) (res int, err error) {
	_req := struct {
		Filter NewsFilter `json:"filter,omitempty"`
	}{
		Filter: filter,
	}

	err = c.client.call(ctx, "news.Count", _req, &res)

	return
}

func (c *svcNews) Get(ctx context.Context, id int) (res *News, err error) {
	_req := struct {
		ID int `json:"id,omitempty"`
	}{
		ID: id,
	}

	err = c.client.call(ctx, "news.Get", _req, &res)

	return
}

func (c *svcNews) List(ctx context.Context, filter NewsFilter, pager Pager) (res []News, err error) {
	_req := struct {
		Filter NewsFilter `json:"filter,omitempty"`
		Pager  Pager      `json:"pager,omitempty"`
	}{
		Filter: filter, Pager: pager,
	}

	err = c.client.call(ctx, "news.List", _req, &res)

	return
}

func (c *svcNews) Tags(ctx context.Context) (res []Tag, err error) {
	_req := struct {
	}{}

	err = c.client.call(ctx, "news.Tags", _req, &res)

	return
}

type rpcClient struct {
	endpoint string
	cl       *http.Client

	requestID uint64
}

func newRPCClient(endpoint string, httpClient *http.Client) *rpcClient {
	return &rpcClient{
		endpoint: endpoint,
		cl:       httpClient,
	}
}

func (rc *rpcClient) call(ctx context.Context, methodName string, request, result interface{}) error {
	// encode params
	bts, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("encode params: %w", err)
	}

	requestID := atomic.AddUint64(&rc.requestID, 1)
	requestIDBts := json.RawMessage(strconv.Itoa(int(requestID)))

	req := zenrpc.Request{
		Version: zenrpc.Version,
		ID:      &requestIDBts,
		Method:  methodName,
		Params:  bts,
	}

	ctx = appkit.NewCallerNameContext(ctx, name)

	res, err := rc.Exec(ctx, req)
	if err != nil {
		return err
	}

	if res == nil {
		return nil
	}

	if res.Error != nil {
		return res.Error
	}

	if res.Result == nil {
		return nil
	}

	if result == nil {
		return nil
	}

	return json.Unmarshal(*res.Result, result)
}

// Exec makes http request to jsonrpc endpoint and returns json rpc response.
func (rc *rpcClient) Exec(ctx context.Context, rpcReq zenrpc.Request) (*zenrpc.Response, error) {
	if appkit.NotificationFromContext(ctx) {
		rpcReq.ID = nil
	}

	c, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, fmt.Errorf("json marshal call failed: %w", err)
	}

	buf := bytes.NewReader(c)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rc.endpoint, buf)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	appkit.SetXRequestIDFromCtx(ctx, req)

	// Do request
	resp, err := rc.cl.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, fmt.Errorf("make request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response (%d)", resp.StatusCode)
	}

	var zresp zenrpc.Response
	if rpcReq.ID == nil {
		return &zresp, nil
	}

	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response body (%s) read failed: %w", bb, err)
	}

	if err = json.Unmarshal(bb, &zresp); err != nil {
		return nil, fmt.Errorf("json decode failed (%s): %w", bb, err)
	}

	return &zresp, nil
}
