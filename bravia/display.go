package bravia

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

const (
	_maxID = 0x7FFFFFFF
)

type Display struct {
	Address      string
	PreSharedKey string
	Log          *zap.Logger

	RequestDelay time.Duration

	once    sync.Once
	limiter *rate.Limiter
}

func (d *Display) init() {
	d.limiter = rate.NewLimiter(rate.Every(d.RequestDelay), 1)
}

type request struct {
	Method  string                   `json:"method"`
	Version string                   `json:"version"`
	Params  []map[string]interface{} `json:"params"`
}

type response struct {
	ID     int           `json:"id"`
	Result []interface{} `json:"result"`
	Error  []interface{} `json:"error"`
}

func (r *response) BuildError() error {
	if len(r.Error) == 0 {
		return nil
	}

	err := &Error{
		code:   -1,
		reason: fmt.Sprintf("unable to parse error %+v", r.Error),
	}

	code, ok := r.Error[0].(float64)
	if !ok {
		return err
	}
	err.code = int(code)

	if len(r.Error) < 2 {
		return err
	}

	reason, ok := r.Error[1].(string)
	if !ok {
		return err
	}
	err.reason = reason

	return err
}

func (d *Display) doRequest(ctx context.Context, service string, req request) ([]interface{}, error) {
	d.once.Do(d.init)

	wrapped := struct {
		request
		ID int `json:"id"`
	}{
		request: req,
		ID:      rand.Intn(_maxID-1) + 1,
	}

	// add an id to the request
	body, err := json.Marshal(wrapped)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s/sony/%s", d.Address, service)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Auth-PSK", d.PreSharedKey)

	if err := d.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("unable to wait for ratelimit: %w", err)
	}

	d.Log.Debug("Doing request", zap.String("url", httpReq.URL.String()), zap.ByteString("body", body))

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("unable to do request: %w", err)
	}
	defer resp.Body.Close()

	var response response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("unable to decode response: %w", err)
	}

	d.Log.Debug("Response", zap.Any("resp", response))

	if wrapped.ID != response.ID {
		return nil, fmt.Errorf("incorrect response id")
	}

	if err := response.BuildError(); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http code %v", resp.StatusCode)
	}

	return response.Result, nil
}
