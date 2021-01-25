package bravia

import (
	"context"
	"fmt"
	"strings"
)

func (d *Display) AudioVideoInputs(ctx context.Context) (map[string]string, error) {
	req := request{
		Version: "1.0",
		Method:  "getPlayingContentInfo",
	}

	res, err := d.doRequest(ctx, "avContent", req)
	switch {
	case err != nil:
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

func (d *Display) SetAudioVideoInput(ctx context.Context, output, input string) error {
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
	return err
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
