package bravia

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.uber.org/zap"
)

func (t *TV) Power(ctx context.Context) (bool, error) {
	var output bool

	payload := SonyTVRequest{
		Params: []map[string]interface{}{},
		Method: "getPowerStatus", Version: "1.0",
		ID: 1,
	}

	response, err := t.PostHTTPWithContext(ctx, "system", payload)
	if err != nil {
		return false, err
	}

	powerStatus := string(response)
	if strings.Contains(powerStatus, "active") {
		output = true
	} else if strings.Contains(powerStatus, "standby") {
		output = false
	} else {
		return false, errors.New("Error getting power status")
	}

	return output, nil
}

func (t *TV) SetPower(ctx context.Context, power bool) error {
	params := make(map[string]interface{})
	params["status"] = power

	payload := SonyTVRequest{
		Params:  []map[string]interface{}{params},
		Method:  "setPowerStatus",
		Version: "1.0",
		ID:      1,
	}

	t.Log.Info("Setting power to", zap.Bool("power", power))

	_, err := t.PostHTTPWithContext(ctx, "system", payload)
	if err != nil {
		return err
	}

	// wait for the display to turn on
	ticker := time.NewTicker(256 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return errors.New("context timed out while waiting for display to turn on")
		case <-ticker.C:
			p, err := t.Power(ctx)
			if err != nil {
				return err
			}

			t.Log.Info("Waiting for display power to change", zap.Bool("target", power), zap.Bool("current", p))

			switch {
			case p && power:
				return nil
			case !p && !power:
				return nil
			}
		}
	}
}
