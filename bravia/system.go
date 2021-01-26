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
