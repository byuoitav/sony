package bravia

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"go.uber.org/zap"
)

const (
	_maxID = 0x7FFFFFFF
)

type Display struct {
	Address      string
	PreSharedKey string
	Log          *zap.Logger
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

func (r *response) ErrorCode() (int, bool) {
	if len(r.Error) < 1 {
		return 0, false
	}

	// status code is the first item in the list
	code, ok := r.Error[0].(float64)
	if !ok {
		return 0, false
	}

	return int(code), true
}

func (r *response) ErrorReason() string {
	if len(r.Error) < 2 {
		return ""
	}

	// reason is the second item in the list
	reason, ok := r.Error[1].(string)
	if !ok {
		return ""
	}

	return reason
}

func (d *Display) doRequest(ctx context.Context, service string, req request) ([]interface{}, error) {
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

	if code, ok := response.ErrorCode(); ok {
		return nil, fmt.Errorf("code %v: %s", code, response.ErrorReason())
	}

	if wrapped.ID != response.ID {
		return nil, fmt.Errorf("incorrect response id")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http code %v", resp.StatusCode)
	}

	return response.Result, nil
}
