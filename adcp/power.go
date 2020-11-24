package adcp

import (
	"context"
	"fmt"
)

var (
	// PowerStatus gets the projector's power status
	PowerStatus = []byte("power_status ?\r\n")

	// PowerOn powers on the projector
	PowerOn = []byte("power \"on\"\r\n")

	// PowerStandby powers off the projector
	PowerStandby = []byte("power \"off\"\r\n")
)

// Power returns the status of the projector
func (p *Projector) Power(ctx context.Context) (bool, error) {
	state := false

	resp, err := p.SendCommand(ctx, p.Address, PowerStatus)
	if err != nil {
		return false, err
	}

	switch resp {
	case `"startup"`:
		state = true
	case `"on"`:
		state = true
	case `"standby"`:
	case `"cooling1"`:
	case `"cooling2"`:
	case `"saving_cooling1"`:
	case `"saving_cooling2"`:
	case `"saving_standby"`:
	default:
		return state, fmt.Errorf("unknown power state '%s'", resp)
	}

	return state, nil
}

// SetPower sets the status of the projector
func (p *Projector) SetPower(ctx context.Context, power bool) error {
	cmd := PowerOn
	if !power {
		cmd = PowerStandby
	}

	resp, err := p.SendCommand(ctx, p.Address, cmd)
	if err != nil {
		return err
	}

	return ResponseError(resp)
}
