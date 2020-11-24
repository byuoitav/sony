package adcp

import (
	"errors"
	"fmt"
)

var responseError = map[string]error{
	"ok":            nil,
	"err_cmd":       errors.New("command format error"),
	"err_option":    errors.New("command option error"),
	"err_inactive":  errors.New("command is temporarily invalid"),
	"err_val":       errors.New("value for command is out of range"),
	"err_auth":      errors.New("network authentication error"),
	"err_internal1": errors.New("internal communication error 1 of the projector"),
	"err_internal2": errors.New("internal communication error 2 of the projector"),
}

func ResponseError(resp string) error {
	if err, ok := responseError[resp]; ok {
		return err
	}

	return fmt.Errorf("unknown response error %q", resp)
}
