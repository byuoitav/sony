package adcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// HardwareInfo contains the common information for device hardware information
type HardwareInfo struct {
	ModelName     string
	SerialNumber  string
	FilterStatus  string
	WarningStatus []string
	ErrorStatus   []string
	PowerStatus   string
	TimerInfo     []map[string]int
	IPAddress     string
	MACAddress    string
	Gateway       string
	DNS           []string
}

var (
	modelName   = []byte("modelname ?\r\n")
	serialNum   = []byte("serialnum ?\r\n")
	ipAddr      = []byte("ipv4_ip_address ?\r\n")
	gateway     = []byte("ipv4_default_gateway ?\r\n")
	dns         = []byte("ipv4_dns_server1 ?\r\n")
	dns2        = []byte("ipv4_dns_server2 ?\r\n")
	macAddr     = []byte("mac_address ?\r\n")
	filter      = []byte("filter_status ?\r\n")
	warnings    = []byte("warning ?\r\n")
	cmdErrors   = []byte("error ?\r\n")
	powerStatus = []byte("power_status ?\r\n")
	timer       = []byte("timer ?\r\n")
)

// Info returns the hardware information of the projector
func (p *Projector) Info(ctx context.Context) (interface{}, error) {
	var info HardwareInfo

	// model name
	resp, err := p.SendCommand(ctx, p.Address, modelName)
	if err != nil {
		return info, err
	}

	info.ModelName = strings.Trim(resp, "\"")

	// ip address
	resp, err = p.SendCommand(ctx, p.Address, ipAddr)
	if err != nil {
		return info, err
	}

	info.IPAddress = strings.Trim(resp, "\"")

	// gateway
	resp, err = p.SendCommand(ctx, p.Address, gateway)
	if err != nil {
		return info, err
	}

	info.Gateway = strings.Trim(resp, "\"")

	// dns
	resp, err = p.SendCommand(ctx, p.Address, dns)
	if err != nil {
		return info, err
	}

	info.DNS = append(info.DNS, strings.Trim(resp, "\""))

	resp, err = p.SendCommand(ctx, p.Address, dns2)
	if err != nil {
		return info, err
	}

	info.DNS = append(info.DNS, strings.Trim(resp, "\""))

	// mac address
	resp, err = p.SendCommand(ctx, p.Address, macAddr)
	if err != nil {
		return info, err
	}

	info.MACAddress = strings.Trim(resp, "\"")

	// serial number
	resp, err = p.SendCommand(ctx, p.Address, serialNum)
	if err != nil {
		return info, err
	}

	info.SerialNumber = strings.Trim(resp, "\"")

	// filter status
	resp, err = p.SendCommand(ctx, p.Address, filter)
	if err != nil {
		return info, err
	}

	info.FilterStatus = strings.Trim(resp, "\"")

	// power status
	resp, err = p.SendCommand(ctx, p.Address, powerStatus)
	if err != nil {
		return info, err
	}

	info.PowerStatus = strings.Trim(resp, "\"")

	// warnings
	resp, err = p.SendCommand(ctx, p.Address, warnings)
	if err != nil {
		return info, err
	}

	err = json.Unmarshal([]byte(resp), &info.WarningStatus)
	if err != nil {
		return info, err
	}

	// errors
	resp, err = p.SendCommand(ctx, p.Address, cmdErrors)
	if err != nil {
		return info, err
	}

	err = json.Unmarshal([]byte(resp), &info.ErrorStatus)
	if err != nil {
		return info, err
	}

	// timer info
	resp, err = p.SendCommand(ctx, p.Address, timer)
	if err != nil {
		return info, err
	}

	err = json.Unmarshal([]byte(resp), &info.TimerInfo)
	if err != nil {
		return info, err
	}

	return info, nil
}

func (p *Projector) Healthy(ctx context.Context) error {
	_, err := p.Power(ctx)
	if err != nil {
		return fmt.Errorf("failed health check: %s", err)
	}

	return nil
}
