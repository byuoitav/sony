package bravia

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

func (d *Display) AudioVideoInputs(ctx context.Context) (map[string]string, error) {
	req := request{
		Version: "1.0",
		Method:  "getPlayingContentInfo",
		Params:  []map[string]interface{}{},
	}

	res, err := d.doRequest(ctx, "avContent", req)
	switch {
	case err != nil:
		var bErr *Error
		if errors.As(err, &bErr) {
			switch bErr.code {
			case _displayOff:
				return nil, nil
			}
		}

		return nil, err
	case len(res) < 1:
		return nil, fmt.Errorf("unexpected response: %+v", res)
	}

	m, ok := res[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response: %+v", res)
	}

	str, ok := m["uri"].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected response: %+v", res)
	}

	return map[string]string{
		"": strings.TrimPrefix(str, "extInput:"),
	}, nil
}

// SetAudioVideoInput sets the input of the display to the given input. Input format is everything tha comes after 'extInput:'. Examples can be found at https://pro-bravia.sony.net/develop/integrate/rest-api/spec/resource-uri-list/index.html.
func (d *Display) SetAudioVideoInput(ctx context.Context, _, input string) error {
	req := request{
		Version: "1.0",
		Method:  "setPlayContent",
		Params: []map[string]interface{}{
			{
				"uri": fmt.Sprintf("extInput:%s", input),
			},
		},
	}

	_, err := d.doRequest(ctx, "avContent", req)
	if err != nil {
		return err
	}

	// wait for input to change
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			inputs, err := d.AudioVideoInputs(ctx)
			switch {
			case err != nil:
				return fmt.Errorf("unable to confirm input set: %w", err)
			case inputs[""] == input:
				return nil
			}
		case <-ctx.Done():
			return fmt.Errorf("unable to confirm power set: %w", ctx.Err())
		}
	}
}

type inputStatus struct {
	URI        string `json:"uri"`
	Title      string `json:"title"`
	Connection bool   `json:"connection"`
	Label      string `json:"label"`
	Icon       string `json:"icon"`
	Status     string `json:"status"`
}

func (d *Display) getCurrentExternalInputsStatus(ctx context.Context) ([]inputStatus, error) {
	req := request{
		Version: "1.1",
		Method:  "getCurrentExternalInputsStatus",
		Params:  []map[string]interface{}{},
	}

	res, err := d.doRequest(ctx, "avContent", req)
	switch {
	case err != nil:
		return nil, err
	case len(res) < 1:
		return nil, fmt.Errorf("unexpected response: %+v", res)
	}

	list, ok := res[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response: %+v", res)
	}

	var statuses []inputStatus
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected response: %+v", res)
		}

		var status inputStatus

		status.URI, ok = m["uri"].(string)
		if !ok {
			return nil, fmt.Errorf("unexpected response: %+v", res)
		}

		status.Title, ok = m["title"].(string)
		if !ok {
			return nil, fmt.Errorf("unexpected response: %+v", res)
		}

		status.Connection, ok = m["connection"].(bool)
		if !ok {
			return nil, fmt.Errorf("unexpected response: %+v", res)
		}

		status.Label, ok = m["label"].(string)
		if !ok {
			return nil, fmt.Errorf("unexpected response: %+v", res)
		}

		status.Icon, ok = m["icon"].(string)
		if !ok {
			return nil, fmt.Errorf("unexpected response: %+v", res)
		}

		status.Status, ok = m["status"].(string)
		if !ok {
			return nil, fmt.Errorf("unexpected response: %+v", res)
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

/*
// GetActiveSignal determines if the current input on the TV is active or not
func GetActiveSignal(address, port string) (structs.ActiveSignal, *nerr.E) {
	var output structs.ActiveSignal

	payload := SonyTVRequest{
		Params:  []map[string]interface{}{},
		Method:  "getCurrentExternalInputsStatus",
		ID:      1,
		Version: "1.1",
	}

	response, err := PostHTTP(address, payload, "avContent")
	if err != nil {
		return output, nerr.Translate(err)
	}

	var outputStruct SonyMultiAVContentResponse
	err = json.Unmarshal(response, &outputStruct)
	if err != nil || len(outputStruct.Result) < 1 {
		return output, nerr.Translate(err)
	}
	//we need to parse the response for the value

	log.L.Debugf("%+v", outputStruct)

	regexStr := `extInput:(.*?)\?port=(.*)`
	re := regexp.MustCompile(regexStr)

	for _, result := range outputStruct.Result[0] {
		if result.Status == "true" {
			matches := re.FindStringSubmatch(result.URI)
			tempActive := fmt.Sprintf("%v!%v", matches[1], matches[2])

			output.Active = (tempActive == port)
		}
	}

	return output, nil
}
*/
