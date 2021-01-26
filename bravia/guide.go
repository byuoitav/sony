package bravia

import (
	"context"
	"fmt"
)

func (d *Display) getSupportedAPIInfo(ctx context.Context) error {
	req := request{
		Version: "1.0",
		Method:  "getSupportedApiInfo",
		Params: []map[string]interface{}{
			{
				"services": nil,
			},
		},
	}

	res, err := d.doRequest(ctx, "guide", req)
	switch {
	case err != nil:
		return err
	case len(res) < 1:
		return fmt.Errorf("unexpected response: %+v", res)
	}

	return nil
}

