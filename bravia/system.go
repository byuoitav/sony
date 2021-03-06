package bravia

import (
	"context"
	"fmt"
	"time"
)

func (d *Display) Power(ctx context.Context) (bool, error) {
	req := request{
		Version: "1.0",
		Method:  "getPowerStatus",
		Params:  []map[string]interface{}{},
	}

	res, err := d.doRequest(ctx, "system", req)
	switch {
	case err != nil:
		return false, err
	case len(res) < 1:
		return false, fmt.Errorf("unexpected response: %+v", res)
	}

	m, ok := res[0].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("unexpected response: %+v", res)
	}

	str, ok := m["status"].(string)
	if !ok {
		return false, fmt.Errorf("unexpected response: %+v", res)
	}

	return str == "active", nil
}

func (d *Display) SetPower(ctx context.Context, power bool) error {
	req := request{
		Version: "1.0",
		Method:  "setPowerStatus",
		Params: []map[string]interface{}{
			{
				"status": power,
			},
		},
	}

	_, err := d.doRequest(ctx, "system", req)
	if err != nil {
		return err
	}

	// wait for display to turn on
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pow, err := d.Power(ctx)
			switch {
			case err != nil:
				return fmt.Errorf("unable to confirm power set: %w", err)
			case pow == power:
				return nil
			}
		case <-ctx.Done():
			return fmt.Errorf("unable to confirm power set: %w", ctx.Err())
		}
	}
}

func (d *Display) Blank(ctx context.Context) (bool, error) {
	req := request{
		Version: "1.0",
		Method:  "getPowerSavingMode",
		Params:  []map[string]interface{}{},
	}

	res, err := d.doRequest(ctx, "system", req)
	switch {
	case err != nil:
		return false, err
	case len(res) < 1:
		return false, fmt.Errorf("unexpected response: %+v", res)
	}

	m, ok := res[0].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("unexpected response: %+v", res)
	}

	str, ok := m["mode"].(string)
	if !ok {
		return false, fmt.Errorf("unexpected response: %+v", res)
	}

	return str == "pictureOff", nil
}

func (d *Display) SetBlank(ctx context.Context, blanked bool) error {
	state := "off"
	if blanked {
		state = "pictureOff"
	}

	req := request{
		Version: "1.0",
		Method:  "setPowerSavingMode",
		Params: []map[string]interface{}{
			{
				"mode": state,
			},
		},
	}

	_, err := d.doRequest(ctx, "system", req)
	if err != nil {
		return err
	}

	// wait for display to blank
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			b, err := d.Blank(ctx)
			switch {
			case err != nil:
				return fmt.Errorf("unable to confirm blank set: %w", err)
			case b == blanked:
				return nil
			}
		case <-ctx.Done():
			return fmt.Errorf("unable to confirm blank set: %w", ctx.Err())
		}
	}
}

type Info struct {
	Product    string `json:"product"`
	Language   string `json:"language"`
	Model      string `json:"model"`
	Serial     string `json:"serial"`
	MACAddress string `json:"macAddr"`
	Name       string `json:"name"`
	Generation string `json:"generation"`
}

func (d *Display) Info(ctx context.Context) (interface{}, error) {
	req := request{
		Version: "1.0",
		Method:  "getSystemInformation",
		Params:  []map[string]interface{}{},
	}

	res, err := d.doRequest(ctx, "system", req)
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

	var info Info

	info.Product, ok = m["product"].(string)
	if !ok {
		return info, fmt.Errorf("unexpected response: %+v", res)
	}

	info.Language, ok = m["language"].(string)
	if !ok {
		return info, fmt.Errorf("unexpected response: %+v", res)
	}

	info.Model, ok = m["model"].(string)
	if !ok {
		return info, fmt.Errorf("unexpected response: %+v", res)
	}

	info.Serial, ok = m["serial"].(string)
	if !ok {
		return info, fmt.Errorf("unexpected response: %+v", res)
	}

	info.MACAddress, ok = m["macAddr"].(string)
	if !ok {
		return info, fmt.Errorf("unexpected response: %+v", res)
	}

	info.Name, ok = m["name"].(string)
	if !ok {
		return info, fmt.Errorf("unexpected response: %+v", res)
	}

	info.Generation, ok = m["generation"].(string)
	if !ok {
		return info, fmt.Errorf("unexpected response: %+v", res)
	}

	return info, nil
}

func (d *Display) Healthy(ctx context.Context) error {
	//get list of supported apis and check to see if the newer api is supported

	req := request{
		Version: "1.0",
		Method:  "getSupportedApiInfo",
		Params: []map[string]interface{}{
			{
				"services": nil,
			},
		},
	}

	found := false
	done := false

	res, err := d.doRequest(ctx, "guide", req)
	switch {
	case err != nil:
		return err
	case len(res) < 1:
		return fmt.Errorf("unexpected response: %+v", res)
	}

	list := res[0].([]interface{})

	for _, item := range list {
		if done {
			break
		}

		tmp, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		apis, ok := tmp["apis"].([]interface{})
		if !ok {
			continue
		}

		for _, api := range apis {
			if done {
				break
			}
			a, ok := api.(map[string]interface{})
			if !ok {
				continue
			}

			if a["name"] == "setAudioVolume" {
				versions, ok := a["versions"].([]interface{})
				if !ok {
					return fmt.Errorf("unable to convert setAudioVolume versions")
				}

				for _, version := range versions {
					v, ok := version.(map[string]interface{})
					if !ok {
						fmt.Errorf("unable to convert setAudioVolume version")
					}

					if v["version"] == "1.2" {
						found = true
					}
				}

				done = true
				break
			}
		}
	}

	if found {
		return nil
	}

	return fmt.Errorf("unsupported api")
}
