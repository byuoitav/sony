package bravia

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type SonyBaseResult struct {
	ID     int                 `json:"id"`
	Result []map[string]string `json:"result"`
	Error  []interface{}       `json:"error"`
}

//Blank gets the blanked status
func (t *TV) Blank(ctx context.Context) (bool, error) {

	var blanked bool

	payload := SonyTVRequest{
		Params:  []map[string]interface{}{},
		Method:  "getPowerSavingMode",
		Version: "1.0",
		ID:      1,
	}

	t.Log.Info("%+v", zap.Any("payload", payload))

	resp, err := t.PostHTTPWithContext(ctx, "system", payload)
	if err != nil {
		t.Log.Info("ERROR: %v", zap.String("error", err.Error()))
		return false, err
	}

	re := SonyBaseResult{}
	err = json.Unmarshal(resp, &re)
	if err != nil {
		return false, errors.New(fmt.Sprintf("failed to unmarshal response from tv: %s", err))
	}

	// make sure there is a result
	if len(re.Result) == 0 {
		return blanked, errors.New(fmt.Sprintf("error response from tv: %s", re.Error))
	}

	if val, ok := re.Result[0]["mode"]; ok {
		if val == "pictureOff" {
			blanked = true
		} else {
			blanked = false
		}
	}

	return blanked, nil
}

func (t *TV) SetBlank(ctx context.Context, blanked bool) error {
	var blankcmd string
	if blanked == true {
		blankcmd = "pictureOff"
	} else if blanked == false {
		blankcmd = "pictureOn"
	}

	params := make(map[string]interface{})
	params["mode"] = blankcmd

	payload := SonyTVRequest{
		Params:  []map[string]interface{}{params},
		Method:  "setPowerSavingMode",
		Version: "1.0",
		ID:      1,
	}

	t.Log.Info("%+v", zap.Any("payload", payload))

	resp, err := t.PostHTTPWithContext(ctx, "system", payload)
	if err != nil {
		t.Log.Info("ERROR: %v", zap.String("error", err.Error()))
		return err
	}

	re := SonyBaseResult{}
	err = json.Unmarshal(resp, &re)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to unmarshal response from tv: %s", err))
	}

	// make sure there is a result
	if len(re.Result) == 0 {
		return errors.New(fmt.Sprintf("error response from tv: %s", re.Error))
	}

	return nil
}
